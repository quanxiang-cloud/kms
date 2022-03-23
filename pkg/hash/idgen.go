package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"reflect"
	"unsafe"
)

// MaxHashConflict presents max time for hash conflict retry
const MaxHashConflict = 5

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Sha256Hash generate sha256 hash, if index>0, it try to avoid hash conflict by salt
func Sha256Hash(index int, elems ...string) string {
	h := sha256.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	if index > 0 {
		h.Write(unsafeStringBytes(fmt.Sprintf("salt_%d", index)))
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// Md5Hash generate sha256 hash, if index>0, it try to avoid hash conflict by salt
func Md5Hash(index int, elems ...string) string {
	h := md5.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	if index > 0 {
		h.Write(unsafeStringBytes(fmt.Sprintf("salt_%d", index)))
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
