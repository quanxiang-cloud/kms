package signature

import (
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {
	demo := map[string]interface{}{
		"id":     "1",
		"name":   "张三",
		"age":    18,
		"status": true,
		"address": map[string]interface{}{
			"contry":   "china",
			"province": "sichuang",
			"city":     "chengdu",
		},
		"interest":  []string{"basketball", "football", "pingpong"},
		"timestamp": "2021-01-29T17:23:05-0800",
	}
	res, err := Signature(demo, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
	_ = res
}
