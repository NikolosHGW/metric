package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/metric/internal/server/handlers"
	"github.com/NikolosHGW/metric/internal/server/services"
	"github.com/NikolosHGW/metric/internal/server/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...zap.Field) {}

func TestWithGzip(t *testing.T) {
	strg := storage.NewMemStorage()
	metricService := services.NewMetricService(strg)
	handler := handlers.NewHandler(metricService, &mockLogger{})

	h := WithGzip(http.HandlerFunc(handler.SetJSONMetric))

	srv := httptest.NewServer(h)
	defer srv.Close()

	requestBody := `{"id":"foo","type":"gauge","value":42.1}`
	successBody := requestBody

	t.Run("положительный тест: application/json запрос без сжатия", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		buf.Write([]byte(requestBody))

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", applicationJSON)

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				t.Errorf("failed body close: %v", err)
			}
		}()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("положительный тест: application/json запрос, в котором сжатый JSON", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", applicationJSON)
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				t.Errorf("failed body close: %v", err)
			}
		}()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("положительный тест: application/json запрос без сжатия, ждущий в ответе сжатый JSON", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		buf.Write([]byte(requestBody))

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", applicationJSON)
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				t.Errorf("failed body close: %v", err)
			}
		}()

		gz, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer func() {
			err := gz.Close()
			if err != nil {
				t.Errorf("failed gzip close: %v", err)
			}
		}()

		b, err := io.ReadAll(gz)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})

	t.Run("положительный тест: application/json запрос, в котором сжатый JSON, ждущий в ответе сжатый JSON", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Type", applicationJSON)
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				t.Errorf("failed body close: %v", err)
			}
		}()

		gz, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer func() {
			err := gz.Close()
			if err != nil {
				t.Errorf("failed gzip close: %v", err)
			}
		}()

		b, err := io.ReadAll(gz)
		require.NoError(t, err)
		require.JSONEq(t, successBody, string(b))
	})
}
