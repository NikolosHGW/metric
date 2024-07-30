package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	return privateKey, &privateKey.PublicKey
}

func encryptMessage(t *testing.T, publicKey *rsa.PublicKey, message []byte) []byte {
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
	assert.NoError(t, err)
	return encryptedMessage
}

func createTempKeyFile(t *testing.T, keyData []byte) string {
	file, err := os.CreateTemp("", "testkey-*.pem")
	assert.NoError(t, err)

	_, err = file.Write(keyData)
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	return file.Name()
}

func TestDecryptMiddleware_Success(t *testing.T) {
	privateKey, publicKey := generateTestKeys(t)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	privateKeyPath := createTempKeyFile(t, privateKeyPEM)
	defer func() {
		err := os.Remove(privateKeyPath)
		if err != nil {
			t.Fatalf("Failed to remove temp file for server private key: %v", err)
		}
	}()

	logger := &mockLogger{}
	middleware := NewDecryptMiddleware(privateKeyPath, logger)

	message := []byte("hello, world")
	encryptedMessage := encryptMessage(t, publicKey, message)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encryptedMessage))
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, message, body)
	})

	middleware.DecryptHandler(nextHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDecryptMiddleware_Failure(t *testing.T) {
	privateKey, _ := generateTestKeys(t)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	privateKeyPath := createTempKeyFile(t, privateKeyPEM)
	defer func() {
		err := os.Remove(privateKeyPath)
		if err != nil {
			t.Fatalf("Failed to remove temp file for server private key: %v", err)
		}
	}()

	logger := &mockLogger{}
	middleware := NewDecryptMiddleware(privateKeyPath, logger)

	encryptedMessage := []byte("invalid data")

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encryptedMessage))
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	middleware.DecryptHandler(nextHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "failed to decrypt data")
}
