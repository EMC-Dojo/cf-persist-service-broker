package model

// ServiceInstance : serializable format for JSON service instance request object
type ServiceInstance struct {
	OrganizationGUID  string            `json:"organization_guid"`
	PlanID            string            `json:"plan_id"`
	ServiceID         string            `json:"service_id"`
	SpaceGUID         string            `json:"space_guid"`
	Parameters        map[string]string `json:"parameters, omitempty"`
	AcceptsIncomplete bool              `json:"accepts_incomplete, omitempty"`
}

// CreateServiceInstanceResponse : the response object containing DashboardURL
type CreateServiceInstanceResponse struct {
	DashboardURL string `json:"dashboard_url"`
}
