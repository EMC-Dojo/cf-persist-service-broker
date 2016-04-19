package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CatalogHandler(c *gin.Context) {
	c.Status(http.StatusOK)
	c.File("templates/catalog.json")
}

func ProvisioningHandler(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{})
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

func main() {
	server := gin.Default()

	authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
		os.Getenv("BROKER_USERNAME"): os.Getenv("BROKER_PASSWORD"),
	}))

	authorized.GET("/v2/catalog", CatalogHandler)
	authorized.PUT("/v2/service_instances/:instanceId", ProvisioningHandler)
	authorized.PUT("/v2/service_instances/:instanceId/service_bindings/:bindingId", BindingHandler)
	authorized.DELETE("/v2/service_instances/:instanceId/service_bindings/:bindingId", UnbindingHandler)
	authorized.DELETE("/v2/service_instances/:instanceId", DeprovisionHandler)

	server.Run(":" + os.Getenv("PORT"))
}
