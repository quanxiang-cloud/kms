package function

import "math/rand"

const (
	randomBasic = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)
const (
	idxBits = 6                 // 6 bits to represent randomBasic (62)
	idxMask = 1<<idxBits - 1    // 1-bits, mask-code
	idxMax  = idxMask / idxBits //
)

// String by random character
func String(l int) string {
	b := make([]byte, l)
	for i, cache, less := l-1, rand.Int63(), idxMax; i >= 0; {
		if less == 0 {
			less, cache = idxMax, rand.Int63()
		}
		if idx := int(cache & idxMask); idx < len(randomBasic) {
			b[i] = randomBasic[idx]
			i--
		}
		cache = cache >> idxBits
		less--
	}
	return string(b)
}
