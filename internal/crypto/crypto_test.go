package crypto

import (
	"crypto/rsa"
	"crypto/x509"
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

	GenerateCrypto(logger, serverPrivateKeyFile.Name(), agentPublicKeyFile.Name())

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

	GenerateCrypto(logger, serverPrivateKeyFile.Name(), "")

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
