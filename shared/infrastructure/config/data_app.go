package config

type AppControllerService struct {
	Port int `json:"port"`
}

type AppAgentService struct {
	Port int `json:"port"`
}

type AppWorkerService struct {
	Port int `json:"port"`
}

type ApplicationServer struct {
	AppControllerService AppControllerService `json:"app_controller_service"`
	AppAgentService      AppAgentService      `json:"app_agent_service"`
	AppWorkerService     AppWorkerService     `json:"app_worker_service"`
}
