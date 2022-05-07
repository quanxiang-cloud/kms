package restful

import (
	"kms/internal/polysign"
	"kms/internal/service"
	"kms/pkg/misc/config"

	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// KeyAgent ext key
type KeyAgent struct {
	e service.KeyAgent
}

// NewKeyAgent new
func NewKeyAgent(conf *config.Config, log logger.AdaptedLogger) (*KeyAgent, error) {
	extKeyManager, err := service.CreateKeyAgent(conf, log)
	if err != nil {
		return nil, err
	}

	return &KeyAgent{
		e: extKeyManager,
	}, nil
}

// Upload upload
func (agent *KeyAgent) Upload(c *gin.Context) {
	req := &service.UploadAgencyReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)
	req.UserName = getUserName(c)
	resp.Format(agent.e.Upload(header.MutateContext(c), req)).Context(c)
}

// Query query
func (agent *KeyAgent) Query(c *gin.Context) {
	req := &service.QueryAgencyReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	resp.Format(agent.e.Query(header.MutateContext(c), req)).Context(c)
}

// List list
func (agent *KeyAgent) List(c *gin.Context) {
	req := &service.ListAgencyReq{
		Active: -1,
	}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	req.Owner = getUserID(c)

	resp.Format(agent.e.List(header.MutateContext(c), req)).Context(c)
}

// Delete delete
func (agent *KeyAgent) Delete(c *gin.Context) {
	req := &service.DeleteAgencyReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	req.UserID = getUserID(c)

	resp.Format(agent.e.Delete(header.MutateContext(c), req)).Context(c)
}

// DeleteInBatch delete in batch
func (agent *KeyAgent) DeleteInBatch(c *gin.Context) {
	req := &service.DeleteInBatchReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	resp.Format(agent.e.DeleteInBatch(header.MutateContext(c), req)).Context(c)
}

// Authorize auth 3party, get access token
func (agent *KeyAgent) Authorize(c *gin.Context) {
	req := &service.AuthReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	req.APIServiceArgs = c.GetHeader(polysign.XHeaderPolyServiceArgs)
	resp.Format(agent.e.Authorize(header.MutateContext(c), req)).Context(c)
}

// ListAuthType list auth type
func (agent *KeyAgent) ListAuthType(c *gin.Context) {
	req := &service.ListAuthTypeReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}

	resp.Format(agent.e.ListAuthType(header.MutateContext(c), req)).Context(c)
}

// GetSample get auth sample
func (agent *KeyAgent) GetSample(c *gin.Context) {
	req := &service.GetSampleReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	resp.Format(agent.e.GetSample(header.MutateContext(c), req)).Context(c)
}

// Update update authType & authContent by service
func (agent *KeyAgent) Update(c *gin.Context) {
	req := &service.UpdateAgencyReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	resp.Format(agent.e.Update(header.MutateContext(c), req)).Context(c)
}

// UpdateInBatch UpdateInBatch
func (agent *KeyAgent) UpdateInBatch(c *gin.Context) {
	req := &service.UpdateInBatchReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	resp.Format(agent.e.UpdateInBatch(header.MutateContext(c), req)).Context(c)
}

// Active update agency key status
func (agent *KeyAgent) Active(c *gin.Context) {
	req := &service.ActiveAgencyReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	resp.Format(agent.e.Active(header.MutateContext(c), req)).Context(c)
}

// DeleteByPrefixPath DeleteByPrefixPath
func (agent *KeyAgent) DeleteByPrefixPath(c *gin.Context) {
	req := &service.DeleteByPrefixPathReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	resp.Format(agent.e.DeleteByPrefixPath(header.MutateContext(c), req)).Context(c)
}

// CheckAuth CheckAuth
func (agent *KeyAgent) CheckAuth(c *gin.Context) {
	req := &service.CheckAuthReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error())).Context(c)
		return
	}
	resp.Format(agent.e.CheckAuth(header.MutateContext(c), req)).Context(c)
}
