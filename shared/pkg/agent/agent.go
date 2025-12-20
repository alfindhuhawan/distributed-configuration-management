package agent

import (
	"bytes"
	"context"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/errorenum"
	"distributed-configuration-management/shared/model/repository"
	"distributed-configuration-management/shared/model/request"
	"distributed-configuration-management/shared/model/response"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var retryRunning int32 // 0 = now working, 1 = on progress , use for loop poll

func AutoRegister(cfg *config.Config) (*repository.RegisterAgentResponse, error) {
	var err error
	var registerResponse response.RegisterResponse

	// check if config agents empty because its needed for register
	agent := cfg.AgentInit
	if agent.Name == "" {
		return nil, errorenum.ErrorEmptyAgent
	}

	// set agent request for hit register
	agentObj := request.AgentRequest{
		Name: agent.Name,
		IP:   agent.IP,
	}

	byteJson, _ := json.Marshal(agentObj)

	url := cfg.ApplicationInternal.ControllerService.RegisterURL

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteJson))
	if err != nil {
		return nil, errorenum.ErrorInternalOnAgent.Var("autoRegister New request : " + err.Error())
	}

	// get authorization
	agentAuth := util.CreateSHA256Signature(cfg.Credentials.AgentCredential.SecretKey)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+agentAuth)

	var client = &http.Client{
		Timeout: time.Second * 60,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errorenum.ErrorInternalOnAgent.Var("autoRegister client do request : " + err.Error())
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errorenum.ErrorInternalOnAgent.Var("autoRegister io readall : " + err.Error())
	}

	err = json.Unmarshal(responseBody, &registerResponse)
	if err != nil {
		return nil, errorenum.ErrorInternalOnAgent.Var("autoRegister json unmarshal : " + err.Error())
	}

	// if response not 200
	if registerResponse.Code != 200 {
		return nil, errorenum.ErrorRegisterAgent.Var(registerResponse.CodeMessage)
	}

	// check data agent for cache start
	urlCache := agent.AgentCacheLoc
	isExist := util.FileExists(urlCache)
	if !isExist {
		// create new file if not exist
		file, err := os.Create(urlCache)
		if err != nil {
			return nil, errorenum.ErrorCreateAgentCache.Var(registerResponse.CodeMessage)
		}
		defer file.Close()
	}

	// save file
	err = util.SaveCache(urlCache, config.AgentCache{
		AgentID:             registerResponse.Data.AgentID,
		AgentName:           agentObj.Name,
		LastVersion:         0,
		PollIntervalSeconds: registerResponse.Data.PollIntervalSeconds,
		URL:                 registerResponse.Data.PollURL,
		LastPollAt:          "",
	})
	if err != nil {
		return nil, errorenum.ErrorSaveAgentCache.Var(registerResponse.CodeMessage)
	}

	return &registerResponse.Data, err
}

func StartDynamicPolling(ctxAgent context.Context, log logger.Logger, cfg *config.Config, initialInterval time.Duration, pathAgentCache string) {
	interval := initialInterval

	for {
		time.Sleep(interval)

		newInterval := pollData(ctxAgent, log, cfg, interval, pathAgentCache)

		// use new interval if poll data return new interval
		if newInterval > 0 && newInterval != interval {
			interval = newInterval
		}
	}
}

func pollData(ctxAgent context.Context, log logger.Logger, cfg *config.Config, currentInterval time.Duration, pathAgentCache string) time.Duration {
	log.Info(ctxAgent, "pollData running, interval: "+currentInterval.String(), nil)

	// get last version based on cache
	var lastVersion float64
	agentCache, _ := util.LoadCache(pathAgentCache)
	if agentCache != nil {
		lastVersion = agentCache.LastVersion
	} else {
		// if agent cache it means there is a probability that agent_cache not create or still empty and controller service not working, need to hit auto register here on polldata
		_, _ = AutoRegister(cfg)

		// reload agent cache
		agentCache, _ = util.LoadCache(pathAgentCache)
	}

	configContResp, errCode, err := getConfigController(ctxAgent, log, cfg, lastVersion)
	if err != nil {
		if errCode != 304 {
			log.Error(ctxAgent, "initial fetch failed, starting infinite backoff retry loop…", nil)

			// not create any other loop if prev loop still exist
			if !atomic.CompareAndSwapInt32(&retryRunning, 0, 1) {
				log.Info(ctxAgent, "retry loop already running, skip creating new one", nil)
				return 0
			}

			// run retry unlimited
			go func() {
				defer atomic.StoreInt32(&retryRunning, 0) // reset flag after success

				newCfg, err := fetchConfigInfiniteBackoff(ctxAgent, log, cfg, pathAgentCache)
				if err == nil {
					log.Info(ctxAgent, "config fetched successfully after backoff", nil)
					// fmt.Println("config fetched successfully after backoff:", newCfg)
					configContResp = newCfg
				}
			}()

			// still using last version
			return 0
		}
	}

	// if success
	if configContResp != nil {
		if configContResp.Version > lastVersion {
			timeNow := time.Now()
			strTimeNow := timeNow.Format(time.RFC3339)

			// save file
			err = util.SaveCache(pathAgentCache, config.AgentCache{
				AgentID:             agentCache.AgentID,
				AgentName:           agentCache.AgentName,
				LastVersion:         configContResp.Version,
				PollIntervalSeconds: configContResp.PollIntervalSeconds,
				URL:                 configContResp.URL,
				LastPollAt:          strTimeNow,
			})
			if err != nil {
				log.Error(ctxAgent, string(errorenum.ErrorSaveAgentCache.Var(err.Error())), nil)
				return 0
			}

			// hit worker POST /config here
			ConfigWorker(ctxAgent, log, cfg, configContResp.URL)

			// set new interval
			if configContResp.PollIntervalSeconds > 0 {
				newInterval := time.Duration(configContResp.PollIntervalSeconds) * time.Second
				if newInterval != currentInterval {
					return newInterval
				}
			}

		}
	}

	return 0
}

// fetchConfigInfiniteBackoff tries forever until success.
// NEVER stops unless program stops.
func fetchConfigInfiniteBackoff(ctxAgent context.Context, log logger.Logger, cfg *config.Config, pathAgentCache string) (*repository.GlobalConfigResponse, error) {
	delay := cfg.BaseDelaySeconds
	maxDelay := cfg.MaxDelaySeconds

	// get last version based on cache
	var lastVersion float64
	agentCache, _ := util.LoadCache(pathAgentCache)
	if agentCache != nil {
		lastVersion = agentCache.LastVersion
	}

	for {
		configContResp, errCode, err := getConfigController(ctxAgent, log, cfg, lastVersion)
		if err == nil || errCode == 304 {
			return configContResp, nil
		}

		log.Info(ctxAgent, "[config fetch] failed: "+err.Error())
		log.Info(ctxAgent, "[config fetch] retrying in:"+strconv.Itoa(delay))
		// fmt.Println("[config fetch] failed:", err)

		time.Sleep(time.Duration(delay) * time.Second)

		// increase delay exponentially
		delay = delay * 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
}

func getConfigController(ctxAgent context.Context, log logger.Logger, cfg *config.Config, lastVersion float64) (*repository.GlobalConfigResponse, int, error) {
	var err error
	var getConfigControllerResp response.GetConfigControllerResponse

	url := cfg.ApplicationInternal.ControllerService.GetConfigURL

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 500, errorenum.ErrorInternalOnAgent.Var("getConfigController new request : " + err.Error())
	}

	// get authorization
	agentAuth := util.CreateSHA256Signature(cfg.Credentials.AgentCredential.SecretKey)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+agentAuth)
	request.Header.Set("ETag", strconv.FormatFloat(lastVersion, 'f', -1, 64))

	var client = &http.Client{
		Timeout: time.Second * 60,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, 500, errorenum.ErrorInternalOnAgent.Var("getConfigController client do : " + err.Error())
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 500, errorenum.ErrorInternalOnAgent.Var("getConfigController io readall : " + err.Error())
	}

	err = json.Unmarshal(responseBody, &getConfigControllerResp)
	if err != nil {
		return nil, 500, errorenum.ErrorInternalOnAgent.Var("getConfigController json unmarshal : " + err.Error())
	}

	if getConfigControllerResp.Code != 200 {
		if getConfigControllerResp.Code == 304 {
			return nil, 304, errorenum.ErrorGetConfigControllerNoChange.Var(getConfigControllerResp.CodeMessage)
		}

		return nil, getConfigControllerResp.Code, errorenum.ErrorGetConfigController.Var(getConfigControllerResp.CodeMessage)
	}

	return &getConfigControllerResp.Data, 200, nil
}

func ConfigWorker(ctxAgent context.Context, log logger.Logger, cfg *config.Config, newURL string) {
	url := cfg.AgentInit.WorkerURL

	reqConfigWorker := request.ConfigWorker{
		URL: newURL,
	}
	byteJson, _ := json.Marshal(reqConfigWorker)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteJson))
	if err != nil {
		log.Error(ctxAgent, string(errorenum.ErrorInternalOnWorker.Var("configWorker New request : "+err.Error())))
	}

	request.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 60,
	}

	responseReq, err := client.Do(request)
	if err != nil {
		log.Error(ctxAgent, string(errorenum.ErrorInternalOnWorker.Var("configWorker client do request : "+err.Error())))
		return // return from this function if error
	}

	defer func() {
		if err := responseReq.Body.Close(); err != nil {
			log.Error(ctxAgent, "failed to close worker response body: "+err.Error())
		}
	}()

	responseBody, err := io.ReadAll(responseReq.Body)
	if err != nil {
		log.Error(ctxAgent, string(errorenum.ErrorInternalOnWorker.Var("configWorker io readall : "+err.Error())))
	}

	var configWorkerResponse response.ConfigWorkerResponse
	err = json.Unmarshal(responseBody, &configWorkerResponse)
	if err != nil {
		log.Error(ctxAgent, string(errorenum.ErrorInternalOnWorker.Var("configWorker json unmarshal : "+err.Error())))
	}

	// if response not 200
	if configWorkerResponse.Code != 200 {
		log.Error(ctxAgent, string(errorenum.ErrorConfigWorker.Var(configWorkerResponse.CodeMessage)))
	}
}
