package sign

import (
	"testing"
)

func TestMd5(t *testing.T) {
	dst, err := MD5([]byte("123"))
	if err != nil {
		t.Fatal(err)
	}
	_ = dst
}

func TestMd5WithSalt(t *testing.T) {
	dst, err := MD5WithSalt([]byte("123"), []byte("456"))
	if err != nil {
		t.Fatal(err)
	}
	_ = dst
}
