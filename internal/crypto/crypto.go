package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type customLogger interface {
	Info(string, ...zap.Field)
}

func GenerateCrypto(l customLogger, serverPrivateKeyPath, agentPublicKeyPath string) error {
	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		l.Info("cannot generate rsa key", zap.Error(err))

		return fmt.Errorf("cannot generate rsa key: %w", err)
	}

	// Кодирование приватного ключа в PEM формат
	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		l.Info("cannot encode private key", zap.Error(err))

		return fmt.Errorf("cannot encode private key: %w", err)
	}

	// Сохранение приватного ключа в файл
	err = os.WriteFile(serverPrivateKeyPath, privateKeyPEM.Bytes(), 0600)
	if err != nil {
		l.Info("cannot write private key to file", zap.Error(err))

		return fmt.Errorf("cannot write private key to file: %w", err)
	}
	l.Info("private key saved successfully", zap.String("path", serverPrivateKeyPath))

	// Генерация публичного ключа
	publicKey := &privateKey.PublicKey

	// Кодирование публичного ключа в PEM формат
	var publicKeyPEM bytes.Buffer
	err = pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})
	if err != nil {
		l.Info("cannot encode public key", zap.Error(err))

		return fmt.Errorf("cannot encode public key: %w", err)
	}

	// Сохранение публичного ключа в файл
	err = os.WriteFile(agentPublicKeyPath, publicKeyPEM.Bytes(), 0644)
	if err != nil {
		l.Info("cannot write public key to file", zap.Error(err))

		return fmt.Errorf("cannot write public key to file: %w", err)
	}
	l.Info("public key saved successfully", zap.String("path", agentPublicKeyPath))

	return nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func EncryptData(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}
