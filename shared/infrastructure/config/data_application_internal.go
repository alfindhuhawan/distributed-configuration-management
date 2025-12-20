package config

type ApplicationInternal struct {
	ControllerService ControllerService `json:"controller"`
}

type ControllerService struct {
	RegisterURL  string `json:"register_url"`
	GetConfigURL string `json:"get_config_url"`
}
