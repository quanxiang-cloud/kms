package encrypt

import (
	"testing"
)

func TestAECBC(t *testing.T) {
	src := "123"
	key := "123456789abcdefg"
	dst, err := AESEncryptCBC([]byte(src), []byte(key), PKCS5)
	if err != nil {
		t.Fatal(err)
	}
	dst2, err := AESDecryptCBC(dst, []byte(key), PKCS5)
	if err != nil {
		t.Fatal(err)
	}

	if src != string(dst2) {
		t.Fatal("aes code fail.")
	}
}
