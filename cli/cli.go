package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	gomail "net/mail"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/uphy/elastalert-mail-gateway/config"
	"github.com/uphy/elastalert-mail-gateway/pkgs/elastalert"
	"github.com/uphy/elastalert-mail-gateway/pkgs/jsonutil"
	"github.com/uphy/elastalert-mail-gateway/pkgs/mail"
	"github.com/uphy/elastalert-mail-gateway/pkgs/server"
	"github.com/urfave/cli"
)

const Version = "0.0.1"

func Run() error {
	app := cli.NewApp()
	app.Name = "elastalert-mail-gateway"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "address,a",
			Value:  "0.0.0.0",
			EnvVar: "GATEWAY_ADDRESS",
		},
		cli.IntFlag{
			Name:   "port,p",
			Value:  2525,
			EnvVar: "GATEWAY_PORT",
		},
	}
	app.ArgsUsage = "[configfile]"
	logger, _ := zap.NewDevelopment()
	app.Action = func(ctx *cli.Context) error {
		configFile := ctx.Args().First()
		cfg, err := config.LoadFile(configFile)
		if err != nil {
			return err
		}

		smtpClient := mail.NewSMTPClient(cfg.SMTPHost, cfg.SMTPPort, "", "")

		address := ctx.String("address")
		port := ctx.Int("port")
		s := server.New(address, port)
		if err := s.Start(); err != nil {
			return err
		}
		logger.Info("Server started.", zap.String("address", address), zap.Int("port", port))

		for receivedMail := range s.Mails {
			logger.Info("Got mail.", zap.Object("header", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				if receivedMail.From != nil {
					enc.AddString("from", receivedMail.From.String())
				}
				if receivedMail.To != nil {
					enc.AddString("to", fmt.Sprint(receivedMail.To))
				}
				if receivedMail.Cc != nil {
					enc.AddString("cc", fmt.Sprint(receivedMail.Cc))
				}
				enc.AddString("subject", receivedMail.Subject)
				return nil
			})))

			// read mail body
			b, err := ioutil.ReadAll(receivedMail.Body)
			if err != nil {
				logger.Error("Failed to read received mail body.", zap.Error(err))
				continue
			}
			// parse mail body; extract documents from the email.
			alerts, err := elastalert.ParseMailBody(string(b))
			if err != nil {
				log.Println(err)
				continue
			}
			// convert mail to map[string]interface{}
			mailObj := receivedMail.Map()

			// Determine recipients based on config file for each alerts
			mails := []*mail.Mail{}
			for _, alert := range alerts {
				logger.Debug("Processing alert.", zap.String("alert", alert.String()))
				scope := jsonutil.Object(map[string]interface{}{
					"doc":  alert.Doc,
					"body": alert.Body,
					"mail": mailObj,
				})
				alertString := alert.String()
				hasMatch := false
				for _, rule := range cfg.Rules {
					match, err := rule.MatchCondition(scope)
					if err != nil {
						logger.Error("failed to send mail.", zap.String("scope", scope.String()))
						continue
					}
					if !match {
						continue
					}
					m := mergeMail(&receivedMail, &rule, strings.NewReader(alertString))
					logger.Debug("Matched.", zap.String("mail", m.JSON()))
					mails = append(mails, m)
					hasMatch = true
				}
				if !hasMatch {
					defaultRule := cfg.Rules.DefaultRule()

					m := mergeMail(&receivedMail, defaultRule, strings.NewReader(alertString))
					logger.Debug("No rules matched.  Applying default rule.", zap.String("mail", m.JSON()))
					mails = append(mails, m)
				}
			}
			// aggregate mails by recipients
			mails, err = mail.AggregateMails(mails)
			if err != nil {
				logger.Error("Failed to aggregate mail", zap.Error(err))
				continue
			}

			// send mails
			for _, m := range mails {
				logger.Debug("Sending mail...", zap.String("mail", m.JSON()))
				if err := smtpClient.Send(m); err != nil {
					logger.Error("Failed to send mail.", zap.String("mail", m.JSON()), zap.Error(err))
				}
			}
		}
		return nil
	}

	return app.Run(os.Args)
}

func mergeMail(received *mail.Mail, rule *config.Rule, body io.Reader) *mail.Mail {
	// generate mail
	var (
		to      []*gomail.Address
		cc      []*gomail.Address
		replyTo *gomail.Address
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

//func sendMail()
