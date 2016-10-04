package config_test

import (
	"github.com/EMC-Dojo/cf-persist-service-broker/config"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var plans []model.Plan
	var services []model.Service
	BeforeEach(func() {
		plans = []model.Plan{
			model.Plan{
				Name:        "isilonservice",
				Description: "An isilon service",
				Metadata: model.PlanMetadata{
					DisplayName: "isilon",
					Bullets: []string{
						"Brings you isilon service",
					},
				},
			},
			model.Plan{
				Name:        "scaleioservice",
				Description: "A scaleio service",
				Metadata: model.PlanMetadata{
					DisplayName: "scaleio",
					Bullets: []string{
						"Brings you scaleio service",
					},
				},
			},
		}

		services = []model.Service{
			model.Service{
				ID:          "fake-uuid",
				Name:        "Persistent-Storage",
				Description: "Supports EMC ScaleIO & Isilon Storage Arrays for use with CloudFoundry",
				Bindable:    true,
				Requires:    []string{"volume_mount"},
				Plans:       plans,
				Metadata: model.ServiceMetadata{
					DisplayName:         "Persistent Storage",
					ImageUrl:            "imageURL",
					ProviderDisplayName: "Dell EMC",
					LongDescription:     "Dell EMC brings you persistent storage on CloudFoundry",
					DocumentationUrl:    "docsURL",
					SupportUrl:          "supportURL",
				},
			},
		}
	})
	Context("When service uuid is not valid", func() {
		It("return error indicate invalid uuid", func() {
			_, err := config.ReadConfig("../fixtures/not-valid-uuid-config.json")
			Expect(err).To(MatchError("the service uuid given is not a valid uuid: not-valid-uuid"))
		})
	})

	Context("Parse Settings file", func() {
		It("return the correct structure", func() {
			services[0].ID = "92e98925-d046-4c72-9598-ba352449a5c7"
			configData, err := config.ReadConfig("../fixtures/valid-uuid-config.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(configData).To(Equal(services))
		})
	})

	Context("When config path is empty", func() {
		It("use default config defined in config.go", func() {
			configData, err := config.ReadConfig("")
			Expect(err).ToNot(HaveOccurred())
			Expect(configData).To(Equal(config.DefaultConfig))
		})
	})

})
