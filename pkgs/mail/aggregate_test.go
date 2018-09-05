package mail

import (
	"fmt"
	gomail "net/mail"
	"reflect"
	"strings"
	"testing"

	"github.com/uphy/elastalert-mail-gateway/pkgs/elastalert"
)

func newMail(t *testing.T, to string, cc string, subject string, body string) *Mail {
	toaddr, err := gomail.ParseAddress(to)
	if err != nil {
		t.FailNow()
		return nil
	}
	ccaddr, err := gomail.ParseAddress(cc)
	if err != nil {
		t.FailNow()
		return nil
	}
	from, _ := gomail.ParseAddress("from@foo.com")
	replyto, _ := gomail.ParseAddress("reply-to@foo.com")
	return &Mail{
		From:    from,
		To:      []*gomail.Address{toaddr},
		Cc:      []*gomail.Address{ccaddr},
		ReplyTo: replyto,
		Subject: subject,
		Body:    strings.NewReader(body),
	}
}

func TestAggregateMails(t *testing.T) {
	type args struct {
		mails []*Mail
	}
	tests := []struct {
		name    string
		args    args
		want    []*Mail
		wantErr bool
	}{
		{
			name: "aggregation 1",
			args: args{
				mails: []*Mail{
					newMail(t, "to@foo.com", "cc@foo.com", "subject", "body1"),
					newMail(t, "to@foo.com", "cc@foo.com", "subject", "body2"),
				},
			},
			want: []*Mail{
				newMail(t, "to@foo.com", "cc@foo.com", "subject", fmt.Sprintf("body1%sbody2%s", elastalert.AlertSeparator, elastalert.AlertSeparator)),
			},
		},
		{
			name: "no aggregation",
			args: args{
				mails: []*Mail{
					newMail(t, "to@foo.com", "cc@foo.com", "subject", "body1"),
					newMail(t, "to@foo.com", "ccc@foo.com", "subject", "body2"),
				},
			},
			want: []*Mail{
				newMail(t, "to@foo.com", "cc@foo.com", "subject", "body1"),
				newMail(t, "to@foo.com", "ccc@foo.com", "subject", "body2"),
			},
		},
		{
			name: "aggregation 2",
			args: args{
				mails: []*Mail{
					newMail(t, "to@foo.com", "cc@foo.com", "subject", "body1"),
					newMail(t, "to@foo.com", "ccc@foo.com", "subject", "body2"),
					newMail(t, "to@foo.com", "cc@foo.com", "subject", "body3"),
				},
			},
			want: []*Mail{
				newMail(t, "to@foo.com", "cc@foo.com", "subject", fmt.Sprintf("body1%sbody3%s", elastalert.AlertSeparator, elastalert.AlertSeparator)),
				newMail(t, "to@foo.com", "ccc@foo.com", "subject", "body2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AggregateMails(tt.args.mails)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateMails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotjson := mailsToString(got)
			wantjson := mailsToString(tt.want)
			if gotjson != wantjson {
				t.Errorf("AggregateMails() = \n%v\nwant\n %v", gotjson, wantjson)
			}
		})
	}
}

func mailsToString(mails []*Mail) string {
	s := []string{}
	for _, m := range mails {
		s = append(s, mailToString(m))
	}
	return strings.Join(s, ", ")
}
func mailToString(m *Mail) string {
	return m.JSONIndent()
}

func Test_aggregateMail(t *testing.T) {
	type args struct {
		mails []*Mail
	}
	tests := []struct {
		name    string
		args    args
		want    *Mail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aggregateMail(tt.args.mails)
			if (err != nil) != tt.wantErr {
				t.Errorf("aggregateMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aggregateMail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addressesToString(t *testing.T) {
	type args struct {
		addresses []*gomail.Address
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addressesToString(tt.args.addresses); got != tt.want {
				t.Errorf("addressesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
