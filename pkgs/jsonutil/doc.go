package jsonutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Object map[string]interface{}

func (d Object) Get(name string) *Value {
	names := strings.Split(name, ".")
	return d.find(d, names)
}

func (d *Object) find(m map[string]interface{}, names []string) *Value {
	v, ok := m[names[0]]
	if !ok {
		return nil
	}
	if len(names) > 1 {
		switch vv := v.(type) {
		case map[string]interface{}:
			return d.find(vv, names[1:])
		case Object:
			return d.find(vv, names[1:])
		}
		return nil
	}
	return NewValue(v)
}

func (d Object) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("<failed to marshal:%v>", err)
	}
	return string(b)
}
