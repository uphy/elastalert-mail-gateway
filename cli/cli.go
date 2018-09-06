package cli

import (
	"os"

	"go.uber.org/zap"

	"github.com/uphy/elastalert-mail-gateway/config"
	"github.com/uphy/elastalert-mail-gateway/pkgs/gateway"
	"github.com/uphy/elastalert-mail-gateway/pkgs/mail"
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
		rulesAlerter := gateway.NewRulesAlerter(cfg.Rules)
		g := gateway.New(smtpClient, []gateway.Alerter{rulesAlerter}, logger)

		address := ctx.String("address")
		port := ctx.Int("port")
		return g.Start(address, port)
	}

	return app.Run(os.Args)
}
