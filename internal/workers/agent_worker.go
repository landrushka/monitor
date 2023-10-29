package workers

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/landrushka/monitor.git/internal/metrics"
	"strings"
	"time"
)

func StartAgent(host string, reportInterval int64, pollInterval int64) error {
	c := resty.New()
	sf := metrics.StatsFloat{}
	si := metrics.StatsInt{}
	count := int64(0)
	var buf bytes.Buffer
	m := metrics.Metrics{}

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
				_, err := c.R().SetHeader("Content-Type", "application/json").
					SetBody(&buf).
					Post(host + "/update")
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
				_, err := c.R().SetHeader("Content-Type", "application/json").
					SetBody(&buf).
					Post(host + "/update")
				buf.Reset()
				if err != nil {
					return err
				}
			}
		}

	}
}
