package restful

import (
	"fmt"
	"kms/internal/polysign"
	"kms/internal/service"
	"kms/pkg/misc/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

type pingResp struct {
	Msg       string `json:"msg"`
	Timestamp string `json:"timestamp"`
}

// PingPong return a pong response for ping
func PingPong(c *gin.Context) {
	r := &pingResp{
		Msg:       "pong",
		Timestamp: time.Now().Format("2006-01-02T15:04:05MST"),
	}
	resp.Format(r, nil).Context(c)
}

// KeyGen keygen
type KeyGen struct {
	k service.KeyGenerator
}

// NewKeyGen new
func NewKeyGen(conf *config.Config, log logger.AdaptedLogger) (*KeyGen, error) {
	generator, err := service.CreateKeyGenerator(conf, log)
	if err != nil {
		return nil, err
	}

	serv := &KeyGen{
		k: generator,
	}

	return serv, nil
}

// CreateKey create Key
func (keyGen *KeyGen) CreateKey(c *gin.Context) {
	req := service.CreateReq{}
	if err := c.ShouldBind(&req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)
	req.UserName = getUserName(c)
	if req.UserID == "" || req.UserName == "" {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}

	resp.Format(keyGen.k.CreateKey(header.MutateContext(c), &req)).Context(c)
}

// DeleteKey delete key
func (keyGen *KeyGen) DeleteKey(c *gin.Context) {
	req := service.DeleteReq{}
	if err := c.ShouldBind(&req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)
	if req.UserID == "" {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}

	resp.Format(keyGen.k.DeleteKey(header.MutateContext(c), &req)).Context(c)
}

// ListKey list
func (keyGen *KeyGen) ListKey(c *gin.Context) {
	req := service.ListReq{}

	if err := c.ShouldBind(&req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)
	resp.Format(keyGen.k.ListKey(header.MutateContext(c), &req)).Context(c)
}

// QueryKey query keyid by userid
func (keyGen *KeyGen) QueryKey(c *gin.Context) {
	req := service.QueryReq{}

	if err := c.ShouldBind(&req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	resp.Format(keyGen.k.QueryKey(header.MutateContext(c), &req)).Context(c)
}

// ActiveKey update key status
func (keyGen *KeyGen) ActiveKey(c *gin.Context) {
	req := service.ActiveReq{}
	if err := c.ShouldBind(&req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)
	if req.UserID == "" {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}

	resp.Format(keyGen.k.ActiveKey(header.MutateContext(c), &req)).Context(c)
}

// UpdateKey update
func (keyGen *KeyGen) UpdateKey(c *gin.Context) {
	req := service.UpdateReq{}
	if err := c.ShouldBind(&req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)
	if req.UserID == "" {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}
	resp.Format(keyGen.k.UpdateKey(header.MutateContext(c), &req)).Context(c)
}

// Signature sign
func (keyGen *KeyGen) Signature(c *gin.Context) {
	req := service.SignatureReq{
		AccessKeyID: c.GetHeader(polysign.XHeaderPolySignKeyID),
	}

	if req.AccessKeyID == "" {
		err := error2.NewErrorWithString(error2.ErrParams, fmt.Sprintf("missing header.%s", polysign.XHeaderPolySignKeyID))
		resp.Format(nil, err).Context(c)
		return
	}

	if err := bindBody(c, &req.FullBody); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	signResp, err := keyGen.k.Signature(header.MutateContext(c), &req)

	if err != nil {
		if err2, ok := err.(error2.Error); ok {
			resp.Format(nil, err2).Context(c)
		}
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.Header(polysign.XInternalHeaderPolySignSignature, signResp.Sign)
	c.Header("User-Id", signResp.UserInfo.UserID)
	c.Header("User-Name", signResp.UserInfo.UserName)
	for _, v := range signResp.UserInfo.DepartmentIDs {
		c.Writer.Header().Add("Department-Id", v)
	}
}
