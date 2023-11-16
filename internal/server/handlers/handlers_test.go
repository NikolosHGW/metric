package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/util"
	"github.com/stretchr/testify/assert"
)

type StorageMock struct{}

func (sm StorageMock) GetGaugeMetric(name string) (util.Gauge, error) {
	return 502.12, nil
}

func (sm StorageMock) GetCounterMetric(name string) (util.Counter, error) {
	return 502, nil
}

func (sm StorageMock) SetGaugeMetric(name string, value util.Gauge) {

}

func (sm StorageMock) SetCounterMetric(name string, value util.Counter) {

}

func TestPostHandle(t *testing.T) {
	type want struct {
		code         int
		contentTypes []string
	}

	strg := StorageMock{}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "положительный тест для метрики типа counter",
			url:  "/update/counter/someMetric/527",
			want: want{
				code:         200,
				contentTypes: []string{"text/plain", "charset=utf-8"},
			},
		},
		{
			name: "положительный тест для метрики типа gauge",
			url:  "/update/gauge/someMetric/527.2",
			want: want{
				code:         200,
				contentTypes: []string{"text/plain", "charset=utf-8"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()

			WithSetMetricHandle(strg)(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)

			assert.ElementsMatch(t, test.want.contentTypes, res.Header.Values("Content-Type"))
		})
	}
}
