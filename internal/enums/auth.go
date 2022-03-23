package enums

// Auth enums describe the type of authentication
var (
	AuthTypeSet   = NewEnumSet(nil)
	AuthNone      = AuthTypeSet.Reg("none")
	AuthSystem    = AuthTypeSet.Reg("system")
	AuthSignature = AuthTypeSet.Reg("signature")
	AuthCookie    = AuthTypeSet.Reg("cookie")
	// AuthOAuth2 = AuthTypeSet.Reg("oauth2")
)

// basic value type
var (
	BasicTypeEnumSet   = NewEnumSet(nil)
	BasicTypeString    = BasicTypeEnumSet.Reg("string")
	BasicTypeNumber    = BasicTypeEnumSet.Reg("number")
	BasicTypeBoolean   = BasicTypeEnumSet.Reg("boolean")
	BasicTypeKeyID     = BasicTypeEnumSet.Reg("keyid")
	BasicTypeKeySecret = BasicTypeEnumSet.Reg("keysecret")
	BasicTypeSignature = BasicTypeEnumSet.Reg("signature")
	BasicTypeExpire    = BasicTypeEnumSet.Reg("expire")
	BasicTypeMethod    = BasicTypeEnumSet.Reg("method")
	BasicTypeAuthURL   = BasicTypeEnumSet.Reg("authurl")
	BasicTypeCookie    = BasicTypeEnumSet.Reg("cookie")
	BasicTypeSignCmd   = BasicTypeEnumSet.Reg("signcmd")
	// BasicTypeTimestamp  = BasicTypeEnumSet.Reg("timestamp")
	BasicTypeFormatDate = BasicTypeEnumSet.Reg("date")
	BasicTypeRandom     = BasicTypeEnumSet.Reg("random")
)

// value type
var (
	ValueTypeEnumSet   = NewEnumSet(nil)
	ValueTypeKeyID     = ValueTypeEnumSet.Reg("keyid")
	ValueTypeKeySecret = ValueTypeEnumSet.Reg("keysecret")
	ValueTypeSignature = ValueTypeEnumSet.Reg("signature")
	ValueTypeExpire    = ValueTypeEnumSet.Reg("expire")
	ValueTypeMethod    = ValueTypeEnumSet.Reg("method")
	ValueTypeAuthURL   = ValueTypeEnumSet.Reg("authurl")
	ValueTypeSignCmd   = ValueTypeEnumSet.Reg("signcmd")
)

// value in
var (
	ValueInEnumSet = NewEnumSet(nil)
	ValueInBody    = ValueInEnumSet.Reg(HTTPBody.Val())
	ValueInQuery   = ValueInEnumSet.Reg(HTTPQuery.Val())
	ValueInHeader  = ValueInEnumSet.Reg(HTTPHeader.Val())
)

// value from
var (
	ValueFromEnumSet   = NewEnumSet(nil)
	ValueFromBodyReq   = ValueFromEnumSet.Reg(HTTPBody.Val())
	ValueFromQueryReq  = ValueFromEnumSet.Reg(HTTPQuery.Val())
	ValueFromHeaderReq = ValueFromEnumSet.Reg(HTTPHeader.Val())
)
