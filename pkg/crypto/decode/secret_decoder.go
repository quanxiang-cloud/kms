// NOTE: **********************NEVER MODIFY THIS FILE***************************

package decode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"reflect"
	"unsafe"
)

var (
	saltP          = []byte("#jK‥y8&%d\\nf%$ǚ$*)dK8^%3☒3 \t?2i\n/L-<9U")
	saltS          = []byte("UyH-j*<&^-o‥p*6ǘ5%/HJ^k,k86\t\n4j ǎguI$#jg\\nhj@!jkJhl")
	base64Encoding = base64.RawURLEncoding
)

const (
	ivSize = aes.BlockSize
)

// SecretDecodeString is decode for secret
func SecretDecodeString(encryptText string, keys ...string) (string, error) {
	iv, data, err := toBuffer(encryptText)
	if err != nil {
		return "", err
	}

	c, err := newCipher(keys, iv)
	if err != nil {
		return "", err
	}

	c.XORKeyStream(data, data)
	return string(data), nil
}

func toBuffer(encryptText string) ([]byte, []byte, error) {
	b, err := base64Encoding.DecodeString(encryptText)
	if err != nil {
		return nil, nil, err
	}
	if len(b) < ivSize {
		return nil, nil, errors.New("encryptText size to short")
	}
	iv, data := b[:ivSize], b[ivSize:]
	return iv, data, nil
}

func newCipher(keys []string, iv []byte) (cipher.Stream, error) {
	aesKey := hashKey(newBuf(), keys...)
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	cc := cipher.NewCTR(c, iv)
	return cc, nil
}

func newBuf() []byte {
	return make([]byte, 0, sha256.Size)
}

func hashKey(buf []byte, elems ...string) []byte {
	h := sha256.New()
	h.Write(saltP)
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))

		// non-stadard hash
		b := h.Sum(buf[:0])
		x := b[(b[0]>>3)&0x0F] ^ 0x5A
		h.Write([]byte{x})
	}
	h.Write(saltS)

	b := h.Sum(buf[:0])
	b[(b[0]>>3)&0x0F] ^= 0x5A // non-stadard hash
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
