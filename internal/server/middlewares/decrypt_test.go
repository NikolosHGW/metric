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

	"github.com/NikolosHGW/metric/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	return privateKey, &privateKey.PublicKey
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

func TestDecryptMiddleware(t *testing.T) {
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

	tests := []struct {
		name             string
		prepareRequest   func() *http.Request
		expectedBody     []byte
		expectedResponse int
		expectNextCalled bool
	}{
		{
			name: "Successful decryption",
			prepareRequest: func() *http.Request {
				message := []byte("hello, world")
				encryptedMessage, err := crypto.EncryptData(publicKey, message)
				assert.NoError(t, err)
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encryptedMessage))
			},
			expectedResponse: http.StatusOK,
			expectedBody:     []byte("hello, world"),
			expectNextCalled: true,
		},
		{
			name: "Invalid data",
			prepareRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid data")))
			},
			expectedResponse: http.StatusInternalServerError,
			expectedBody:     nil,
			expectNextCalled: false,
		},
		{
			name: "Empty data",
			prepareRequest: func() *http.Request {
				encryptedMessage, err := crypto.EncryptData(publicKey, []byte{})
				assert.NoError(t, err)
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encryptedMessage))
			},
			expectedResponse: http.StatusOK,
			expectedBody:     []byte{},
			expectNextCalled: true,
		},
		{
			name: "Too short data",
			prepareRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(make([]byte, 10)))
			},
			expectedResponse: http.StatusInternalServerError,
			expectedBody:     nil,
			expectNextCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.prepareRequest()
			rr := httptest.NewRecorder()

			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, body)
			})

			middleware.DecryptHandler(nextHandler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedResponse, rr.Code)
			if !tt.expectNextCalled {
				assert.False(t, nextCalled, "next handler should not be called")
			} else {
				assert.True(t, nextCalled, "next handler should be called")
			}
		})
	}
}
