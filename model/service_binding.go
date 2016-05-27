package model

type ServiceBinding struct {
	ServiceId    string                 `json:"service_id"`
	AppId        string                 `json:"app_guid"`
	PlanId       string                 `json:"plan_id"`
	BindResource map[string]string      `json:"bind_resource"`
	Parameters   map[string]interface{} `json:"parameters"`
}

type CreateServiceBindingResponse struct {
	Credentials  CreateServiceBindingCredentials `json:"credentials"`
	VolumeMounts []VolumeMount                   `json:"volume_mounts"`
}

type CreateServiceBindingCredentials struct {
	URI string `json:"uri"`
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
