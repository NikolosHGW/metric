package interceptor

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func (m *CheckIP) UnaryCheckIPInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if m.trustedSubnet == "" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.logger.Info("No metadata found in context")
		return nil, status.Errorf(status.Code(nil), "No metadata found in context")
	}

	var clientIP string
	if xRealIP := md.Get("x-real-ip"); len(xRealIP) > 0 {
		clientIP = xRealIP[0]
	} else if xForwardedFor := md.Get("x-forwarded-for"); len(xForwardedFor) > 0 {
		clientIP = xForwardedFor[0]
	} else {
		m.logger.Info("No client IP found in metadata")
		return nil, status.Errorf(status.Code(nil), "No client IP found in metadata")
	}

	ip := net.ParseIP(clientIP)
	if ip == nil {
		m.logger.Info("Invalid IP address", zap.String("clientIP", clientIP))
		return nil, status.Errorf(status.Code(nil), "Invalid IP address")
	}

	_, cidr, err := net.ParseCIDR(m.trustedSubnet)
	if err != nil {
		m.logger.Info("Invalid CIDR", zap.Error(err))
		return nil, status.Errorf(status.Code(err), "Invalid CIDR: %v", err)
	}

	if !cidr.Contains(ip) {
		m.logger.Info("IP address not trusted", zap.String("clientIP", clientIP))
		return nil, status.Errorf(status.Code(nil), "IP address not trusted")
	}

	return handler(ctx, req)
}
