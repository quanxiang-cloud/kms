package code

import (
	"bytes"
)

// AppendEnd append end
func AppendEnd(src []byte, additional []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(src)
	buf.Write(additional)
	return buf.Bytes(), nil
}

// AppendBegin append end
func AppendBegin(src []byte, additional []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(additional)
	buf.Write(src)
	return buf.Bytes(), nil
}
