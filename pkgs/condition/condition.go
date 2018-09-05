package condition

import (
	"reflect"

	"github.com/uphy/elastalert-mail-gateway/pkgs/jsonutil"
)

type (
	Condition interface {
		Eval(scope Scope) (bool, error)
	}
	Scope interface {
		Get(name string) *jsonutil.Value
	}
	Conditions struct {
		Match      *Match      `json:"match,omitempty"`
		Constant   *Constant   `json:"constant,omitempty"`
		Wildcard   *Wildcard   `json:"wildcard,omitempty"`
		Not        *Not        `json:"not,omitempty"`
		And        *And        `json:"and,omitempty"`
		Eq         *Eq         `json:"eq,omitempty"`
		Gt         *Gt         `json:"gt,omitempty"`
		Gte        *Gte        `json:"gte,omitempty"`
		Lt         *Lt         `json:"lt,omitempty"`
		Lte        *Lte        `json:"lte,omitempty"`
		conditions []Condition `json:"-"`
	}
)

func (c *Conditions) Eval(scope Scope) (bool, error) {
	if c.conditions == nil {
		t := reflect.TypeOf(c).Elem()
		v := reflect.ValueOf(*c)
		conditionInterface := reflect.TypeOf((*Condition)(nil)).Elem()
		conditions := make([]Condition, 0)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Type.AssignableTo(conditionInterface) {
				value := v.Field(i)
				if value.IsNil() {
					continue
				}
				cond := value.Interface().(Condition)
				if cond != nil {
					conditions = append(conditions, cond)
				}
			}
		}
		c.conditions = conditions
	}
	for _, cond := range c.conditions {
		v, err := cond.Eval(scope)
		if err != nil {
			return false, err
		}
		if !v {
			return false, nil
		}
	}
	return true, nil
}
