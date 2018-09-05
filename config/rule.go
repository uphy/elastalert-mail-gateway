package config

import (
	"encoding/json"
	"fmt"

	"github.com/uphy/elastalert-mail-gateway/pkgs/condition"
)

type (
	Rules []Rule
	Rule  struct {
		To        *AddressList          `json:"to,omitempty"`
		Cc        *AddressList          `json:"cc,omitempty"`
		ReplyTo   *Address              `json:"reply_to,omitempty"`
		Condition *condition.Conditions `json:"condition,omitempty"`
		Default   bool                  `json:"default,omitempty"`
	}
)

func (r Rules) DefaultRule() *Rule {
	for _, mail := range r {
		if mail.Default {
			return &mail
		}
	}
	return nil
}

func (r *Rule) MatchCondition(scope condition.Scope) (bool, error) {
	return r.Condition.Eval(scope)
}

func (r Rule) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<failed to unmarshal:%v>", err)
	}
	return string(b)
}
