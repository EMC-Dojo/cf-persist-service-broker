package libstoragewrapper_test

import (
	"os"
	"strings"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/client"

	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	instanceID := "3d7e25a9-849a-4e19-bdb1-baddaf878f1c"

	Describe("Libstorage Client Integration", func() {
		var libsClient types.Client
		var volumeID string
		var storagePoolName string
		var serviceInstance model.ServiceInstance

		ctx := context.Background()

		BeforeEach(func() {
			config, err := model.GetConfig(strings.NewReader(""))
			Expect(err).NotTo(HaveOccurred())

			libsClient, err = client.New(config)
			Expect(err).ToNot(HaveOccurred())

			storagePoolName = os.Getenv("SCALEIO_STORAGE_POOL_NAME")
			Expect(storagePoolName).ToNot(BeEmpty())
			serviceInstance = model.ServiceInstance{
				Parameters: map[string]string{
					"storage_pool_name": storagePoolName,
				},
			}

			volumeName, err := utils.GenerateVolumeName(instanceID, serviceInstance)
			Expect(err).ToNot(HaveOccurred())

			volumeCreateOpts, err := utils.CreateVolumeOpts(serviceInstance)
			Expect(err).ToNot(HaveOccurred())

			volume, err := libstoragewrapper.CreateVolume(libsClient, ctx, volumeName, volumeCreateOpts)
			Expect(err).ToNot(HaveOccurred())

			volumeID = volume.ID
		})

		AfterEach(func() {
			err := libstoragewrapper.RemoveVolume(libsClient, ctx, volumeID)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when inspecting a volume", func() {
			It("returns the volume", func() {
				opts := &types.VolumeInspectOpts{
					Attachments: true,
				}
				volume, err := libsClient.Storage().VolumeInspect(ctx, volumeID, opts)
				Expect(err).ToNot(HaveOccurred())
				Expect(volume.Name).To(Equal("3d7e25a9849a4e19bdb1baddaf878f1"))
				Expect(volume.Size).To(Equal(int64(8)))
			})
		})

		Context("When passing in an instanceID", func() {
			It("return a volume ID if instanceID exist", func() {
        service_id, plan_id := "x", "y"
				getVolumeID, err := libstoragewrapper.GetVolumeID(libsClient, instanceID, service_id, plan_id)
				Expect(err).ToNot(HaveOccurred())
				Expect(getVolumeID).To(Equal(volumeID))
			})
		})
	})
})
