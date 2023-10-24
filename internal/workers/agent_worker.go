package workers

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/landrushka/monitor.git/internal/metrics"
	"strconv"
	"strings"
	"time"
)

func StartWorker(host string, reportInterval int64, pollInterval int64) {
	c := resty.New()
	sf := metrics.StatsFloat{}
	si := metrics.StatsInt{}
	count := int64(0)

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
				_, err := c.R().SetPathParams(map[string]string{
					"name":  name,
					"value": fmt.Sprintf("%.2f", value),
				}).SetHeader("Content-Type", "text/plain").
					Post(host + "/update/gauge/{name}/{value}")
				if err != nil {
					panic(err)
				}
			}
			for name, value := range si {
				if !strings.Contains(host, "http://") {
					host = "http://" + host
				}
				_, err := c.R().SetPathParams(map[string]string{
					"name":  name,
					"value": strconv.FormatInt(value, 10),
				}).SetHeader("Content-Type", "text/plain").
					Post(host + "/update/counter/{name}/{value}")
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
