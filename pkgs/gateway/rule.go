package gateway

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/uphy/elastalert-mail-gateway/pkgs/condition"
	"github.com/uphy/elastalert-mail-gateway/pkgs/elastalert"
	"github.com/uphy/elastalert-mail-gateway/pkgs/jsonutil"
	"github.com/uphy/elastalert-mail-gateway/pkgs/mail"
	"go.uber.org/zap"
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
	RulesAlerter struct {
		Rules Rules
	}
)

func NewRulesAlerter(rules Rules) *RulesAlerter {
	return &RulesAlerter{rules}
}

func (r *RulesAlerter) Alert(ctx *AlertContext, alert *elastalert.Alert) ([]*mail.Mail, error) {
	mails := []*mail.Mail{}
	scope := jsonutil.Object(map[string]interface{}{
		"doc":  alert.Doc,
		"body": alert.Body,
		"mail": ctx.ReceivedMailJSON,
	})
	alertString := alert.String()
	hasMatch := false
	for _, rule := range r.Rules {
		match, err := rule.MatchCondition(scope)
		if err != nil {
			ctx.Logger.Error("failed to evaluate match condition.", zap.String("scope", scope.String()), zap.Error(err))
			continue
		}
		if !match {
			continue
		}
		m := mergeMail(ctx.ReceivedMail, &rule, strings.NewReader(alertString))
		ctx.Logger.Debug("Matched.", zap.String("mail", m.JSON()))
		mails = append(mails, m)
		hasMatch = true
	}
	if !hasMatch {
		defaultRule := r.Rules.DefaultRule()

		m := mergeMail(ctx.ReceivedMail, defaultRule, strings.NewReader(alertString))
		ctx.Logger.Debug("No rules matched.  Applying default rule.", zap.String("mail", m.JSON()))
		mails = append(mails, m)
	}
	return mails, nil
}

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
