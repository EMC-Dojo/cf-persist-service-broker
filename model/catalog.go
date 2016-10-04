package model

// Catalog : type that defines json struct for returning catalog to CC
type Catalog struct {
	Services []Service `json:"services"`
}

// Service : struct to nest inside of Catalog, holds service information
type Service struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Bindable    bool            `json:"bindable"`
	Requires    []string        `json:"requires"`
	Plans       []Plan          `json:"plans"`
	Metadata    ServiceMetadata `json:"metadata"`
}

// Plan : struct to nest inside of Service, hold ID/Name/Description of Plan
type Plan struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Metadata    PlanMetadata `json:"metadata"`
}

type ServiceMetadata struct {
	DisplayName         string `json:"displayName"`
	ImageUrl            string `json:"imageUrl"`
	LongDescription     string `json:"longDescription"`
	ProviderDisplayName string `json:"providerDisplayName"`
	DocumentationUrl    string `json:"documentationUrl"`
	SupportUrl          string `json:"supportUrl"`
}

type PlanMetadata struct {
	Bullets     []string `json:"bullets"`
	DisplayName string   `json:"displayName"`
}
