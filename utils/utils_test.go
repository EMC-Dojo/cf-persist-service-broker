package utils_test

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils/Utils", func() {
	var instanceID, parsedInstanceID, storagePool, libstorage_uri string
	var size int64
	var err error

	BeforeSuite(func() {
		instanceID = os.Getenv("TEST_INSTANCE_ID") //3c653bce-8752-451b-96d9-a8a1a925b118
		Expect(instanceID).ToNot(BeEmpty())
		parsedInstanceID = os.Getenv("PARSED_INSTANCE_ID") //3c653bce8752451b96d9a8a1a925b118
		Expect(parsedInstanceID).ToNot(BeEmpty())
		storagePool = os.Getenv("STORAGE_POOL_NAME")
		Expect(storagePool).ToNot(BeEmpty())
		libstorage_uri = os.Getenv("LIBSTORAGE_URI")
		Expect(libstorage_uri).ToNot(BeEmpty())
		size, err = strconv.ParseInt(os.Getenv("TEST_SIZE"), 10, 64)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Generate Volume Name", func() {
		It("Should Create Volume Name using instanceId and serviceInstance", func() {

			cooked, err := utils.CreateNameForVolume(instanceID)
			Expect(err).ToNot(HaveOccurred())
			Expect(cooked).To(Equal(parsedInstanceID))
		})
		It("Should return an error if it is longer than 32 characters after removing all hyphens", func() {
			_, err := utils.CreateNameForVolume(instanceID + "1234567890abcdefghijklmnopqrstuv")
			Expect(err).To(MatchError("Volume name cannot exceed 32 characters when all hyphens are removed."))
		})
		It("Should leave the volume name alone if it is shorter than 31 characters", func() {
			cooked, err := utils.CreateNameForVolume("123-456-789")
			Expect(err).ToNot(HaveOccurred())
			Expect(cooked).To(Equal("123-456-789"))
		})
	})

	Describe("CreateVolumeRequest", func() {
		Context("When scaleio", func() {
			It("Should Create Volume Options using a serviceInstance", func() {

				volumeCreateRequest, err := utils.CreateVolumeRequest(instanceID, storagePool, size)
				Expect(err).ToNot(HaveOccurred())

				Expect(volumeCreateRequest.Name).To(Equal(parsedInstanceID))
				Expect(*volumeCreateRequest.AvailabilityZone).To(Equal(""))
				Expect(*volumeCreateRequest.IOPS).To(Equal(int64(0)))
				Expect(*volumeCreateRequest.Size).To(Equal(size))
				Expect(*volumeCreateRequest.Type).To(Equal(storagePool))
			})
		})
	})

	Describe("CreatePlanIDString and UnmarshalPlanID", func() {
		Context("When given libstorage uri and service name", func() {
			It("contructs plan id as a json base on those fields", func() {
				planIDJson, err := utils.CreatePlanIDString("isilonservice", libstorage_uri)
				Expect(err).ToNot(HaveOccurred())

				planID, err := utils.UnmarshalPlanID(planIDJson)
				Expect(err).ToNot(HaveOccurred())
				Expect(planID).To(Equal(model.PlanID{
					LibsHostName:    libstorage_uri,
					LibsServiceName: "isilonservice",
				}))
			})
		})
	})

	Describe("FileExists", func() {
		Context("File exists after creation and does not exist after deletion", func() {
			It("return true or false", func() {
				path := "test.test"
				Expect(utils.FileExists(path)).To(BeFalse())

				err := ioutil.WriteFile(path, []byte("abc"), 0644)
				Expect(err).ToNot(HaveOccurred())
				Expect(utils.FileExists(path)).To(BeTrue())

				err = os.Remove(path)
				Expect(err).ToNot(HaveOccurred())
				Expect(utils.FileExists(path)).To(BeFalse())
			})
		})
	})
})
