package decode

import (
	"fmt"
	"kms/pkg/crypto/encode"
	"testing"
)

func TestCrypto(t *testing.T) {
	type testCase struct {
		src  string
		keys []string
	}
	var testCases = []*testCase{
		&testCase{
			src:  "hello, world",
			keys: []string{"foo", "bar"},
		},
		&testCase{
			src:  "hello, world",
			keys: []string{"foo", "bar"},
		},
		&testCase{
			src:  "hello, world",
			keys: []string{"123", "456"},
		},
		&testCase{
			src:  "hello, world2",
			keys: []string{"foo", "bar"},
		},
		&testCase{
			src:  "hello, world2",
			keys: []string{"123", "456"},
		},
	}
	for i, v := range testCases {
		c, err := encode.SecretEncodeString(v.src, v.keys...)
		if err != nil {
			t.Errorf("case %d SecretEncodeString error: %s", i+1, err.Error())
		}
		d, err := SecretDecodeString(c, v.keys...)
		if err != nil {
			t.Errorf("case %d SecretDecodeString error: %s", i+1, err.Error())
		}
		if d != v.src {
			t.Errorf("case %d SecretDecodeString mismatch: expect %s got %s", i+1, v.src, d)
		}

		if false {
			fmt.Printf("%d/%d %-14s %-24s %-16s\n", i+1, len(testCases), v.src, c, d)
		}
	}
}

func TestKey(t *testing.T) {
	type testCase struct {
		src    string
		keys   []string
		expect string
	}
	var testCases = []*testCase{
		&testCase{
			src:    "hello, world",
			keys:   []string{"foo", "bar"},
			expect: "AAAAAAAAAAAAAAAAAAAAAPo1uT6laOirkqJJwQ",
		},
		&testCase{
			src:    "hello, world",
			keys:   []string{"foo", "bar"},
			expect: "AAAAAAAAAAAAAAAAAAAAAPo1uT6laOirkqJJwQ",
		},
		&testCase{
			src:    "hello, world",
			keys:   []string{"123", "456"},
			expect: "AAAAAAAAAAAAAAAAAAAAACuAQpUZMra1zJnfqA",
		},
		&testCase{
			src:    "hello, world2",
			keys:   []string{"foo", "bar"},
			expect: "AAAAAAAAAAAAAAAAAAAAAPo1uT6laOirkqJJwd4",
		},
		&testCase{
			src:    "hello, world2",
			keys:   []string{"123", "456"},
			expect: "AAAAAAAAAAAAAAAAAAAAACuAQpUZMra1zJnfqDU",
		},
	}
	for i, v := range testCases {
		d, err := SecretDecodeString(v.expect, v.keys...)
		if d != v.src || err != nil {
			t.Errorf("case %d SecretDecodeString mismatch: expect %s got %s, err:%v", i+1, v.expect, d, err)
		}
	}
}
