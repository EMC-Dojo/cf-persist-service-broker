package server_test

import (
	// . "github.com/EMC-Dojo/cf-persist-service-broker"

	"fmt"
	"os"

	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/server"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server/ServerHelper", func() {

	Describe("ModifyPlanID", func() {
		var libstorageHost string
		BeforeEach(func() {
			libstorageHost = os.Getenv("LIBSTORAGE_URI")
			Expect(libstorageHost).ToNot(BeEmpty())
		})
		Context("Modify Plan ID in config file", func() {
			It("appends libstorage host to plan's name", func() {
				catalogPlans := []model.Plan{
					model.Plan{
						Name: "isilonservice",
					},
					model.Plan{
						Name: "scaleioservice",
					},
				}
				isilonPlanID, err := utils.CreatePlanIDString("isilonservice", libstorageHost)
				Expect(err).ToNot(HaveOccurred())
				scaleioPlanID, err := utils.CreatePlanIDString("scaleioservice", libstorageHost)
				Expect(err).ToNot(HaveOccurred())

				catalogPlans, err = server.AddCatalogPlanIDs(catalogPlans, libstorageHost)
				Expect(err).ToNot(HaveOccurred())
				Expect(catalogPlans).To(Equal(
					[]model.Plan{
						model.Plan{
							ID:   isilonPlanID,
							Name: "isilonservice",
						},
						model.Plan{
							ID:   scaleioPlanID,
							Name: "scaleioservice",
						},
					},
				))
			})
		})

		Describe("DoPlansExistInLibstorage", func() {
			Context("If all the plans from catalog exist in libstorage server", func() {
				It("return true", func() {
					catalogPlans := []model.Plan{
						model.Plan{
							Name: "isilonservice",
						},
						model.Plan{
							Name: "scaleioservice",
						},
					}
					libstorageServices := map[string]*types.ServiceInfo{
						"1": &types.ServiceInfo{
							Name: "isilonservice",
						},
						"2": &types.ServiceInfo{
							Name: "scaleioservice",
						},
					}
					exist := server.DoPlansExistInLibstorage(catalogPlans, libstorageServices)
					Expect(exist).To(BeTrue())
				})
			})

			Context("If plans does not exist in libstorage server", func() {
				It("return false", func() {
					catalogPlans := []model.Plan{
						model.Plan{
							Name: "isilonservice",
						},
						model.Plan{
							Name: "scaleioservice",
						},
						model.Plan{
							Name: "fakeservice",
						},
					}
					libstorageServices := map[string]*types.ServiceInfo{
						"1": &types.ServiceInfo{
							Name: "isilonservice",
						},
						"2": &types.ServiceInfo{
							Name: "scaleioservice",
						},
					}
					exist := server.DoPlansExistInLibstorage(catalogPlans, libstorageServices)
					Expect(exist).To(BeFalse())
				})
			})
		})

		Describe("VerifyEnvironmentVariable", func() {
			Context("If one of the environment is empty", func() {
				It("returns an environment not exist err", func() {
					notexistenv := "NOTEXISTENV"
					existenv := "EXISTENV"
					os.Unsetenv(notexistenv)
					os.Setenv(existenv, "exist")
					err := server.VerifyEnvironmentVariable([]string{existenv, notexistenv})
					Expect(err).To(MatchError(fmt.Sprintf("Environment variable %s is expected", notexistenv)))
				})
			})

			Context("If all environment vars are set", func() {
				It("does not return err", func() {
					existenva := "EXISTENVA"
					existenvb := "EXISTENVB"
					os.Setenv(existenva, "exista")
					os.Setenv(existenvb, "existb")
					err := server.VerifyEnvironmentVariable([]string{existenva, existenvb})
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
