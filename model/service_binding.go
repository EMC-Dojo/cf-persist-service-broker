package model

type ServiceBinding struct {
	ServiceID    string               `json:"service_id"`
	AppID        string               `json:"app_guid"`
	PlanID       string               `json:"plan_id"`
	BindResource map[string]string    `json:"bind_resource"`
	Parameters   ServiceBindingParams `json:"parameters"`
}

type ServiceBindingParams struct {
	Driver string `json:"volume_driver"`
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
	Driver       string `json:"driver"`
	ContainerDir string `json:"container_dir"`
	Mode         string `json:"mode"`
	DeviceType   string `json:"device_type"`
	Device       Device `json:"device"`
}

type Device struct {
	VolumeID    string            `json:"volume_id"`
	MountConfig map[string]string `json:"mount_config"`
}
