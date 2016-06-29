package model

type ServiceBinding struct {
	ServiceID    string                 `json:"service_id"`
	AppID        string                 `json:"app_guid"`
	PlanID       string                 `json:"plan_id"`
	BindResource map[string]string      `json:"bind_resource"`
	Parameters   map[string]interface{} `json:"parameters"`
}

type CreateServiceBindingResponse struct {
	Credentials  CreateServiceBindingCredentials `json:"credentials"`
	VolumeMounts []VolumeMount                   `json:"volume_mounts"`
}

type CreateServiceBindingCredentials struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	URI      string `json:"uri"`
	Username string `json:"username"`
}

type VolumeMount struct {
	ContainerPath string                    `json:"container_path"`
	Mode          string                    `json:"mode"`
	Private       VolumeMountPrivateDetails `json:"private"`
}

type VolumeMountPrivateDetails struct {
	Driver  string `json:"driver"`
	GroupId string `json:"group_id"`
	Config  string `json:"config"`
}
