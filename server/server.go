package server

import (
  "os"

  "github.com/gin-gonic/gin"
)

type Server struct {
}

func (s Server) Run(port string) {
  server := gin.Default()

  authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
    os.Getenv("BROKER_USERNAME"): os.Getenv("BROKER_PASSWORD"),
  }))

  authorized.GET("/v2/catalog", CatalogHandler)
  authorized.PUT("/v2/service_instances/:instanceId", ProvisioningHandler)
  authorized.PUT("/v2/service_instances/:instanceId/service_bindings/:bindingId", BindingHandler)
  authorized.DELETE("/v2/service_instances/:instanceId/service_bindings/:bindingId", UnbindingHandler)
  authorized.DELETE("/v2/service_instances/:instanceId", DeprovisionHandler)

  server.Run(":" + port)
}