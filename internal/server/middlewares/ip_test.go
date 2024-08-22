package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckIPMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		trustedSubnet  string
		xRealIP        string
		expectedStatus int
	}{
		{
			name:           "No trusted subnet, any IP allowed",
			trustedSubnet:  "",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid IP in trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid IP outside trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.2.10",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid CIDR format",
			trustedSubnet:  "invalid-cidr",
			xRealIP:        "192.168.1.10",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Empty X-Real-IP header",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "IP is within the subnet boundary",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.255",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "IP is outside the subnet boundary",
			trustedSubnet:  "192.168.1.0/25",
			xRealIP:        "192.168.1.128",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := NewCheckIP(tt.trustedSubnet, &mockLogger{})

			handler := middleware.WithCheckIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
