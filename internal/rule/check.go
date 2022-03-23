package rule

import (
	"kms/pkg/misc/code"

	error2 "github.com/quanxiang-cloud/cabin/error"
)

// ********************************************************************

// IsActive verify active
func IsActive(active int) bool {
	return active == ActiveEnable
}

// ValidateActive ValidateActive
func ValidateActive(active int, op Operation) error {
	switch op {
	case OpUpdate, OpDelete:
		if IsActive(active) {
			return error2.New(code.ErrActiveKey)
		}
	case OpSignature:
		if !IsActive(active) {
			return error2.New(code.ErrNotActiveKey)
		}
	default:
	}
	return nil
}

// ********************************************************************

// CheckCharSet check data whether contains unsupport character
func CheckCharSet(data ...string) error {
	for _, v := range data {
		if !charSetExpr.Match([]byte(v)) {
			return error2.New(code.ErrCharacterSet)
		}
	}
	return nil
}

// ********************************************************************

// CheckParse check key status
func CheckParse(status int) bool {
	return status == Parsed
}
