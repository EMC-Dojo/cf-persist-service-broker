package libstoragewrapper_test

//
// import (
// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
//
// 	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
// 	"github.com/EMC-Dojo/cf-persist-service-broker/mocks"
// )
//
// var _ = Describe("Controller", func() {
// 	var (
// 		t           GinkgoTestReporter
// 		mockCtrl    *gomock.Controller
// 		mockClient  *mocks.MockAPIClient
// 		serviceName string
// 	)
// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(t)
// 		mockClient = mocks.NewMockAPIClient(mockCtrl)
// 		mockContext = mocks.NewMockContext(mockCtrl)
// 		serviceName = "mockService"
// 	})
//
// 	AfterEach(func() {
// 		mockCtrl.Finish()
// 	})
//
// 	Describe("RemoveVolume", func() {
// 		Context("when deletion succeeded", func() {
// 			It("returns no error", func() {
// 				mockClient.EXPECT().VolumeRemove(serviceName, "volumeID").Return(nil)
//
// 				err := libstoragewrapper.RemoveVolume(mockClient, serviceName, "volumeID")
// 				Expect(err).ToNot(HaveOccurred())
// 			})
// 		})
// 	})
// })
