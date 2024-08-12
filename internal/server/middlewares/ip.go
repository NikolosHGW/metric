package middlewares

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

type CheckIP struct {
	logger        customLogger
	trustedSubnet string
}

func NewCheckIP(trustedSubnet string, logger customLogger) *CheckIP {
	return &CheckIP{
		trustedSubnet: trustedSubnet,
		logger:        logger,
	}
}

func (m *CheckIP) isIPTrusted(r *http.Request) bool {
	if m.trustedSubnet == "" {
		return true
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP == "" {
		return false
	}

	ip := net.ParseIP(xRealIP)
	_, cidr, err := net.ParseCIDR(m.trustedSubnet)
	if err != nil {
		m.logger.Info("Invalid CIDR", zap.Error(err))
		return false
	}

	return cidr.Contains(ip)
}

func (m *CheckIP) WithCheckIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.isIPTrusted(r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
