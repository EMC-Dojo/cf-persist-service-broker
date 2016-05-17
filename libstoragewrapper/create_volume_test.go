package libstoragewrapper_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/EMC-CMD/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-CMD/cf-persist-service-broker/mocks"
	"github.com/emccode/libstorage/api/types"
)

var _ = Describe("Controller", func() {
	var (
		t                 GinkgoTestReporter
		mockCtrl          *gomock.Controller
		mockClient        *mocks.MockClient
		mockContext       *mocks.MockContext
		mockStorageDriver *mocks.MockStorageDriver
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(t)
		mockClient = mocks.NewMockClient(mockCtrl)
		mockContext = mocks.NewMockContext(mockCtrl)
		mockStorageDriver = mocks.NewMockStorageDriver(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("CreateVolume", func() {
		Context("when provision succeeded", func() {
			It("returns volumes", func() {
				size := int64(8)
				iops := int64(100)
				volumeType := "type"
				availabilityZone := "az"
				volumeReturned := &types.Volume{
					Name:             "mock",
					Size:             size,
					AvailabilityZone: availabilityZone,
					IOPS:             iops,
					Type:             volumeType,
				}
				volumeCreateOpts := types.VolumeCreateOpts{
					AvailabilityZone: &availabilityZone,
					IOPS:             &iops,
					Type:             &volumeType,
					Size:             &size,
				}

				mockClient.EXPECT().Storage().Return(mockStorageDriver)
				mockStorageDriver.EXPECT().VolumeCreate(mockContext, "mock", &volumeCreateOpts).Return(volumeReturned, nil)

				volume, err := libstoragewrapper.CreateVolume(mockClient, mockContext, "mock", volumeCreateOpts)
				Expect(err).ToNot(HaveOccurred())
				Expect(volume).To(Equal(volumeReturned))
			})
		})
	})
})
