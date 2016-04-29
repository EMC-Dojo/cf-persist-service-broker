package server

import (
	"fmt"

	"github.com/EMC-CMD/cf-persist-service-broker/storage"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		t                 GinkgoTestReporter
		mockCtrl          *gomock.Controller
		mockStorageDriver *storage.MockStorageDriver
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(t)
		mockStorageDriver = storage.NewMockStorageDriver(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("CreateVolume", func() {
		Context("when provision succeeded", func() {
			It("returns volumes", func() {
				mockStorageDriver.EXPECT().VolumeCreate(storage.Context{}, "", &storage.VolumeCreateOpts{}).Return(&storage.Volume{}, nil)

				volume, err := CreateVolume(mockStorageDriver)
				Expect(err).ToNot(HaveOccurred())
				Expect(*volume).To(Equal(storage.Volume{}))
			})
		})
	})

	Describe("RemoveVolume", func() {
		Context("when deletion succeeded", func() {
			It("returns no error", func() {
				mockStorageDriver.EXPECT().VolumeRemove(storage.Context{}, "", &storage.VolumeCreateOpts{}).Return(nil)

				err := RemoveVolume(mockStorageDriver)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
