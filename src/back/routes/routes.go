package routes

import (
	"github.com/gin-gonic/gin"
	"backend/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/telemetry", handlers.IngestTelemetry)

	return r
}