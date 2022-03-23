package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

// Sha256 Sha256
func Sha256(entity []byte, key []byte) ([]byte, error) {
	return sha(entity, hmac.New(sha256.New, key))
}

// Sha1 Sha1
func Sha1(entity []byte, key []byte) ([]byte, error) {
	return sha(entity, hmac.New(sha1.New, key))

}

func sha(entity []byte, hash hash.Hash) ([]byte, error) {
	_, err := hash.Write(entity)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
