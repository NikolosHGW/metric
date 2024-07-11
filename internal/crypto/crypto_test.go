package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"go.uber.org/zap"
)

type testLogger struct {
	logs []string
}

func (l *testLogger) Info(msg string, fields ...zap.Field) {
	l.logs = append(l.logs, msg)
}

func TestGenerateCrypto(t *testing.T) {
	// Создание временных файлов для хранения ключей
	serverPrivateKeyFile, err := os.CreateTemp("", "server_private_key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file for server private key: %v", err)
	}
	defer os.Remove(serverPrivateKeyFile.Name())

	agentPublicKeyFile, err := os.CreateTemp("", "agent_public_key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file for agent public key: %v", err)
	}
	defer os.Remove(agentPublicKeyFile.Name())

	// Создание фейкового логгера
	logger := &testLogger{}

	// Вызов тестируемой функции
	GenerateCrypto(logger, serverPrivateKeyFile.Name(), agentPublicKeyFile.Name())

	// Проверка логов на наличие ошибок
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

	// Проверка наличия и правильности сохранения приватного ключа
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

	// Проверка наличия и правильности сохранения публичного ключа
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

	// Сравнение публичного ключа с публичным ключом из приватного ключа
	if !publicKey.Equal(&privateKey.PublicKey) {
		t.Fatalf("Public key does not match private key's public key")
	}
}
