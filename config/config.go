package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/uphy/elastalert-mail-gateway/pkgs/gateway"
)

type (
	Config struct {
		SMTPHost string        `json:"smtp_host,omitempty"`
		SMTPPort int           `json:"smtp_port,omitempty"`
		Rules    gateway.Rules `json:"rules,omitempty"`
	}
)

func LoadFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

func Load(reader io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) Save(w io.Writer) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
