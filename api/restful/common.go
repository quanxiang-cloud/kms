package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	error2 "github.com/quanxiang-cloud/cabin/error"
)

// ============================ common ============================
// there are common funs used by restful controller

const (
	// HeaderUserID parm User-Id in header
	HeaderUserID = "User-Id"
	// HeaderUserName param User-Name in header
	HeaderUserName = "User-Name"
)

// getUserID getUserID
func getUserID(c *gin.Context) string {
	return c.Request.Header.Get(HeaderUserID)
}

func getUserName(c *gin.Context) string {
	return c.Request.Header.Get(HeaderUserName)
}

// bindBody repeatedly bind body
func bindBody(c *gin.Context, d interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	bb, ok := b.(binding.BindingBody)
	if !ok {
		return error2.NewErrorWithString(error2.ErrParams, "binding type error")
	}
	if err := c.ShouldBindBodyWith(d, bb); err != nil {
		return error2.NewErrorWithString(error2.ErrParams, err.Error())
	}
	return nil
}
