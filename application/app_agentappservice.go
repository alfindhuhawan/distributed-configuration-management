package application

import (
	"context"
	"distributed-configuration-management/shared/driver"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/server"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/pkg/agent"
	"fmt"
	"time"
)

type appagent struct {
	httpHandler *server.GinHTTPHandler
}

func (c appagent) RunApplication() {
	c.httpHandler.RunApplication()
}

func NewAppAgent() func() driver.RegistryContract {
	return func() driver.RegistryContract {

		cfg := config.ReadConfig()

		appID := util.GenerateID(4)

		appData := driver.NewApplicationData("appagent", appID)

		log := logger.NewSimpleJSONLogger(appData)

		appAddress := fmt.Sprintf(":%d", cfg.ApplicationServer.AppAgentService.Port)
		httpHandler := server.NewGinHTTPHandler(log, appAddress, appData)

		traceIDAgent := util.GenerateID(16)
		ctxAgent := logger.SetTraceID(context.Background(), traceIDAgent)

		// init agentcacheloc
		agentCacheLoc := cfg.AgentInit.AgentCacheLoc
		/*
			since worker use in-memory , its need for handling scenario :
				- controller service not start
				- agent service start
				- worker service start

				1. agent_cache.json exist and url is not empty
				2. code below handle worker using previous config from agent_cache.json even the controller service not working on the start
		*/
		isExist := util.FileExists(agentCacheLoc)
		if isExist {
			// get data
			agentCacheObj := config.ReadAgentCache(agentCacheLoc)
			if agentCacheObj.URL != "" {
				// hit config worker
				agent.ConfigWorker(ctxAgent, log, cfg, agentCacheObj.URL)
			}
		}

		// auto register
		respAutoRegister, err := agent.AutoRegister(cfg)
		if err != nil {
			log.Error(ctxAgent, "error auto register : "+err.Error())
		} else {
			// return response
			log.Info(ctxAgent, util.DebuggingStructWithoutError(respAutoRegister), nil)
		}

		// dynamic polling
		// set poll interval seconds
		pollIntervalSeconds := cfg.PollIntervalSeconds
		if respAutoRegister != nil {
			if respAutoRegister.PollIntervalSeconds > 0 {
				pollIntervalSeconds = respAutoRegister.PollIntervalSeconds
			}
		}

		go agent.StartDynamicPolling(ctxAgent, log, cfg, time.Duration(pollIntervalSeconds)*time.Second, agentCacheLoc)

		return &appagent{
			httpHandler: &httpHandler,
		}
	}
}
