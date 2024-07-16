package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/NikolosHGW/metric/internal/crypto"
	"go.uber.org/zap"
)

type customLogger interface {
	Info(string, ...zap.Field)
}

type DecryptMiddleware struct {
	privateKey *rsa.PrivateKey
	logger     customLogger
}

func NewDecryptMiddleware(privateKeyPath string, logger customLogger) (*DecryptMiddleware, error) {
	privateKey, err := crypto.LoadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return &DecryptMiddleware{
		privateKey: privateKey,
		logger:     logger,
	}, nil
}

func (dm *DecryptMiddleware) DecryptHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			dm.logger.Info("failed to read request body", zap.Error(err))
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				dm.logger.Info("err close body", zap.Error(err))
			}
		}()

		decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, dm.privateKey, encryptedData)
		if err != nil {
			dm.logger.Info("failed to decrypt data", zap.Error(err))
			http.Error(w, "failed to decrypt data", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(decryptedData))
		r.ContentLength = int64(len(decryptedData))

		next.ServeHTTP(w, r)
	})
}
