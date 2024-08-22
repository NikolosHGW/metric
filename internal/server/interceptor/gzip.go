package interceptor

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"strings"

	"github.com/NikolosHGW/metric/internal/server/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryGzipInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if containsGzip(md["content-encoding"]) {
		decompressedReq, err := decompressRequest(req)
		if err != nil {
			logger.Log.Info("failed to decompress request", zap.Error(err))
			return nil, status.Errorf(status.Code(err), "failed to decompress request: %v", err)
		}
		req = decompressedReq
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	if containsGzip(md["accept-encoding"]) {
		compressedResp, err := compressResponse(resp)
		if err != nil {
			logger.Log.Info("failed to compress response", zap.Error(err))
			return nil, status.Errorf(status.Code(err), "failed to compress response: %v", err)
		}
		resp = compressedResp
	}

	return resp, nil
}

func containsGzip(encodings []string) bool {
	for _, encoding := range encodings {
		if strings.Contains(encoding, "gzip") {
			return true
		}
	}
	return false
}

func decompressRequest(req interface{}) (interface{}, error) {
	reqBytes, ok := req.([]byte)
	if !ok {
		return nil, status.Errorf(status.Code(nil), "request is not in bytes format")
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := gzipReader.Close()
		if err != nil {
			logger.Log.Info("err close gzipReader", zap.Error(err))
		}
	}()

	decompressedBytes, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	return decompressedBytes, nil
}

func compressResponse(resp interface{}) (interface{}, error) {
	respBytes, ok := resp.([]byte)
	if !ok {
		return nil, status.Errorf(status.Code(nil), "response is not in bytes format")
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	defer func() {
		err := gzipWriter.Close()
		if err != nil {
			logger.Log.Info("err close gzipWriter", zap.Error(err))
		}
	}()

	if _, err := gzipWriter.Write(respBytes); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
