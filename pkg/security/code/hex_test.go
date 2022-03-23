package code

import (
	"testing"
)

func TestHex(t *testing.T) {
	src := []byte{97, 244, 162, 200, 108, 176, 44, 73, 10, 198, 190, 13, 29, 232, 204, 195, 209, 2, 160, 140, 250, 159, 98, 161, 210, 195, 76, 111, 112, 17, 15, 35}
	code, _ := HexEncode([]byte(src))
	dst, err := HexDecode(code)
	if err != nil {
		t.Fatal(err)
	}

	if string(src) != string(dst) {
		t.Fatal("hex code fail.")
	}
}
