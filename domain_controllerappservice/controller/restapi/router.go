package restapi

import (
	"distributed-configuration-management/domain_controllerappservice/usecase/getconfig"
	"distributed-configuration-management/domain_controllerappservice/usecase/registeragent"
	"distributed-configuration-management/domain_controllerappservice/usecase/updateconfig"

	"github.com/gin-gonic/gin"

	"distributed-configuration-management/shared/helper"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/logger"
)

type Controller struct {
	Router gin.IRouter
	Config *config.Config
	Log    logger.Logger
	Helper helper.HTTPHelper

	GetConfigInport     getconfig.Inport
	RegisterAgentInport registeragent.Inport
	UpdateConfigInport  updateconfig.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	// privateAdminOnly := r.Router.Group("/api/v1", r.authenticatedAdmin())
	public := r.Router.Group("/api/v1")
	// private := r.Router.Group("/api/v1", r.authenticated())

	// config
	public.GET("/config", r.authenticatedAgent(), r.getConfigHandler(r.GetConfigInport))
	public.POST("/config", r.authenticatedAdmin(), r.updateConfigHandler(r.UpdateConfigInport))
	public.POST("/register", r.authenticatedAgent(), r.registerAgentHandler(r.RegisterAgentInport))
}
