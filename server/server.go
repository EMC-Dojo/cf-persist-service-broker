package server

import (
  "net/http"
  "io/ioutil"
  "os"

  "github.com/gin-gonic/gin"
)

type Server struct {
}

func catalogHandler(c *gin.Context) {
  c.Status(http.StatusOK)
  c.File("templates/catalog.json")
}

func provisioningHandler(c *gin.Context) {
  c.JSON(http.StatusCreated, gin.H{})
}

func deprovisionHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}

func bindingHandler(c *gin.Context) {
  body, _ := ioutil.ReadFile("fixtures/create_binding_response.json")
  c.String(http.StatusCreated, string(body))
}

func unbindingHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}

func (s Server) Run(port string) {
  server := gin.Default()

  authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
    os.Getenv("BROKER_USERNAME"): os.Getenv("BROKER_PASSWORD"),
  }))

  authorized.GET("/v2/catalog", catalogHandler)
  authorized.PUT("/v2/service_instances/:instanceId", provisioningHandler)
  authorized.PUT("/v2/service_instances/:instanceId/service_bindings/:bindingId", bindingHandler)
  authorized.DELETE("/v2/service_instances/:instanceId/service_bindings/:bindingId", unbindingHandler)
  authorized.DELETE("/v2/service_instances/:instanceId", deprovisionHandler)

  server.Run(":" + port)
}