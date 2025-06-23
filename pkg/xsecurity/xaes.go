package xsecurity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type (
	EncryptionAES struct {
		Key []byte // Must be 32 bytes for AES-256
	}
)

const (
	nonceSize = 12
	tagSize   = 16
)

// NewEncryptionAES creates a new instance with the given key (must be 32 bytes)
func NewEncryptionAES(key []byte) (*EncryptionAES, error) {
	if len(key) != 32 {
		return nil, errors.New("AES-256 requires a 32-byte key")
	}
	return &EncryptionAES{Key: key}, nil
}

// Encrypt returns base64(nonce + tag + ciphertext)
func (e *EncryptionAES) Encrypt(plaintext string, key ...string) (string, error) {
	_key := e.Key
	if len(key) > 0 {
		decoded, err := base64.StdEncoding.DecodeString(key[0])
		if err != nil {
			return "", fmt.Errorf("%s:%w", "invalid aes key", err)
		}
		_key = decoded
	}

	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and get ciphertext with tag at the end
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Split tag from ciphertext
	tag := ciphertext[len(ciphertext)-tagSize:]
	body := ciphertext[:len(ciphertext)-tagSize]

	// Final format: nonce + tag + body
	final := append(nonce, append(tag, body...)...)
	return base64.StdEncoding.EncodeToString(final), nil
}

// Decrypt takes base64(nonce + tag + ciphertext) and returns plaintext
func (e *EncryptionAES) Decrypt(encoded string, key ...string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	if len(data) < nonceSize+tagSize {
		return "", errors.New("invalid encrypted data")
	}

	nonce := data[:nonceSize]
	tag := data[nonceSize : nonceSize+tagSize]
	body := data[nonceSize+tagSize:]

	_key := e.Key
	if len(key) > 0 {
		decKey, err := base64.StdEncoding.DecodeString(key[0])
		if err != nil {
			return "", fmt.Errorf("%s:%w", "invalid aes key", err)
		}
		_key = decKey
	}

	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Rebuild ciphertext (GCM expects tag at the end)
	fullCiphertext := append(body, tag...)

	plaintext, err := gcm.Open(nil, nonce, fullCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
