package model

const (
	ScaleIOPlanGuid = "92798c7d-e7b0-49d6-8872-4aeafbb193ef"
	IsilonPlanGuid  = "6bf9bd41-d436-4ddc-9f7b-7a597a9ccbbb"
)

type ServicePlan struct {
	Name        string      `json:"name"`
	Id          string      `json:"id"`
	Description string      `json:"description"`
	Metadata    interface{} `json:"metadata, omitempty"`
	Free        bool        `json:"free, omitempty"`
}
