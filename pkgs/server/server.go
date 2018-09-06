package server

import (
	"encoding/base64"
	"fmt"
	gomail "net/mail"
	"strings"

	"github.com/flashmob/go-guerrilla/backends"

	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/mail"
	umail "github.com/uphy/elastalert-mail-gateway/pkgs/mail"
)

type Server struct {
	addr  string
	port  int
	Mails chan umail.Mail
}

func New(addr string, port int) *Server {
	mails := make(chan umail.Mail)
	return &Server{addr, port, mails}
}

func (s *Server) Start() error {
	cfg := &guerrilla.AppConfig{
		BackendConfig: backends.BackendConfig{
			"save_process": "gateway",
		},
		LogFile: "server.log",
	}
	cfg.AllowedHosts = []string{"."}
	sc := guerrilla.ServerConfig{
		ListenInterface: fmt.Sprintf("%s:%d", s.addr, s.port),
		IsEnabled:       true,
	}
	cfg.Servers = append(cfg.Servers, sc)
	d := &guerrilla.Daemon{
		Config: cfg,
	}
	d.AddProcessor("gateway", func() backends.Decorator {
		return func(p backends.Processor) backends.Processor {
			return backends.ProcessWith(func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
				if task == backends.TaskValidateRcpt {
					return p.Process(e, task)
				} else if task == backends.TaskSaveMail {
					s.Mails <- *parseMail(e)
					return p.Process(e, task)
				}
				return p.Process(e, task)
			})
		}
	})
	return d.Start()
}

func parseMail(e *mail.Envelope) *umail.Mail {
	msg, _ := gomail.ReadMessage(e.NewReader())

	var (
		from    *umail.Address
		to      []*umail.Address
		cc      []*umail.Address
		replyTo *umail.Address
	)

	fromHeader := msg.Header.Get("From")
	if len(fromHeader) > 0 {
		from, _ = umail.ParseAddress(fromHeader)
	}

	toHeader := msg.Header.Get("To")
	if len(toHeader) > 0 {
		to, _ = umail.ParseAddressList(toHeader)
	} else {
		to = []*umail.Address{}
	}

	ccHeader := msg.Header.Get("Cc")
	if len(ccHeader) > 0 {
		cc, _ = umail.ParseAddressList(ccHeader)
	} else {
		cc = []*umail.Address{}
	}

	replyToHeader := msg.Header.Get("Reply-To")
	if len(replyToHeader) > 0 {
		replyTo, _ = umail.ParseAddress(replyToHeader)
	}

	body := msg.Body
	if "base64" == strings.ToLower(msg.Header.Get("Content-Transfer-Encoding")) {
		body = base64.NewDecoder(base64.StdEncoding, msg.Body)
	}

	subject := msg.Header.Get("Subject")
	return &umail.Mail{
		From:    from,
		To:      to,
		Cc:      cc,
		Subject: subject,
		ReplyTo: replyTo,
		Body:    body,
	}
}
