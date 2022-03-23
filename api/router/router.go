package router

import (
	restful "kms/api/restful"
	"kms/pkg/misc/config"
	"kms/pkg/probe"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	ginLogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router Router
type Router struct {
	*probe.Probe

	c      *config.Config
	engine *gin.Engine
}

// NewRouter create a router
func NewRouter(c *config.Config, log logger.AdaptedLogger) (*Router, error) {
	engine, err := newRouter(c, log)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/kms")

	v1.Any("/ping", restful.PingPong)

	keyGen, err := restful.NewKeyGen(c, log)
	if err != nil {
		return nil, err
	}

	rawGroup := v1.Group("/key")
	{
		rawGroup.POST("/create", keyGen.CreateKey)
		rawGroup.POST("/delete", keyGen.DeleteKey)
		rawGroup.POST("/list", keyGen.ListKey)
		rawGroup.POST("/query", keyGen.QueryKey)
		rawGroup.POST("/active", keyGen.ActiveKey)
		rawGroup.POST("/update", keyGen.UpdateKey)
		rawGroup.POST("/signature", keyGen.Signature)
	}

	keyAgent, err := restful.NewKeyAgent(c, log)
	if err != nil {
		return nil, err
	}

	agentGroup := v1.Group("/ext")
	{
		agentGroup.POST("/upload", keyAgent.Upload)
		agentGroup.POST("/query", keyAgent.Query)
		agentGroup.POST("/list", keyAgent.List)
		agentGroup.POST("/sample", keyAgent.GetSample)
		agentGroup.POST("/delete", keyAgent.Delete)
		agentGroup.POST("/deleteInBatch", keyAgent.DeleteInBatch)
		agentGroup.POST("/authorize", keyAgent.Authorize)
		agentGroup.POST("/authTypes", keyAgent.ListAuthType)
		agentGroup.POST("/update", keyAgent.Update)
		agentGroup.POST("/updateInBatch", keyAgent.UpdateInBatch)
		agentGroup.POST("/active", keyAgent.Active)
		agentGroup.POST("/deleteByPrefix", keyAgent.DeleteByPrefixPath)
		agentGroup.POST("/checkAuth", keyAgent.CheckAuth)
	}

	keyConfig, err := restful.NewKeyConfig(c, log)
	if err != nil {
		return nil, err
	}

	configGroup := v1.Group("/config")
	{
		configGroup.POST("/update", keyConfig.Update)
		configGroup.POST("/query", keyConfig.Query)
		configGroup.POST("/delete", keyConfig.Delete)
	}

	r := &Router{
		Probe:  probe.New(log),
		c:      c,
		engine: engine,
	}

	r.probe()
	return r, nil
}

func newRouter(c *config.Config, log logger.AdaptedLogger) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(ginLogger.LoggerFunc(), ginLogger.RecoveryFunc())

	return engine, nil
}

func (r *Router) probe() {
	r.engine.GET("liveness", func(c *gin.Context) {
		r.Probe.LivenessProbe(c.Writer, c.Request)
	})

	r.engine.Any("readiness", func(c *gin.Context) {
		r.Probe.ReadinessProbe(c.Writer, c.Request)
	})
}

// Run router
func (r *Router) Run() {
	r.Probe.SetRunning()
	r.engine.Run(r.c.Port)
}

// Close router
func (r *Router) Close() {
}

func (r *Router) router() {

}
