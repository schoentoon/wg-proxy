package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/schoentoon/wg-proxy/pkg/dialer"
)

func main() {
	configFile := flag.String("config", "config.yml", "What config file to use?")
	flag.Parse()

	if *configFile == "" {
		logrus.Fatal("No config file specified")
	}

	cfg, err := ReadConfig(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	if cfg.Debug {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Debugf("%+v", cfg)

	dial, err := dialer.NewDialer(cfg.Interface, cfg.Peers...)
	if err != nil {
		logrus.Fatal(err)
	}

	client := http.Client{
		Transport: &http.Transport{
			DialContext: dial.DialContext,
		},
	}
	resp, err := client.Get("http://192.168.2.1/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
