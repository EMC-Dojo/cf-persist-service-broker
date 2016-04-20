package server

import (
  "io/ioutil"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/EMC-CMD/cf-persist-service-broker/storage"
)

func CatalogHandler(c *gin.Context) {
  c.Status(http.StatusOK)
  c.File("templates/catalog.json")
}

func ProvisioningHandler(c *gin.Context) {
  _, err := CreateVolume(&storage.ScaleIODriver{})
  if err != nil {
    c.JSON(422, gin.H{})
  }
  c.JSON(http.StatusCreated, gin.H{})
}

func CreateVolume(driver storage.StorageDriver) (*storage.Volume, error) {
  return driver.VolumeCreate(storage.Context{}, "", &storage.VolumeCreateOpts{})
}

func DeprovisionHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}

func BindingHandler(c *gin.Context) {
  body, _ := ioutil.ReadFile("fixtures/create_binding_response.json")
  c.String(http.StatusCreated, string(body))
}

func UnbindingHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}
