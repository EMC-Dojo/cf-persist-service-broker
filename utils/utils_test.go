package utils_test

import (
	"github.com/EMC-CMD/cf-persist-service-broker/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "github.com/EMC-CMD/cf-persist-service-broker/model"
)

var _ = Describe("Utils/Utils", func() {
  var serviceInstance model.ServiceInstance

  BeforeEach(func() {
    serviceInstance = model.ServiceInstance{
      Parameters: map[string]interface{}{
        "storage_pool_name": "default",
      },
    }
  })

  Describe("Generate Volume Name", func() {
    Context("When scaleio", func() {
      It("Should Create Volume Name using instanceId and serviceInstance", func () {
        instanceId := "3c653bce-8752-451b-96d9-a8a1a925b118"
        expected := "3c653bce8752451b96d9a8a1a925b11"
        cooked, err := utils.GenerateVolumeName(instanceId, serviceInstance)
        Expect(err).ToNot(HaveOccurred())
        Expect(cooked).To(Equal(expected))
      })
      It("Should return an error if it is longer than 32 characters after removing all hyphens", func () {
        instanceId := "3c653bce-8752-451b-96d9-a8a1a925b1189"
        _, err := utils.GenerateVolumeName(instanceId, serviceInstance)
        Expect(err).To(MatchError("Volume name cannot exceed 32 characters when all hyphens are removed."))
      })
      It("Should leave the volume name alone if it is shorter than 31 characters", func () {
        cooked, err := utils.GenerateVolumeName("123-456-789", serviceInstance)
        Expect(err).ToNot(HaveOccurred())
        Expect(cooked).To(Equal("123-456-789"))
      })
    })
  })

  Describe("CreateVolumeOpts", func() {
    Context("When scaleio", func() {
      It("Should Create Volume Options using a serviceInstance", func() {
        volumeCreateOptions, err := utils.CreateVolumeOpts(serviceInstance)
        Expect(err).ToNot(HaveOccurred())
        Expect(*volumeCreateOptions.AvailabilityZone).To(Equal("az"))
        Expect(*volumeCreateOptions.IOPS).To(Equal(int64(100)))
        Expect(*volumeCreateOptions.Size).To(Equal(int64(8)))
        Expect(*volumeCreateOptions.Type).To(Equal("default"))
      })
    })
  })
})
