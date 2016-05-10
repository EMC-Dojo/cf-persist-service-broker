package server

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strings"

	"github.com/EMC-CMD/cf-persist-service-broker/model"
	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/client"
)

var _ = Describe("Integration", func() {
	volumeName := "volumeName"

	Describe("Libstorage Client Integration", func() {
		var libsClient types.Client
		var volumeID string
		ctx := context.Background()

		BeforeEach(func() {
			config, err := model.GetConfig(strings.NewReader(""))
			Expect(err).NotTo(HaveOccurred())

			libsClient, err = client.New(config)
			Expect(err).ToNot(HaveOccurred())

			volume, err := CreateVolume(libsClient, ctx, volumeName, "az", "pool1", 100, 8)
			Expect(err).ToNot(HaveOccurred())

			volumeID = volume.ID
		})

		AfterEach(func() {
			err := RemoveVolume(libsClient, ctx, volumeID)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when inspecting a volume", func() {
			It("returns the volume", func() {
				opts := &types.VolumeInspectOpts{
					Attachments: true,
				}
				volume, err := libsClient.Storage().VolumeInspect(ctx, volumeID, opts)
				Expect(err).ToNot(HaveOccurred())
				Expect(volume.Name).To(Equal(volumeName))
				Expect(volume.Size).To(Equal(int64(8)))
			})
		})
	})

	Describe("Service Broker Integration", func() {
		Describe("Provision", func() {
			Context("when provisioning request is valid", func() {
				It("returns 200", func() {

				})
			})
		})
	})
})
