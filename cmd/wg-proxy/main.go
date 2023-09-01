package main

import (
	"flag"
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
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

	var wg sync.WaitGroup

	if cfg.Proxy.HTTP.Addr != "" {
		wg.Add(1)
		proxy := goproxy.NewProxyHttpServer()
		proxy.Verbose = cfg.Debug
		proxy.Tr = &http.Transport{
			DialContext: dial.DialContext,
		}

		go func(wg *sync.WaitGroup, proxy *goproxy.ProxyHttpServer) {
			err := http.ListenAndServe(cfg.Proxy.HTTP.Addr, proxy)
			if err != nil {
				logrus.Error(err)
			}
			defer wg.Done()
		}(&wg, proxy)
	}

	wg.Wait()
}
