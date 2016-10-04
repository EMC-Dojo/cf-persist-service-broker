package server

import (
	"fmt"
	"os"

	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/types"
)

func AddCatalogPlanIDs(catalogPlans []model.Plan, libstorageHost string) ([]model.Plan, error) {
	for i, plan := range catalogPlans {
		newPlanID, err := utils.CreatePlanIDString(plan.Name, libstorageHost)
		if err != nil {
			return catalogPlans, err
		}
		catalogPlans[i].ID = newPlanID
	}
	return catalogPlans, nil
}

func DoPlansExistInLibstorage(catalogPlans []model.Plan, libstorageServices map[string]*types.ServiceInfo) bool {
	dict := map[string]bool{}

	for _, libService := range libstorageServices {
		dict[libService.Name] = true
	}

	for _, plan := range catalogPlans {
		if !dict[plan.Name] {
			return false
		}
	}

	return true
}

func VerifyEnvironmentVariable(envs []string) error {
	for _, item := range envs {
		if os.Getenv(item) == "" {
			return fmt.Errorf("Environment variable %s is expected", item)
		}
	}
	return nil
}
