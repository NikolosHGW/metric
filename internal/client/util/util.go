package util

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NikolosHGW/metric/internal/client/metrics"
	"github.com/NikolosHGW/metric/internal/util"
)

const (
	pollInterval   = 2
	reportInterval = 10
)

func CollectMetrics(m metrics.ClientMetrics) {
	for {
		m.RefreshMetrics()
		m.IncPollCount()
		m.UpdateRandomValue()

		time.Sleep(pollInterval * time.Second)
	}
}

func SendMetrics(m metrics.ClientMetrics) {
	metricTypeMap := util.GetMetricTypeMap()
	sb := strings.Builder{}
	for {
		// fmt.Println(m.GetMetrics())
		for k, v := range m.GetMetrics() {
			fmt.Println(k, v)
			gaugeValue, gaugeOk := v.(util.Gauge)
			counterValue, counterOk := v.(util.Counter)
			fmt.Println(gaugeOk, counterOk)
			if gaugeOk || counterOk {
				result := ""
				if gaugeOk {
					result = strconv.FormatFloat(float64(gaugeValue), 'f', -1, 64)
				} else {

					result = strconv.Itoa(int(counterValue))
				}
				sb.WriteString("http://localhost:8080/")
				sb.WriteString(metricTypeMap[k])
				sb.WriteString("/")
				sb.WriteString(k)
				sb.WriteString("/")
				sb.WriteString(result)

				resp, _ := http.Post(sb.String(), "application/json", nil)
				// if err != nil {
				// 	return err
				// }
				resp.Body.Close()
			}
		}

		time.Sleep(reportInterval * time.Second)
	}
}
