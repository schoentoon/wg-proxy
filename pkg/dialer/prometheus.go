package dialer

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var txBytesDesc = prometheus.NewDesc(
	"tx_bytes",
	"Transmitted bytes",
	[]string{"public_key"}, nil,
)
var rxBytesDesc = prometheus.NewDesc(
	"rx_bytes",
	"Received bytes",
	[]string{"public_key"}, nil,
)
var lastHandshake = prometheus.NewDesc(
	"last_handshake",
	"Last handshake in seconds",
	[]string{"public_key"}, nil,
)

func (d *Dialer) Describe(ch chan<- *prometheus.Desc) {
	ch <- txBytesDesc
	ch <- rxBytesDesc
	ch <- lastHandshake
}

func (d *Dialer) Collect(ch chan<- prometheus.Metric) {
	str, err := d.dev.IpcGet()
	if err != nil {
		logrus.Error(err)
		return
	}

	public_key := ""
	for _, line := range strings.Split(str, "\n") {
		if line == "" {
			return
		}

		split := strings.Split(line, "=")
		if len(split) != 2 {
			logrus.Errorf("Illegal data from ipc", line)
			return
		}

		key, value := split[0], split[1]

		switch key {
		case "public_key":
			decoded, err := hex.DecodeString(value)
			if err != nil {
				logrus.Error(err)
				return
			}
			public_key = base64.StdEncoding.EncodeToString(decoded)
		case "tx_bytes":
			bytes, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				logrus.Error(err)
			}
			ch <- prometheus.MustNewConstMetric(
				txBytesDesc,
				prometheus.GaugeValue,
				float64(bytes),
				public_key,
			)
		case "rx_bytes":
			bytes, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				logrus.Error(err)
			}
			ch <- prometheus.MustNewConstMetric(
				rxBytesDesc,
				prometheus.GaugeValue,
				float64(bytes),
				public_key,
			)
		case "last_handshake_time_sec":
			when, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				logrus.Error(err)
			}
			ch <- prometheus.MustNewConstMetric(
				lastHandshake,
				prometheus.GaugeValue,
				float64(when),
				public_key,
			)
		}
	}
}
