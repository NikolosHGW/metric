package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testLogger struct {
	logs []string
}

func (l *testLogger) Info(msg string, fields ...zap.Field) {
	l.logs = append(l.logs, msg)
}

func TestGeneratePrivateKey(t *testing.T) {
	tests := []struct {
		name                 string
		serverPrivateKeyPath string
		expectError          bool
	}{
		{
			name:                 "Successful key generation",
			serverPrivateKeyPath: "test_private_key.pem",
			expectError:          false,
		},
		{
			name:                 "Invalid file path",
			serverPrivateKeyPath: "",
			expectError:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &testLogger{}

			err := GeneratePrivateKey(l, tt.serverPrivateKeyPath)
			if (err != nil) != tt.expectError {
				t.Errorf("GeneratePrivateKey() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError {
				err = os.Remove(tt.serverPrivateKeyPath)
				if err != nil {
					t.Fatalf("Failed to remove serverPrivateKeyPath: %v", err)
				}
			}
		})
	}
}

func TestGeneratePublicKey(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	tests := []struct {
		name               string
		privateKey         *rsa.PrivateKey
		agentPublicKeyPath string
		expectError        bool
	}{
		{
			name:               "Successful key generation",
			privateKey:         privateKey,
			agentPublicKeyPath: "test_public_key.pem",
			expectError:        false,
		},
		{
			name:               "Invalid file path",
			privateKey:         privateKey,
			agentPublicKeyPath: "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GeneratePublicKey(tt.privateKey, tt.agentPublicKeyPath)
			if (err != nil) != tt.expectError {
				t.Errorf("GeneratePublicKey() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError {
				err = os.Remove(tt.agentPublicKeyPath)
				if err != nil {
					t.Fatalf("Failed to remove agentPublicKeyPath: %v", err)
				}
			}
		})
	}
}

func TestLoadPrivateKey_Success(t *testing.T) {
	privateKeyPath := "test_private_key.pem"
	serverPrivateKeyFile, err := os.CreateTemp("", privateKeyPath)
	if err != nil {
		t.Fatalf("Failed to create temp file for server private key: %v", err)
	}
	defer func() {
		err := os.Remove(serverPrivateKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file for server private key: %v", err)
		}
	}()

	logger := &testLogger{}

	err = GeneratePrivateKey(logger, serverPrivateKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed GenerateCrypto: %v", err)
	}

	privateKey, err := LoadPrivateKey(serverPrivateKeyFile.Name())

	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.IsType(t, &rsa.PrivateKey{}, privateKey)
}

func TestLoadPrivateKey_InvalidFormat(t *testing.T) {
	invalidKeyPath := "test_invalid_key.pem"
	serverPrivateKeyFile, err := os.CreateTemp("", invalidKeyPath)
	if err != nil {
		t.Fatalf("Failed to create temp file for server private key: %v", err)
	}
	defer func() {
		err := os.Remove(serverPrivateKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file for server private key: %v", err)
		}
	}()

	privateKey, err := LoadPrivateKey(serverPrivateKeyFile.Name())

	assert.Error(t, err)
	assert.Nil(t, privateKey)
	assert.EqualError(t, err, "failed to decode PEM block containing private key")
}

func TestLoadPrivateKey_FileNotFound(t *testing.T) {
	nonExistentPath := "non_existent_file.pem"

	privateKey, err := LoadPrivateKey(nonExistentPath)

	assert.Error(t, err)
	assert.Nil(t, privateKey)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestLoadPrivateKey_CorruptedFile(t *testing.T) {
	corruptedKeyPath := "test_corrupted_key.pem"
	err := os.WriteFile(corruptedKeyPath, []byte("invalid data"), 0600)
	assert.NoError(t, err)

	defer func() {
		err := os.Remove(corruptedKeyPath)
		if err != nil {
			t.Fatalf("Failed to remove temp file for server private key: %v", err)
		}
	}()

	privateKey, err := LoadPrivateKey(corruptedKeyPath)

	assert.Error(t, err)
	assert.Nil(t, privateKey)
}

func TestLoadPublicKey(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	publicKeyPath := "agent_public_key.pem"

	err := GeneratePublicKey(privateKey, publicKeyPath)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	publicKey, err := LoadPublicKey(publicKeyPath)
	if err != nil {
		t.Fatalf("Failed to load public key: %v", err)
	}

	if publicKey == nil {
		t.Fatal("Public key is nil")
	}

	err = os.Remove(publicKeyPath)
	if err != nil {
		t.Fatalf("Cannot public remove file: %v", err)
	}
}

func TestLoadPublicKey_FileNotFound(t *testing.T) {
	_, err := LoadPublicKey("non_existent_file.pem")
	if err == nil {
		t.Fatal("Expected an error for non-existent file, but got nil")
	}
}

func TestLoadPublicKey_InvalidPEM(t *testing.T) {
	invalidPEMPath := "invalid.pem"
	err := os.WriteFile(invalidPEMPath, []byte("invalid data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid PEM file: %v", err)
	}
	defer func() {
		err := os.Remove(invalidPEMPath)
		if err != nil {
			t.Fatalf("Failed to remove temp file invalidPEMPath: %v", err)
		}
	}()

	_, err = LoadPublicKey(invalidPEMPath)
	if err == nil {
		t.Fatal("Expected an error for invalid PEM data, but got nil")
	}
}

func TestLoadPublicKey_InvalidKeyType(t *testing.T) {
	invalidKeyTypePath := "invalid_key_type.pem"
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: []byte("invalid key data"),
	}
	var pemData bytes.Buffer
	err := pem.Encode(&pemData, block)
	if err != nil {
		t.Fatalf("Failed to encode invalid key type: %v", err)
	}
	err = os.WriteFile(invalidKeyTypePath, pemData.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid key type file: %v", err)
	}
	defer func() {
		err := os.Remove(invalidKeyTypePath)
		if err != nil {
			t.Fatalf("Failed to remove temp file invalidKeyTypePath: %v", err)
		}
	}()

	_, err = LoadPublicKey(invalidKeyTypePath)
	if err == nil {
		t.Fatal("Expected an error for invalid key type, but got nil")
	}
}

func TestLoadPublicKey_InvalidKeyFormat(t *testing.T) {
	invalidKeyFormatPath := "invalid_key_format.pem"
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: []byte("invalid key data"),
	}
	var pemData bytes.Buffer
	err := pem.Encode(&pemData, block)
	if err != nil {
		t.Fatalf("Failed to encode invalid key format: %v", err)
	}
	err = os.WriteFile(invalidKeyFormatPath, pemData.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid key format file: %v", err)
	}
	defer func() {
		err := os.Remove(invalidKeyFormatPath)
		if err != nil {
			t.Fatalf("Failed to remove temp file invalidKeyFormatPath: %v", err)
		}
	}()

	_, err = LoadPublicKey(invalidKeyFormatPath)
	if err == nil {
		t.Fatal("Expected an error for invalid key format, but got nil")
	}
}

func TestEncryptData(t *testing.T) {
	// Генерация ключей для тестов
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)

	publicKey := &privateKey.PublicKey

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "Валидные данные",
			data:    []byte("test data"),
			wantErr: false,
		},
		{
			name:    "Пустые данные",
			data:    []byte(""),
			wantErr: false,
		},
		{
			name:    "Большие данные",
			data:    make([]byte, 1000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedData, err := EncryptData(publicKey, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, encryptedData)
			} else {
				assert.NoError(t, err)
				if len(tt.data) == 0 {
					assert.Equal(t, []byte{}, encryptedData)
				} else {
					assert.NotNil(t, encryptedData)
					assert.Greater(t, len(encryptedData), 0)
				}
			}
		})
	}
}
