package restful

import (
	"kms/internal/service"
	"kms/pkg/misc/config"

	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// KeyConfig key config
type KeyConfig struct {
	s service.KeyConfigServe
}

// NewKeyConfig new key config
func NewKeyConfig(conf *config.Config, log logger.AdaptedLogger) (*KeyConfig, error) {
	kcServe, err := service.NewKeyConfigServe(conf, log)
	if err != nil {
		return nil, err
	}

	return &KeyConfig{
		s: kcServe,
	}, nil
}

// Update update
func (kc *KeyConfig) Update(c *gin.Context) {
	req := &service.UpdateKConfigReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}

	resp.Format(kc.s.Update(header.MutateContext(c), req)).Context(c)
}

// Query query
func (kc *KeyConfig) Query(c *gin.Context) {
	req := &service.QueryKConfigReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}

	resp.Format(kc.s.Query(header.MutateContext(c), req)).Context(c)
}

// Delete delete
func (kc *KeyConfig) Delete(c *gin.Context) {
	req := &service.DeleteKConfigReq{}
	if err := c.ShouldBind(req); err != nil {
		resp.Format(nil, error2.New(error2.ErrParams)).Context(c)
		return
	}

	resp.Format(kc.s.Delete(header.MutateContext(c), req)).Context(c)
}
