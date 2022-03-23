package security

import (
	"errors"
	"fmt"
	CRC32 "hash/crc32"
	CRC64 "hash/crc64"
	"kms/pkg/security/code"
	"kms/pkg/security/encrypt"
	"kms/pkg/security/sign"
	"kms/pkg/security/sort"
	"reflect"
	"strings"
)

// err 错误处理
var (
	ErrInvalidArg = errors.New("invalid args")
	ErrInvalidCMD = errors.New("invalid cmd")
)

const (
	encode = "encode"
	decode = "decode"
)

func none(entity Type, args ...string) (Type, error) {
	return nil, ErrInvalidCMD
}
func base64(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 2 {
		return nil, ErrInvalidArg
	}
	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	_t := strings.ToLower(args[0])
	_d := strings.ToLower(args[1])
	if _t == "std" {
		if _d == encode {
			return code.Base64StdEncode(value)
		} else if _d == decode {
			return code.Base64StdDecode(value)
		}
	} else if _t == "url" {
		if _d == encode {
			return code.Base64URLEncode(value)
		} else if _d == decode {
			return code.Base64URLDecode(value)
		}
	}
	return nil, ErrInvalidArg
}

func sha256(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 1 {
		return nil, ErrInvalidArg
	}

	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	return code.Sha256(value, []byte(args[0]))
}

func sha1(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 1 {
		return nil, ErrInvalidArg
	}

	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	return code.Sha1(value, []byte(args[0]))
}

func urlCode(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 1 {
		return nil, ErrInvalidArg
	}

	switch entity.(type) {
	case []byte:
		return code.URL(string(entity.([]byte)), args[0])
	case Ans:
		return code.URL(string(entity.(Ans)), args[0])
	case string:
		return code.URL(entity.(string), args[0])
	default:
		if reflect.TypeOf(entity).Kind() == reflect.Ptr {
			v := reflect.ValueOf(entity).Elem()
			switch v.Type().Kind() {
			case reflect.Struct:
				for i := 0; i < v.NumField(); i++ {
					field := v.Type().Field(i)
					if field.Type.Kind() == reflect.String && v.FieldByName(field.Name).CanSet() {
						code, _ := code.URL(v.Field(i).String(), args[0])
						v.FieldByName(field.Name).SetString(string(code))
					}
				}
			}

			return entity, nil
		}
	}
	return nil, ErrUnSupportToByte
}

func hex(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 1 {
		return nil, ErrInvalidArg
	}
	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	_d := strings.ToLower(args[0])
	if _d == "encode" {
		return code.HexEncode(value)
	} else if _d == "decode" {
		return code.HexDecode(value)
	}

	return nil, ErrInvalidArg
}

func appendString(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 2 {
		return nil, ErrInvalidArg
	}
	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	_d := strings.ToLower(args[0])
	if _d == "begin" {
		return code.AppendBegin(value, []byte(args[1]))
	} else if _d == "end" {
		return code.AppendEnd(value, []byte(args[1]))
	}

	return nil, ErrInvalidArg
}

func md5(entity Type, args ...string) (Type, error) {
	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	size := len(args)
	if size == 0 {
		return sign.MD5(value)
	} else if size == 1 {
		return sign.MD5WithSalt(value, []byte(args[0]))
	}

	return nil, ErrInvalidArg
}

func crc32(entity Type, args ...string) (Type, error) {
	size := len(args)
	var (
		temporary uint32
		err       error
	)

	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}
	switch {
	case size <= 0:
		temporary, err = sign.CRC32(value, 0)
	case strings.ToUpper(args[0]) == "IEEE":
		temporary, err = sign.CRC32(value, CRC32.IEEE)
	case strings.ToLower(args[0]) == "CASTAGNOLI":
		temporary, err = sign.CRC32(value, CRC32.Castagnoli)
	default:
		return nil, ErrInvalidArg
	}

	if err != nil {
		return nil, err
	}

	return []byte(uint32ToString(temporary)), nil
}

func crc64(entity Type, args ...string) (Type, error) {
	size := len(args)
	var (
		temporary uint64
		err       error
	)

	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}
	switch {
	case size <= 0:
		temporary, err = sign.CRC64(value, 0)
	case strings.ToUpper(args[0]) == "ISO":
		temporary, err = sign.CRC64(value, CRC64.ISO)
	case strings.ToLower(args[0]) == "ECMA":
		temporary, err = sign.CRC64(value, CRC64.ECMA)
	default:
		return nil, ErrInvalidArg
	}

	if err != nil {
		return nil, err
	}

	return []byte(uint64ToString(temporary)), nil
}

func uint32ToString(v uint32) string {
	return fmt.Sprintf("%d", v)
}

func uint64ToString(v uint64) string {
	return fmt.Sprintf("%d", v)
}

func aes(entity Type, args ...string) (Type, error) {
	size := len(args)
	if size < 1 {
		return nil, ErrInvalidArg
	}

	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}
	model := strings.ToUpper(args[0])
	if model == "CBC" {
		if size < 4 {
			return nil, ErrInvalidArg
		}

		return aesCBC(value, args[1:])
	}

	return nil, ErrInvalidArg
}

func aesCBC(entity Type, args []string) (Ans, error) {
	padding := strings.ToLower(args[0])
	_t := strings.ToLower(args[1])

	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	var paddingType encrypt.PaddingType
	if padding == "pkcs5" {
		paddingType = encrypt.PKCS5
	} else if padding == "pkcs7" {
		paddingType = encrypt.PKCS7
	} else {
		return nil, ErrInvalidArg
	}

	key := []byte(args[2])
	if _t == decode {
		return encrypt.AESDecryptCBC(value, key, paddingType)
	} else if _t == encode {
		return encrypt.AESEncryptCBC(value, key, paddingType)
	}

	return nil, ErrInvalidArg
}

func rsa(entity Type, args ...string) (Type, error) {
	value, err := Byte(entity)
	if err != nil {
		return nil, err
	}

	size := len(args)
	if size < 2 {
		return nil, ErrInvalidArg
	}
	_t := strings.ToLower(args[0])
	key := make([]byte, 0)
	for _, arg := range args[1:] {
		key = append(key, []byte(arg)...)
		key = append(key, []byte(" ")...)
	}

	if _t == decode {
		return encrypt.RSADecrypt(value, key)
	} else if _t == encode {
		return encrypt.RSAEncrypt(value, key)
	}

	return nil, ErrInvalidArg
}

func sorting(entity Type, args ...string) (Type, error) {
	var (
		format = sort.DefaultFormat
		order  = sort.DefaultSort
		mapper = sort.DefaultMapperType
	)

	for _, arg := range args {
		arg = strings.ToLower(arg)
		if f, ok := sort.OneofFormat(arg); ok {
			format = f
			continue
		}
		if m, ok := sort.OneofMapper(arg); ok {
			mapper = m
			continue
		}

		if o, ok := sort.OneofOrder(arg); ok {
			order = o
			continue
		}
	}
	if order == sort.SortDESC {
		return sort.WordSortDESC(entity, format, mapper)
	}
	return sort.WordSortASC(entity, format, mapper)
}
