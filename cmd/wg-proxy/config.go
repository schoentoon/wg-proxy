package main

import (
	"os"

	"gitlab.com/schoentoon/wg-proxy/pkg/dialer"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Interface dialer.Interface `yaml:"interface"`
	Peers     []dialer.Peer    `yaml:"peer"`

	Proxy struct {
		HTTP struct {
			Addr string `yaml:"addr"`
		} `yaml:"http"`
		Socks5 struct {
			Addr string `yaml:"addr"`
		} `yaml:"socks5"`
	} `yaml:"proxy"`

	Debug   bool   `yaml:"debug"`
	Metrics string `yaml:"metrics"`
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
