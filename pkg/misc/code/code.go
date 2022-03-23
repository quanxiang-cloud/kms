package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

// error code definition
const (
	ErrServiceWithKeys = 130010020010
	ErrInvalidAuthType = 130010020009
	ErrAuthContent     = 130010020008
	ErrCharacterSet    = 130010020007
	ErrKeyIsExists     = 130010020006
	ErrKeyDelFaild     = 130010020005
	ErrInvalidKey      = 130010020004
	ErrActiveKey       = 130010020003
	ErrNotActiveKey    = 130010020002
	ErrMaximumHold     = 130010020001

	ErrInvalidOAuth2ClientID = 130020020003
	ErrInvalidOAuth2State    = 130020020002
	ErrInvalidOAuth2Code     = 130020020001
)

// CodeTable codeTable
var CodeTable = map[int64]string{
	ErrServiceWithKeys: "该分组设置包含私钥",
	ErrInvalidAuthType: "不支持的鉴权方法，仅支持 %v",
	ErrAuthContent:     "鉴权方法配置不合法: %v",
	ErrCharacterSet:    "无效的字符",
	ErrMaximumHold:     "密钥持有数达到最大",
	ErrActiveKey:       "对象启用中，不能进行该操作",
	ErrNotActiveKey:    "对象未启用，不能进行该操作",
	ErrKeyDelFaild:     "failded to delete key", //?
	ErrInvalidKey:      "无效的密钥",
	ErrKeyIsExists:     "对象已存在",

	ErrInvalidOAuth2Code:     "invalid code, due to code is empty",
	ErrInvalidOAuth2ClientID: "invalid client id, due to client id is empty",
	ErrInvalidOAuth2State:    "invalid state, due to state is not defined",
}
