package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// PaddingType 追加类型
type PaddingType string

const (
	// PKCS5 pksc5
	PKCS5 PaddingType = "pkcs5"

	// PKCS7 pkcs7
	PKCS7 PaddingType = "pkcs7"
)

// AESEncryptCBC aes encrypt CBC
func AESEncryptCBC(origData []byte, key []byte, _t PaddingType) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = padding(_t)(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AESDecryptCBC aes decrypt CBC
func AESDecryptCBC(encrypted []byte, key []byte, _t PaddingType) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)
	decrypted = unPadding(_t)(decrypted)
	return decrypted, nil
}

func padding(_t PaddingType) func([]byte, int) []byte {
	switch _t {
	case PKCS5:
		return pkcs5Padding
	case PKCS7:
		return pkcs7Padding
	}
	return pkcs5Padding
}

func unPadding(_t PaddingType) func([]byte) []byte {
	switch _t {
	case PKCS5:
		return pkcs5UnPadding
	case PKCS7:
		return pkcs5UnPadding
	}
	return pkcs5UnPadding
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)

}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
