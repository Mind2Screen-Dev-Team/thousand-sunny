package xrsa

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey, xhash hash.Hash) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(xhash, nil, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey, xhash hash.Hash) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(xhash, nil, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func LoadPrivateKey(xpath string) (*rsa.PrivateKey, error) {
	f, err := os.Open(xpath)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	defer f.Close()
	defer b.Reset()

	if _, _err := f.WriteTo(&b); _err != nil && !errors.Is(_err, io.EOF) {
		return nil, _err
	}

	var (
		block, _  = pem.Decode(b.Bytes())
		key, _err = x509.ParsePKCS8PrivateKey(block.Bytes)
	)
	if _err != nil {
		return nil, _err
	}

	res, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("failed to load rsa private key")
	}

	return res, nil
}
