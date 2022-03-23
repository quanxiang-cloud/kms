package code

import (
	"testing"
)

func TestBase64Std(t *testing.T) {
	src := "123456"
	code, _ := Base64StdEncode([]byte(src))
	dst, err := Base64StdDecode(code)
	if err != nil {
		t.Fatal(err)
	}

	if src != string(dst) {
		t.Fatal("base64 code fail.")
	}
}

func TestBase64URL(t *testing.T) {
	src := "123"
	code, _ := Base64URLEncode([]byte(src))
	dst, err := Base64URLDecode(code)
	if err != nil {
		t.Fatal(err)
	}

	if src != string(dst) {
		t.Fatal("base64 code fail.")
	}
}
