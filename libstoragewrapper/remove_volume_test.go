package libstoragewrapper_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/mocks"
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

	Describe("RemoveVolume", func() {
		Context("when deletion succeeded", func() {
			It("returns no error", func() {
				mockClient.EXPECT().Storage().Return(mockStorageDriver)
				mockStorageDriver.EXPECT().VolumeRemove(mockContext, "volumeID", nil).Return(nil)

				err := libstoragewrapper.RemoveVolume(mockClient, mockContext, "volumeID")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
