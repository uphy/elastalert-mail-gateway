package condition

import (
	"path/filepath"

	"github.com/uphy/elastalert-mail-gateway/pkgs/jsonutil"
)

type (
	Match struct {
		jsonutil.KeyValue
	}
	Wildcard struct {
		jsonutil.KeyValue
	}
	Not struct {
		Conditions
	}
	And                 []*Conditions
	Or                  []*Conditions
	arithmeticCondition struct {
		jsonutil.KeyValue
	}
	Eq struct {
		arithmeticCondition
	}
	Gt struct {
		arithmeticCondition
	}
	Gte struct {
		arithmeticCondition
	}
	Lt struct {
		arithmeticCondition
	}
	Lte struct {
		arithmeticCondition
	}
	Constant bool
)

func (m *Match) Eval(scope Scope) (bool, error) {
	v := scope.Get(m.Key)
	if v == nil {
		return false, nil
	}
	return v.Equals(m.Value), nil
}

func (c Constant) Eval(scope Scope) (bool, error) {
	return bool(c), nil
}

func (w *Wildcard) Eval(scope Scope) (bool, error) {
	v := scope.Get(w.Key)
	if v == nil {
		return false, nil
	}
	pattern := w.Value.String()
	value := v.String()
	return filepath.Match(pattern, value)
}

func (n *Not) Eval(scope Scope) (bool, error) {
	b, err := n.Conditions.Eval(scope)
	if err != nil {
		return false, err
	}
	return !b, nil
}

func (a And) Eval(scope Scope) (bool, error) {
	for _, c := range a {
		matched, err := c.Eval(scope)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func (o Or) Eval(scope Scope) (bool, error) {
	for _, c := range o {
		matched, err := c.Eval(scope)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func (a *arithmeticCondition) values(scope Scope) (float64, float64, error) {
	f1, err := a.Value.Float64()
	if err != nil {
		return 0, 0, err
	}
	v := scope.Get(a.Key)
	if err != nil {
		return 0, 0, err
	}
	f2, err := v.Float64()
	return f1, f2, nil
}

func (e *Eq) Eval(scope Scope) (bool, error) {
	f1, f2, err := e.arithmeticCondition.values(scope)
	if err != nil {
		return false, err
	}
	return f1 == f2, nil
}

func (e *Gt) Eval(scope Scope) (bool, error) {
	f1, f2, err := e.arithmeticCondition.values(scope)
	if err != nil {
		return false, err
	}
	return f1 > f2, nil
}

func (e *Gte) Eval(scope Scope) (bool, error) {
	f1, f2, err := e.arithmeticCondition.values(scope)
	if err != nil {
		return false, err
	}
	return f1 >= f2, nil
}

func (e *Lt) Eval(scope Scope) (bool, error) {
	f1, f2, err := e.arithmeticCondition.values(scope)
	if err != nil {
		return false, err
	}
	return f1 < f2, nil
}

func (e *Lte) Eval(scope Scope) (bool, error) {
	f1, f2, err := e.arithmeticCondition.values(scope)
	if err != nil {
		return false, err
	}
	return f1 <= f2, nil
}
