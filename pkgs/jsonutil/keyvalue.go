package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
)

type (
	KeyValue struct {
		Key   string
		Value *Value
	}
)

func (kv KeyValue) MarshalJSON() ([]byte, error) {
	v := map[string]interface{}{
		kv.Key: kv.Value,
	}
	return json.Marshal(v)
}

func (kv *KeyValue) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) != 1 {
		return fmt.Errorf("must be a key value pair: %v", v)
	}
	for k, val := range v {
		*kv = KeyValue{
			Key:   k,
			Value: NewValue(val),
		}
		return nil
	}
	return errors.New("never come here")
}
