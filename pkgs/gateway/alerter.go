package gateway

import (
	"io"

	"github.com/uphy/elastalert-mail-gateway/pkgs/elastalert"
	"github.com/uphy/elastalert-mail-gateway/pkgs/mail"
	"go.uber.org/zap"
)

type (
	Alerter interface {
		Alert(ctx *AlertContext, a *elastalert.Alert) ([]*mail.Mail, error)
	}
	AlertContext struct {
		ReceivedMailJSON map[string]interface{}
		ReceivedMail     *mail.Mail
		Logger           *zap.Logger
	}
)

func mergeMail(received *mail.Mail, rule *Rule, body io.Reader) *mail.Mail {
	// generate mail
	var (
		to      []*mail.Address
		cc      []*mail.Address
		replyTo *mail.Address
	)
	if rule.To != nil {
		to = rule.To.Addresses
	}
	if len(to) == 0 {
		to = received.To
	}
	if rule.Cc != nil {
		cc = rule.Cc.Addresses
	}
	if len(cc) == 0 {
		cc = received.Cc
	}
	if rule.ReplyTo != nil {
		replyTo = rule.ReplyTo.Address
	}
	if replyTo == nil {
		replyTo = received.ReplyTo
	}
	return &mail.Mail{
		From:    received.From,
		To:      to,
		Cc:      cc,
		ReplyTo: replyTo,
		Subject: received.Subject,
		Body:    body,
	}
}
