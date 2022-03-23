package code

import (
	"net/url"
)

// URL url编码
func URL(entity string, t string) ([]byte, error) {
	if "path" == t {
		entity = url.PathEscape(string(entity))
	} else if "query" == t {
		entity = url.QueryEscape(string(entity))
	}

	return []byte(entity), nil
}
