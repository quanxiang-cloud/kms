// NOTE: **********************NEVER MODIFY THIS FILE***************************
// crypto: sha256(orgKeys)+iv(1234567890abcdef) -> AES(CTR mode)

package netcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"reflect"
	"unsafe"
)

var (
	iv             = []byte("1234567890abcdef")
	base64Encoding = base64.RawURLEncoding
)

// EasyEncodeString is encode for secret
func EasyEncodeString(plantText string, keys ...string) (string, error) {
	c, err := newCipher(keys)
	if err != nil {
		return "", err
	}
	b := []byte(plantText)
	c.XORKeyStream(b, b)
	return base64Encoding.EncodeToString(b), nil
}

// EasyDecodeString is decode for secret
func EasyDecodeString(encryptText string, keys ...string) (string, error) {
	c, err := newCipher(keys)
	if err != nil {
		return "", err
	}
	b, err := base64Encoding.DecodeString(encryptText)
	if err != nil {
		return "", err
	}
	c.XORKeyStream(b, b)
	return string(b), nil
}

// EasyDecodeStringByAesKey is decode for secret
func EasyDecodeStringByAesKey(encryptText string, key []byte) (string, error) {
	c, err := newCipherByAesKey(key)
	if err != nil {
		return "", err
	}
	b, err := base64Encoding.DecodeString(encryptText)
	if err != nil {
		return "", err
	}
	c.XORKeyStream(b, b)
	return string(b), nil
}

// EasyKey generate aes key (hash256) from keys
func EasyKey(keys ...string) []byte {
	return hashKey(nil, keys...)
}

func newCipher(keys []string) (cipher.Stream, error) {
	key := hashKey(nil, keys...)
	c, err := newCipherByAesKey(key)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newCipherByAesKey(key []byte) (cipher.Stream, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cc := cipher.NewCTR(c, iv)
	return cc, nil
}

func hashKey(buf []byte, elems ...string) []byte {
	h := sha256.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	b := h.Sum(buf)
	return b
}

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func newBuf() []byte {
	return make([]byte, sha256.Size)
}
