package gateway

import (
	"fmt"
	"io/ioutil"

	"github.com/uphy/elastalert-mail-gateway/pkgs/elastalert"
	"github.com/uphy/elastalert-mail-gateway/pkgs/mail"
	"github.com/uphy/elastalert-mail-gateway/pkgs/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Gateway struct {
	smtpClient *mail.SMTPClient
	alerters   []Alerter
	logger     *zap.Logger
}

func New(smtpClient *mail.SMTPClient, alerters []Alerter, logger *zap.Logger) *Gateway {
	return &Gateway{smtpClient, alerters, logger}
}

func (g *Gateway) Start(address string, port int) error {
	s := server.New(address, port)
	if err := s.Start(); err != nil {
		return err
	}
	g.logger.Info("Server started.", zap.String("address", address), zap.Int("port", port))

	for receivedMail := range s.Mails {
		g.logger.Info("Got mail.", zap.Object("header", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
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
			g.logger.Error("Failed to read received mail body.", zap.Error(err))
			continue
		}
		// parse mail body; extract documents from the email.
		alerts, err := elastalert.ParseMailBody(string(b))
		if err != nil {
			g.logger.Error("Failed to parse mail body.", zap.Error(err))
			continue
		}
		// convert mail to map[string]interface{}
		mailObj := receivedMail.Map()

		// Determine recipients based on config file for each alerts
		mails := []*mail.Mail{}
		ctx := &AlertContext{
			Logger:           g.logger,
			ReceivedMailJSON: mailObj,
			ReceivedMail:     &receivedMail,
		}
		for _, al := range alerts {
			g.logger.Debug("Processing alert.", zap.String("alert", al.String()))
			for _, alerter := range g.alerters {
				alertMails, err := alerter.Alert(ctx, al)
				if err != nil {
					g.logger.Error("Failed to process alert.", zap.Error(err))
				}
				mails = append(mails, alertMails...)
			}
		}
		// aggregate mails by recipients
		mails, err = mail.AggregateMails(mails)
		if err != nil {
			g.logger.Error("Failed to aggregate mail", zap.Error(err))
			continue
		}

		// send mails
		for _, m := range mails {
			g.logger.Debug("Sending mail...", zap.String("mail", m.JSON()))
			if err := g.smtpClient.Send(m); err != nil {
				g.logger.Error("Failed to send mail.", zap.String("mail", m.JSON()), zap.Error(err))
			}
		}
	}
	return nil
}
