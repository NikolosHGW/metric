package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

type customLogger interface {
	Info(string, ...zap.Field)
}

func GeneratePrivateKey(l customLogger, serverPrivateKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		l.Info("cannot generate rsa key", zap.Error(err))

		return fmt.Errorf("cannot generate rsa key: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		l.Info("cannot encode private key", zap.Error(err))

		return fmt.Errorf("cannot encode private key: %w", err)
	}

	err = os.WriteFile(serverPrivateKeyPath, privateKeyPEM.Bytes(), 0600)
	if err != nil {
		l.Info("cannot write private key to file", zap.Error(err))

		return fmt.Errorf("cannot write private key to file: %w", err)
	}
	l.Info("private key saved successfully", zap.String("path", serverPrivateKeyPath))

	return nil
}

func GeneratePublicKey(privateKey *rsa.PrivateKey, agentPublicKeyPath string) error {
	publicKey := &privateKey.PublicKey

	var publicKeyPEM bytes.Buffer
	err := pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})
	if err != nil {
		log.Printf("cannot encode public key: %v", err)

		return fmt.Errorf("cannot encode public key: %w", err)
	}

	err = os.WriteFile(agentPublicKeyPath, publicKeyPEM.Bytes(), 0644)
	if err != nil {
		log.Printf("cannot write public key to file: %v", err)

		return fmt.Errorf("cannot write public key to file: %w", err)
	}
	log.Printf("public key saved successfully, path: %v", agentPublicKeyPath)

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
	if len(data) == 0 {
		return []byte{}, nil
	}

	// Генерация нового симметричного ключа AES-256
	aesKey := make([]byte, 32) // 32 байта для AES-256
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("cannot generate AES key: %w", err)
	}

	// Шифрование данных с использованием AES
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("cannot generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Шифрование симметричного ключа с использованием RSA
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt AES key: %w", err)
	}

	// Возвращаем зашифрованный симметричный ключ и зашифрованные данные
	result := append(encryptedKey, ciphertext...)
	return result, nil
}
