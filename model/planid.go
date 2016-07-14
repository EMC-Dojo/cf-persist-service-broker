package model

// PlanID : defines format for PlanID as plan HostName and PlanID. Used for serializing in JSON format
type PlanID struct {
	LibsHostName    string `json:"h"`
	LibsServiceName string `json:"p"`
	LibsDriverName  string `json:"d"`
}
