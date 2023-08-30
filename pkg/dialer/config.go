package dialer

import (
	"encoding/base64"
	"errors"
	"fmt"

	_ "gopkg.in/yaml.v3"
)

type Interface struct {
	Address    string `yaml:"address"`
	PrivateKey string `yaml:"private_key"`
	Dns        string `yaml:"dns"`
}

type Peer struct {
	PublicKey string `yaml:"public_key"`
	Endpoint  string `yaml:"endpoint"`
	AllowedIP string `yaml:"allowedip"`
}

func base64KeyToHex(in string) (string, error) {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Invalid base64: %s", in), err)
	}

	return fmt.Sprintf("%x", out), nil
}

func (p *Peer) toIpcString() (string, error) {
	key, err := base64KeyToHex(p.PublicKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("public_key=%s\nallowed_ip=%s\nendpoint=%s\n", key, p.AllowedIP, p.Endpoint), nil
}
