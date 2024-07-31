package middlewares

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/NikolosHGW/metric/internal/crypto"
	"go.uber.org/zap"
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

func (dm *DecryptMiddleware) DecryptHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dm.privateKeyPath == "" {
			next.ServeHTTP(w, r)
			return
		}
		privateKey, err := crypto.LoadPrivateKey(dm.privateKeyPath)
		if err != nil {
			dm.logger.Info("failed LoadPrivateKey", zap.Error(err))
			http.Error(w, "cannot load private key", http.StatusInternalServerError)
			return
		}
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

		if len(encryptedData) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Извлекаем зашифрованный симметричный ключ
		encryptedKeySize := privateKey.Size()
		if len(encryptedData) < encryptedKeySize {
			dm.logger.Info("encrypted data is too short")
			http.Error(w, "encrypted data is too short", http.StatusInternalServerError)
			return
		}

		encryptedKey := encryptedData[:encryptedKeySize]
		ciphertext := encryptedData[encryptedKeySize:]

		// Расшифровка симметричного ключа с использованием RSA
		aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedKey, nil)
		if err != nil {
			dm.logger.Info("failed to decrypt AES key", zap.Error(err))
			http.Error(w, "failed to decrypt AES key", http.StatusInternalServerError)
			return
		}

		// Расшифровка данных с использованием AES
		block, err := aes.NewCipher(aesKey)
		if err != nil {
			dm.logger.Info("failed to create AES cipher", zap.Error(err))
			http.Error(w, "failed to create AES cipher", http.StatusInternalServerError)
			return
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			dm.logger.Info("failed to create GCM", zap.Error(err))
			http.Error(w, "failed to create GCM", http.StatusInternalServerError)
			return
		}

		nonceSize := gcm.NonceSize()
		if len(ciphertext) < nonceSize {
			dm.logger.Info("ciphertext too short")
			http.Error(w, "ciphertext too short", http.StatusInternalServerError)
			return
		}

		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			dm.logger.Info("failed to decrypt data", zap.Error(err))
			http.Error(w, "failed to decrypt data", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(plaintext))
		r.ContentLength = int64(len(plaintext))

		next.ServeHTTP(w, r)
	})
}
