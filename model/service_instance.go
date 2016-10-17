package model

// ServiceInstance : serializable format for JSON service instance request object
type ServiceInstance struct {
	OrganizationGUID  string     `json:"organization_guid"`
	PlanID            string     `json:"plan_id"`
	ServiceID         string     `json:"service_id"`
	SpaceGUID         string     `json:"space_guid"`
	Parameters        Parameters `json:"parameters, omitempty"`
	AcceptsIncomplete bool       `json:"accepts_incomplete, omitempty"`
}

// Parameters : gives the size in Gigabytes as specified by user
type Parameters struct {
	SizeInGB        string `json:"sizeinGB, omitempty"`
	StoragePoolName string `json:"storage_pool_name"`
}

// CreateServiceInstanceResponse : the response object containing DashboardURL
type CreateServiceInstanceResponse struct {
	DashboardURL string `json:"dashboard_url"`
}
