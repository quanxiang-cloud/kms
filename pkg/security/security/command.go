package security

import (
	"errors"
)

// CMD cmd
type CMD struct {
	Name        string
	CommandProc Command
}

type arg []string

// Ans ans 结果集
type Ans []byte

func (a Ans) String() string {
	return string(a)
}

// Byte byte
func (a Ans) Byte() []byte {
	return a
}

var (
	// ErrUnSupportToByte not byte
	ErrUnSupportToByte = errors.New("type can not convert to byte")
)

// Type Type
type Type interface{}

// Byte byte
func Byte(t Type) ([]byte, error) {
	switch t.(type) {
	case []byte:
		return t.([]byte), nil
	case Ans:
		return t.(Ans), nil
	}

	return nil, ErrUnSupportToByte
}

// Command 命令模板函数
type Command func(Type, ...string) (Type, error)

var cmdTable = []*CMD{
	{Name: "none", CommandProc: none},

	{Name: "base64", CommandProc: base64},
	{Name: "hex", CommandProc: hex},
	{Name: "sha1", CommandProc: sha1},
	{Name: "sha256", CommandProc: sha256},
	{Name: "append", CommandProc: appendString},
	{Name: "url", CommandProc: urlCode},

	{Name: "md5", CommandProc: md5},
	{Name: "crc32", CommandProc: crc32},
	{Name: "crc64", CommandProc: crc64},

	{Name: "aes", CommandProc: aes},
	{Name: "rsa", CommandProc: rsa},

	{Name: "sort", CommandProc: sorting},
}

type commandTable map[string]*CMD

func populateCommandTable() commandTable {
	commands := make(map[string]*CMD, len(cmdTable))
	for _, command := range cmdTable {
		commands[command.Name] = command
	}

	return commands
}
