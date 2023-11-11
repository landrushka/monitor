package workers

import (
	"bytes"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"github.com/landrushka/monitor.git/internal/archiver"
	"github.com/landrushka/monitor.git/internal/storage"
	"log"
	"strings"
	"time"
)

func StartAgent(host string, reportInterval int64, pollInterval int64) error {
	c := resty.New()
	sf := storage.GaugeMetric{}
	si := storage.CounterMetric{}
	count := int64(0)
	var buf bytes.Buffer
	m := storage.Metrics{}

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		sf.GetRandomMetrics()
		sf.GetGaugeMetrics()
		si.GetCount()
		count += pollInterval
		if count/reportInterval >= 1 {
			count -= reportInterval
			if !strings.Contains(host, "http://") {
				host = "http://" + host
			}
			for name, value := range sf {
				m.ID = name
				m.MType = "gauge"
				m.Value = &value
				json.NewEncoder(&buf).Encode(m)
				b, _ := archiver.Compress(buf.Bytes())
				err := retry.Do(
					func() error {
						var err error
						_, err = c.R().SetHeader("Content-Type", "application/json").
							SetHeader("Content-Encoding", "gzip").
							SetBody(b).
							Post(host + "/update")
						return err
					},
					retry.Attempts(10),
					retry.OnRetry(func(n uint, err error) {
						log.Printf("Retrying request after error: %v", err)
					}),
				)
				buf.Reset()
				if err != nil {
					return err
				}
			}
			for name, value := range si {
				if !strings.Contains(host, "http://") {
					host = "http://" + host
				}
				m.ID = name
				m.MType = "counter"
				m.Delta = &value
				json.NewEncoder(&buf).Encode(m)
				b, _ := archiver.Compress(buf.Bytes())
				err := retry.Do(
					func() error {
						var err error
						_, err = c.R().SetHeader("Content-Type", "application/json").
							SetHeader("Content-Encoding", "gzip").
							SetBody(b).
							Post(host + "/update")
						return err
					},
					retry.Attempts(3),
					retry.OnRetry(func(n uint, err error) {
						log.Printf("Retrying request after error: %v", err)
					}),
				)
				buf.Reset()
				if err != nil {
					return err
				}
			}
		}

	}
}
