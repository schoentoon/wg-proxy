package main

import (
	"os"

	"gitlab.com/schoentoon/wg-proxy/pkg/dialer"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interface dialer.Interface `yaml:"interface"`
	Peers     []dialer.Peer    `yaml:"peer"`

	Debug bool `yaml:"debug"`
}

// ReadConfig reads a file into the config structure
func ReadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	out := &Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&out)
	if err != nil {
		return nil, err
	}

	return out, err
}
