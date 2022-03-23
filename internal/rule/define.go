package rule

import "regexp"

// active value
const (
	ActiveAny     = -1
	ActiveEnable  = 1
	ActiveDisable = 0
	ActiveDefault = ActiveEnable
)

// value of parse
const (
	Parsed    = 1
	NotParsed = 0
)

// Operation is enum of operation
type Operation uint

// Opeartions
const (
	OpCreate Operation = iota + 1
	OpUpdate
	OpQuery
	OpDelete
	OpSignature
)

var (
	charSetExpr = regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9,./;'\\[\\]\\\\<>?:\"{}|`~!@#$%^&*()_+-=\\s\\n，。；‘’【】、《》？：“”{}·~！￥…（）]*$")
)
