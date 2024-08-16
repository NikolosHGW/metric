package interceptor

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"

	"github.com/NikolosHGW/metric/internal/crypto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type customLogger interface {
	Info(string, ...zap.Field)
}

type DecryptMiddleware struct {
	logger         customLogger
	privateKeyPath string
}

func NewDecryptMiddleware(privateKeyPath string, logger customLogger) *DecryptMiddleware {
	return &DecryptMiddleware{
		privateKeyPath: privateKeyPath,
		logger:         logger,
	}
}

func (dm *DecryptMiddleware) UnaryDecryptInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if dm.privateKeyPath == "" {
		return handler(ctx, req)
	}

	privateKey, err := crypto.LoadPrivateKey(dm.privateKeyPath)
	if err != nil {
		dm.logger.Info("failed to load private key", zap.Error(err))
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	reqBytes, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		dm.logger.Info("failed to marshal request", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	encryptedKeySize := privateKey.Size()
	if len(reqBytes) < encryptedKeySize {
		dm.logger.Info("encrypted data is too short")
		return nil, fmt.Errorf("encrypted data is too short")
	}

	encryptedKey := reqBytes[:encryptedKeySize]
	ciphertext := reqBytes[encryptedKeySize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedKey, nil)
	if err != nil {
		dm.logger.Info("failed to decrypt AES key", zap.Error(err))
		return nil, fmt.Errorf("failed to decrypt AES key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		dm.logger.Info("failed to create AES cipher", zap.Error(err))
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		dm.logger.Info("failed to create GCM", zap.Error(err))
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		dm.logger.Info("ciphertext too short")
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		dm.logger.Info("failed to decrypt data", zap.Error(err))
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	decryptedReq := proto.Clone(req.(proto.Message))
	err = proto.Unmarshal(plaintext, decryptedReq)
	if err != nil {
		dm.logger.Info("failed to unmarshal decrypted data", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal decrypted data: %w", err)
	}

	return handler(ctx, decryptedReq)
}
