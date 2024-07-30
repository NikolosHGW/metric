package crypto

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"
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

func TestGenerateCrypto(t *testing.T) {
	serverPrivateKeyFile, err := os.CreateTemp("", "server_private_key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file for server private key: %v", err)
	}
	defer func() {
		err := os.Remove(serverPrivateKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file for server private key: %v", err)
		}
	}()

	agentPublicKeyFile, err := os.CreateTemp("", "agent_public_key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file for agent public key: %v", err)
	}
	defer func() {
		err := os.Remove(agentPublicKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file for agent public key: %v", err)
		}
	}()

	logger := &testLogger{}

	err = GenerateCrypto(logger, serverPrivateKeyFile.Name(), agentPublicKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed GenerateCrypto: %v", err)
	}

	for _, log := range logger.logs {
		cantGenerateKey := log == "cannot generate rsa key"
		cantEncodePrivKey := log == "cannot encode private key"
		cantWritePrivKeyToFile := log == "cannot write private key to file"
		cantEncodePubKey := log == "cannot encode public key"
		cantWritePubKeyToFile := log == "cannot write public key to file"

		if cantGenerateKey || cantEncodePrivKey || cantWritePrivKeyToFile || cantEncodePubKey || cantWritePubKeyToFile {
			t.Fatalf("Error occurred during GenerateCrypto execution: %s", log)
		}
	}

	privateKeyData, err := os.ReadFile(serverPrivateKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed to read server private key file: %v", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		t.Fatalf("Failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	publicKeyData, err := os.ReadFile(agentPublicKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed to read agent public key file: %v", err)
	}

	block, _ = pem.Decode(publicKeyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		t.Fatalf("Failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse public key: %v", err)
	}

	if !publicKey.Equal(&privateKey.PublicKey) {
		t.Fatalf("Public key does not match private key's public key")
	}
}

func TestGenerateCrypto_TableDriven(t *testing.T) {
	tests := []struct {
		name                 string
		serverPrivateKeyPath string
		agentPublicKeyPath   string
		expectedError        string
	}{
		{
			name:                 "Пустые пути",
			serverPrivateKeyPath: "",
			agentPublicKeyPath:   "",
			expectedError:        "cannot write private key to file",
		},
		{
			name:                 "Только приватный путь",
			serverPrivateKeyPath: "testfile",
			agentPublicKeyPath:   "",
			expectedError:        "cannot write public key to file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &testLogger{}
			err := GenerateCrypto(logger, tt.serverPrivateKeyPath, tt.agentPublicKeyPath)
			defer func() {
				if tt.serverPrivateKeyPath != "" {
					err := os.Remove(tt.serverPrivateKeyPath)
					if err != nil {
						t.Fatalf("Failed to remove temp file for server private key: %v", err)
					}
				}
			}()
			if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
				t.Fatalf("Expected error containing '%v', got '%v'", tt.expectedError, err)
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

	publicKeyPath := "test_public_key.pem"
	publicKeyFile, err := os.CreateTemp("", publicKeyPath)
	if err != nil {
		t.Fatalf("Failed to create temp file for public key: %v", err)
	}
	defer func() {
		err := os.Remove(publicKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file for public key: %v", err)
		}
	}()

	logger := &testLogger{}

	err = GenerateCrypto(logger, serverPrivateKeyFile.Name(), publicKeyFile.Name())
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
	privateKeyPath := "server_private_key.pem"
	publicKeyPath := "agent_public_key.pem"

	logger := &testLogger{}

	err := GenerateCrypto(logger, privateKeyPath, publicKeyPath)
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

	err = os.Remove(privateKeyPath)
	if err != nil {
		t.Fatalf("Cannot private remove file: %v", err)
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
	serverPrivateKeyFile, err := os.CreateTemp("", "server_private_key.pem")
	assert.NoError(t, err)
	defer func() {
		err := os.Remove(serverPrivateKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file serverPrivateKeyFile: %v", err)
		}
	}()

	agentPublicKeyFile, err := os.CreateTemp("", "agent_public_key.pem")
	assert.NoError(t, err)
	defer func() {
		err := os.Remove(agentPublicKeyFile.Name())
		if err != nil {
			t.Fatalf("Failed to remove temp file agentPublicKeyFile: %v", err)
		}
	}()

	logger := &testLogger{}
	err = GenerateCrypto(logger, serverPrivateKeyFile.Name(), agentPublicKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed GenerateCrypto: %v", err)
	}

	publicKey, err := LoadPublicKey(agentPublicKeyFile.Name())
	assert.NoError(t, err)

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
			wantErr: true,
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
				assert.NotNil(t, encryptedData)
			}
		})
	}
}
