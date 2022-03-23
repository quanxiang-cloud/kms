package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RSAEncrypt rsa encrypt
func RSAEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	block, rest := pem.Decode(publicKey)
	if rest != nil && len(rest) != 0 {
		return nil, errors.New("pem decode public key err")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to parse public key (use ParsePKIXPublicKey instead for this key format")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, key, origData)
}

// RSADecrypt rsa decrypt
func RSADecrypt(cipherText []byte, privateKey []byte) ([]byte, error) {
	block, rest := pem.Decode(privateKey)
	if rest != nil && len(rest) != 0 {
		return nil, errors.New("pem decode privite key err")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pri, cipherText)
}
