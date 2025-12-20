package application

import (
	"distributed-configuration-management/domain_workerappservice/controller/restapi"
	"distributed-configuration-management/domain_workerappservice/usecase/getconfig"
	"distributed-configuration-management/domain_workerappservice/usecase/runconfig"
	"distributed-configuration-management/shared/driver"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/server"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/pkg/workersdk"
	"fmt"
	"time"
)

type appworker struct {
	httpHandler *server.GinHTTPHandler
	controller  driver.Controller
}

func (c appworker) RunApplication() {
	c.controller.RegisterRouter()
	c.httpHandler.RunApplication()
}

func NewAppWorker() func() driver.RegistryContract {
	return func() driver.RegistryContract {

		cfg := config.ReadConfig()

		appID := util.GenerateID(4)

		appData := driver.NewApplicationData("appworker", appID)

		workerSource := workersdk.New(workersdk.Options{
			Timeout: time.Duration(cfg.Worker.WorkerTimeout) * time.Second,
		})

		log := logger.NewSimpleJSONLogger(appData)

		appAddress := fmt.Sprintf(":%d", cfg.ApplicationServer.AppWorkerService.Port)
		httpHandler := server.NewGinHTTPHandler(log, appAddress, appData)

		return &appworker{
			httpHandler: &httpHandler,
			controller: &restapi.Controller{
				Log:             log,
				Config:          cfg,
				Router:          httpHandler.Router,
				RunConfigInport: runconfig.NewUsecase(workerSource),
				GetConfigInport: getconfig.NewUsecase(workerSource),
			},
		}
	}
}
