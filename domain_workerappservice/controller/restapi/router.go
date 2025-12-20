package restapi

import (
	"distributed-configuration-management/domain_workerappservice/usecase/getconfig"
	"distributed-configuration-management/domain_workerappservice/usecase/runconfig"

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

	GetConfigInport getconfig.Inport
	RunConfigInport runconfig.Inport
}

// RegisterRouter registering all the router
func (r *Controller) RegisterRouter() {
	// privateAdminOnly := r.Router.Group("/api/v1", r.authenticatedAdmin())
	public := r.Router.Group("/api/v1")
	// private := r.Router.Group("/api/v1", r.authenticated())

	// config
	public.GET("/hit", r.authorized(), r.getConfigHandler(r.GetConfigInport))
	public.POST("/config", r.authorized(), r.runConfigHandler(r.RunConfigInport))
}
