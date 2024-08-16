package interceptor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type HashMiddleware struct {
	key string
}

func NewHashMiddleware(key string) *HashMiddleware {
	return &HashMiddleware{
		key: key,
	}
}

func (hm *HashMiddleware) UnaryHashInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if hashValues := md.Get("hashsha256"); len(hashValues) > 0 && hm.key != "" {
		reqBytes, err := serializeRequest(req)
		if err != nil {
			return nil, status.Errorf(status.Code(err), "failed to serialize request: %v", err)
		}

		requestHash := hashValues[0]
		if !checkHash(reqBytes, hm.key, requestHash) {
			return nil, status.Errorf(status.Code(nil), "hash mismatch")
		}
	}

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	if hm.key != "" {
		respBytes, err := serializeResponse(resp)
		if err != nil {
			return nil, status.Errorf(status.Code(err), "failed to serialize response: %v", err)
		}
		respHash := getHash(respBytes, hm.key)
		newMD := metadata.Pairs("hashsha256", respHash)
		err = grpc.SetTrailer(ctx, newMD)
		if err != nil {
			return nil, status.Errorf(status.Code(err), "failed to SetTrailer: %v", err)
		}
	}

	return resp, nil
}

func checkHash(data []byte, key string, requestHash string) bool {
	return hmac.Equal([]byte(getHash(data, key)), []byte(requestHash))
}

func getHash(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func serializeRequest(req interface{}) ([]byte, error) {
	pb, ok := req.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to convert request to proto.Message")
	}

	return proto.Marshal(pb)
}

func serializeResponse(resp interface{}) ([]byte, error) {
	pb, ok := resp.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to convert response to proto.Message")
	}

	return proto.Marshal(pb)
}
