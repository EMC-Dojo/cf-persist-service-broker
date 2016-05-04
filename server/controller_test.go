package server

import (
  "fmt"

  "github.com/golang/mock/gomock"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/emccode/libstorage/api/types"
  "github.com/EMC-CMD/cf-persist-service-broker/mocks"
)

type GinkgoTestReporter struct{}

func (g GinkgoTestReporter) Errorf(format string, args ...interface{}) {
  Fail(fmt.Sprintf(format, args))
}

func (g GinkgoTestReporter) Fatalf(format string, args ...interface{}) {
  Fail(fmt.Sprintf(format, args))
}

var _ = Describe("Controller", func() {
  var (
    t GinkgoTestReporter
    mockCtrl   *gomock.Controller
    mockClient *mocks.MockClient
    mockContext *mocks.MockContext
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
          Name: "mock",
          Size: size,
          AvailabilityZone: availabilityZone,
          IOPS: iops,
          Type: volumeType,
        }
        volumeCreateOpts := &types.VolumeCreateOpts{
          AvailabilityZone: &availabilityZone,
          IOPS: &iops,
          Type: &volumeType,
          Size: &size,
        }

        mockClient.EXPECT().Storage().Return(mockStorageDriver)
        mockStorageDriver.EXPECT().VolumeCreate(mockContext, "mock", volumeCreateOpts).Return(volumeReturned, nil)

        volume, err := CreateVolume(mockClient, mockContext, "mock", availabilityZone, volumeType, iops, size)
        Expect(err).ToNot(HaveOccurred())
        Expect(volume).To(Equal(volumeReturned))
      })
    })
  })

  Describe("RemoveVolume", func() {
    Context("when deletion succeeded", func() {
      It("returns no error", func() {
        mockClient.EXPECT().Storage().Return(mockStorageDriver)
        mockStorageDriver.EXPECT().VolumeRemove(mockContext, "volumeID", nil).Return(nil)

        err := RemoveVolume(mockClient, mockContext, "volumeID")
        Expect(err).ToNot(HaveOccurred())
      })
    })
  })
})
