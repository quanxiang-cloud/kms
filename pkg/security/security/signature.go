package security

import (
	"fmt"
	"strings"
)

// ConfigSecretKeyName is the key name in cmd config
const ConfigSecretKeyName = "<SECRET_KEY>"

// DefaultInst is the default instance of signature
var DefaultInst = New()

// Signature generate a signature of given data
func Signature(cmd string, d interface{}) (string, error) {
	r, err := DefaultInst.LookupCommandString(cmd, d)
	return r.String(), err
}

const (
	cmdSep   = "|"
	spaceSep = " "
)

// SignatureCMD is command of signature
type SignatureCMD struct {
	CMD string
	Arg []string
}

// SplitSignatureCMD split
func SplitSignatureCMD(cmd string) []SignatureCMD {
	commands := strings.Split(cmd, cmdSep)
	ans := make([]SignatureCMD, 0, len(commands))
	for _, command := range commands {
		params := strings.Split(command, spaceSep)

		size := len(params)
		if size == 0 {
			continue
		}

		cmd := SignatureCMD{
			CMD: params[0],
		}
		if size > 1 {
			cmd.Arg = params[1:]
		}

		ans = append(ans, cmd)
	}

	return ans
}

var (
	signatureCmds = map[string]interface{}{
		"none":   nil,
		"base64": nil,
		"hex":    nil,
		"sha1":   nil,
		"sha256": nil,
		"append": nil,
		"url":    nil,
		"md5":    nil,
		"crc32":  nil,
		"crc64":  nil,
		"aes":    nil,
		"rsa":    nil,
		"sort":   nil,
	}
)

// ValidSignatureCmds valid cmd
func ValidSignatureCmds(cmd string) error {
	if _, ok := signatureCmds[cmd]; !ok {
		return fmt.Errorf("invalid cmd %s", cmd)
	}
	return nil
}
