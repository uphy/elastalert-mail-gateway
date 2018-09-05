package jsonutil

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Value struct {
	v interface{}
}

func NewValue(v interface{}) *Value {
	return &Value{v}
}

func (v Value) String() string {
	return fmt.Sprint(v.v)
}

func (v *Value) Int64() (int64, error) {
	return strconv.ParseInt(v.String(), 10, 64)
}

func (v *Value) Float64() (float64, error) {
	return strconv.ParseFloat(v.String(), 64)
}

func (v *Value) Bool() bool {
	switch v.String() {
	case "true", "yes":
		return true
	default:
		return false
	}
}

func (v *Value) Equals(o *Value) bool {
	return v.String() == o.String()
}

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.v)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*v = Value{i}
	return nil
}
