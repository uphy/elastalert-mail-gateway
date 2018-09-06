package mail

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/uphy/elastalert-mail-gateway/pkgs/elastalert"
)

func AggregateMails(mails []*Mail) ([]*Mail, error) {
	type Key struct {
		To string
		Cc string
	}
	m := map[Key][]*Mail{}
	for _, mail := range mails {
		k := Key{
			To: addressesToString(mail.To),
			Cc: addressesToString(mail.Cc),
		}
		m[k] = append(m[k], mail)
	}
	ret := []*Mail{}
	for _, aggregatingMails := range m {
		if len(aggregatingMails) == 0 {
			ret = append(ret, aggregatingMails[0])
		} else {
			a, err := aggregateMail(aggregatingMails)
			if err != nil {
				return nil, err
			}
			ret = append(ret, a)
		}
	}
	return ret, nil
}

func aggregateMail(mails []*Mail) (*Mail, error) {
	buf := new(bytes.Buffer)
	for _, mail := range mails {
		b, err := ioutil.ReadAll(mail.Body)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
		buf.WriteString(elastalert.AlertSeparator)
	}
	base := mails[0]
	return &Mail{
		From:    base.From,
		To:      base.To,
		Cc:      base.Cc,
		ReplyTo: base.ReplyTo,
		Subject: base.Subject,
		Body:    strings.NewReader(buf.String()),
	}, nil
}

func addressesToString(addresses []*Address) string {
	buf := new(bytes.Buffer)
	for _, addr := range addresses {
		buf.WriteString(addr.Address.Address)
	}
	return buf.String()
}
