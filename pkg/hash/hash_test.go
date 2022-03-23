package hash

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	h1 := Md5Hash(0, "x", "y")
	h2 := Md5Hash(0, "xy")
	fmt.Println("Md5Hash", h1)
	fmt.Println("Md5Hash", h2)
}
