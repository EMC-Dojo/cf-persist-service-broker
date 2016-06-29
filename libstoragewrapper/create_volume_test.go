package libstoragewrapper_test

//
// import (
// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
//
// 	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
// 	"github.com/EMC-Dojo/cf-persist-service-broker/mocks"
// 	"github.com/emccode/libstorage/api/types"
// )
//
// var _ = Describe("Controller", func() {
// 	var (
// 		t           GinkgoTestReporter
// 		mockCtrl    *gomock.Controller
// 		mockClient  *mocks.MockAPIClient
// 		mockContext *mocks.MockContext
// 	)
//
// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(t)
// 		mockClient = mocks.NewMockAPIClient(mockCtrl)
// 	})
//
// 	AfterEach(func() {
// 		mockCtrl.Finish()
// 	})
//
// 	Describe("CreateVolume", func() {
// 		Context("when provision succeeded", func() {
// 			It("returns volumes", func() {
// 				serviceName := "mockService"
// 				size := int64(8)
// 				iops := int64(100)
// 				volumeType := "type"
// 				availabilityZone := "az"
// 				volumeReturned := &types.Volume{
// 					Name:             "mock",
// 					Size:             size,
// 					AvailabilityZone: availabilityZone,
// 					IOPS:             iops,
// 					Type:             volumeType,
// 				}
// 				volumeCreateRequest := types.VolumeCreateRequest{
// 					Name:             "mock",
// 					AvailabilityZone: &availabilityZone,
// 					IOPS:             &iops,
// 					Type:             &volumeType,
// 					Size:             &size,
// 				}
//
// 				mockClient.EXPECT().VolumeCreate(serviceName, &volumeCreateRequest).Return(volumeReturned, nil)
//
// 				volume, err := libstoragewrapper.CreateVolume(mockClient, serviceName, volumeCreateRequest)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(volume).To(Equal(volumeReturned))
// 			})
// 		})
// 	})
// })
