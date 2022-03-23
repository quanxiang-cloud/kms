package apipath

import (
	"path"
	"strings"
)

// Split parse full path with path and name
func Split(full string) (string, string) {
	if !strings.HasPrefix(full, "/") {
		full = "/" + full
	}
	path, name := "", full
	if index := strings.LastIndex(full, "/"); index >= 0 {
		path = full[:index]
		name = full[index+1:]
	}
	return path, name
}

// Join join the namespace and name as full path
func Join(namespace, name string) string {
	if !strings.HasPrefix(namespace, "/") {
		namespace = "/" + namespace
	}
	return path.Join(namespace, name)
}
