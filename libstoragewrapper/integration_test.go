package libstoragewrapper_test

import (
	"os"
	"strconv"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/server"
	"github.com/emccode/libstorage/api/types"

	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	Describe("Libstorage Client Integration", func() {
		var libsClient types.APIClient
		var instanceID, storagePool, volumeID, libStorageURI, serviceName string
		var size int64
		var err error

		BeforeSuite(func() {
			instanceID = os.Getenv("TEST_INSTANCE_ID") //InstanceID comes from CC we translate into VolumeName 3c653bce-8752-451b-96d9-a8a1a925b118
			Expect(instanceID).ToNot(BeEmpty())
			storagePool = os.Getenv("STORAGE_POOL_NAME")
			size, err = strconv.ParseInt(os.Getenv("TEST_SIZE"), 10, 64)
			Expect(err).ToNot(HaveOccurred())
			serviceName = os.Getenv("LIB_STOR_SERVICE")
			Expect(serviceName).ToNot(BeEmpty())
			libStorageURI = os.Getenv("LIBSTORAGE_URI")
			Expect(libStorageURI).ToNot(BeEmpty())

			libsClient = server.NewLibsClient()
			Expect(err).ToNot(HaveOccurred())
		})

		BeforeEach(func() {
			volumeRequest, err := utils.CreateVolumeRequest(instanceID, storagePool, int64(8))
			Expect(err).ToNot(HaveOccurred())
			volume, err := libstoragewrapper.CreateVolume(libsClient, serviceName, volumeRequest)
			Expect(err).ToNot(HaveOccurred())
			volumeID = volume.ID
		})

		AfterEach(func() {
			err := libstoragewrapper.RemoveVolume(libsClient, serviceName, volumeID)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when inspecting a volume", func() {
			It("returns the volume", func() {
				volume, err := libstoragewrapper.GetVolumeByID(libsClient, serviceName, volumeID)
				Expect(err).ToNot(HaveOccurred())
				parsedVolumeName, err := utils.CreateNameForVolume(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(volume.Name).To(Equal(parsedVolumeName))
				Expect(volume.Size).To(Equal(int64(size)))
			})
		})

		Context("When passing in an instanceID", func() {
			It("return a volume ID if instanceID exist", func() {
				getVolumeID, err := libstoragewrapper.GetVolumeID(libsClient, serviceName, instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(getVolumeID).To(Equal(volumeID))
			})
		})
	})
})
