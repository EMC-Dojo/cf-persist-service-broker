package model

type ServiceInstance struct {
	OrganizationGuid string            `json:"organization_guid"`
	PlanId           string            `json:"plan_id"`
	ServiceId        string            `json:"service_id"`
	SpaceGuid        string            `json:"space_guid"`
	Parameters       map[string]string `json:"parameters, omitempty"`
	AcceptIncomplete bool              `json:"accepts_incomplete, omitempty"`
}

type CreateServiceInstanceResponse struct {
	DashboardUrl string `json:"dashboard_url"`
}
