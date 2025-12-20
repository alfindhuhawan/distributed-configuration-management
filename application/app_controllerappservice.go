package application

import (
	"distributed-configuration-management/domain_controllerappservice/controller/restapi"
	"distributed-configuration-management/domain_controllerappservice/gateway/prod"
	"distributed-configuration-management/domain_controllerappservice/usecase/getconfig"
	"distributed-configuration-management/domain_controllerappservice/usecase/registeragent"
	"distributed-configuration-management/domain_controllerappservice/usecase/updateconfig"
	"distributed-configuration-management/shared/driver"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/server"
	"distributed-configuration-management/shared/infrastructure/util"
	"fmt"
)

type appcontroller struct {
	httpHandler *server.GinHTTPHandler
	controller  driver.Controller
}

func (c appcontroller) RunApplication() {
	c.controller.RegisterRouter()
	c.httpHandler.RunApplication()
}

func NewAppController() func() driver.RegistryContract {
	return func() driver.RegistryContract {

		cfg := config.ReadConfig()

		appID := util.GenerateID(4)

		appData := driver.NewApplicationData("appcontroller", appID)

		log := logger.NewSimpleJSONLogger(appData)

		appAddress := fmt.Sprintf(":%d", cfg.ApplicationServer.AppControllerService.Port)
		httpHandler := server.NewGinHTTPHandler(log, appAddress, appData)

		datasource := prod.NewGateway(log, appData, cfg)

		return &appcontroller{
			httpHandler: &httpHandler,
			controller: &restapi.Controller{
				Log:                 log,
				Config:              cfg,
				Router:              httpHandler.Router,
				RegisterAgentInport: registeragent.NewUsecase(datasource),
				GetConfigInport:     getconfig.NewUsecase(datasource),
				UpdateConfigInport:  updateconfig.NewUsecase(datasource),
			},
		}
	}
}
