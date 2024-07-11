package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"go.uber.org/zap"
)

type customLogger interface {
	Info(string, ...zap.Field)
}

func GenerateCrypto(l customLogger, serverPrivateKeyPath, agentPublicKeyPath string) {
	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		l.Info("cannot generate rsa key", zap.Error(err))
	}

	// Кодирование приватного ключа в PEM формат
	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		l.Info("cannot encode private key", zap.Error(err))
		return
	}

	// Сохранение приватного ключа в файл
	err = os.WriteFile(serverPrivateKeyPath, privateKeyPEM.Bytes(), 0600)
	if err != nil {
		l.Info("cannot write private key to file", zap.Error(err))
		return
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
		return
	}

	// Сохранение публичного ключа в файл
	err = os.WriteFile(agentPublicKeyPath, publicKeyPEM.Bytes(), 0644)
	if err != nil {
		l.Info("cannot write public key to file", zap.Error(err))
		return
	}
	l.Info("public key saved successfully", zap.String("path", agentPublicKeyPath))
}
