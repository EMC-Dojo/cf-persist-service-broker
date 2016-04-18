package main

import (
	"net/http"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func CatalogHandler(c *gin.Context)  {
	c.Status(http.StatusOK)
	c.File("seeds/catalog.json")
}

func ProvisioningHandler(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{})
}

func DeprovisionHandler(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{})
}

func BindingHandler(c *gin.Context)  {
	body, _ := ioutil.ReadFile("fixtures/create_binding_response.json")
	c.String(http.StatusCreated, string(body))
}

func UnbindingHandler(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{})
}

func main() {
	server := gin.Default()

	server.GET("/v2/catalog", CatalogHandler)
	server.PUT("/v2/service_instances/:instanceId", ProvisioningHandler)
	server.PUT("/v2/service_instances/:instanceId/service_bindings/:bindingId", BindingHandler)
	server.DELETE("/v2/service_instances/:instanceId/service_bindings/:bindingId", UnbindingHandler)
	server.DELETE("/v2/service_instances/:instanceId", DeprovisionHandler)

	server.Run(":"+os.Getenv("PORT"))
}
