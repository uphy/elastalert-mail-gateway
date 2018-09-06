package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/smtp"
)

type (
	Mail struct {
		From    *Address   `json:"from"`
		To      []*Address `json:"to"`
		Cc      []*Address `json:"cc"`
		ReplyTo *Address   `json:"reply_to"`
		Subject string     `json:"subject"`
		Body    io.Reader  `json:"-"`
	}
	SMTPClient struct {
		host     string
		port     int
		user     string
		password string
	}
)

func (m *Mail) Map() map[string]interface{} {
	b, _ := json.Marshal(m)
	var v map[string]interface{}
	json.Unmarshal(b, &v)
	return v
}

func (m *Mail) JSON() string {
	b, _ := json.Marshal(m)
	return string(b)
}
func (m *Mail) JSONIndent() string {
	b, _ := json.MarshalIndent(m, "", "   ")
	return string(b)
}

func NewSMTPClient(host string, port int, user string, password string) *SMTPClient {
	return &SMTPClient{host, port, user, password}
}

func (s *SMTPClient) Send(m *Mail) error {
	msg := new(bytes.Buffer)
	msg.WriteString("From: " + m.From.String())
	msg.WriteString("\r\n")
	msg.WriteString("To: " + joinMailAddress(m.To))
	msg.WriteString("\r\n")
	if len(m.Cc) > 0 {
		msg.WriteString("Cc: " + joinMailAddress(m.Cc))
		msg.WriteString("\r\n")
	}
	msg.WriteString("Subject: " + m.Subject)
	msg.WriteString("\r\n\r\n")
	io.Copy(msg, m.Body)

	var auth smtp.Auth
	if len(s.user) > 0 && len(s.password) > 0 {
		auth = smtp.PlainAuth("", s.user, s.password, s.host)
	}
	return smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.port), auth, m.From.Address.Address, recipients(m), msg.Bytes())
}

func recipients(m *Mail) []string {
	recipients := []string{}
	for _, a := range m.To {
		recipients = append(recipients, a.Address.Address)
	}
	for _, a := range m.Cc {
		recipients = append(recipients, a.Address.Address)
	}
	return recipients
}

func joinMailAddress(a []*Address) string {
	buf := new(bytes.Buffer)
	for i, aa := range a {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(aa.Address.Address)
	}
	return buf.String()
}
