package enums

// Http enums describe the position of params in http request & response
var (
	HTTPEnumSet = NewEnumSet(nil)
	HTTPBody    = HTTPEnumSet.Reg("body")
	HTTPHeader  = HTTPEnumSet.Reg("header")
	HTTPQuery   = HTTPEnumSet.Reg("query")
	HTTPPath    = HTTPEnumSet.Reg("path")
	HTTPCookie  = HTTPEnumSet.Reg("cookie")
)

// http method set
var (
	HTTPMethodSet = NewEnumSet(nil)
	HTTPPost      = HTTPMethodSet.Reg("POST")
	HTTPGet       = HTTPMethodSet.Reg("GET")
	HTTPPut       = HTTPMethodSet.Reg("PUT")
)
