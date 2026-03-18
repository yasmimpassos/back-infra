package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/models"
)

func validateRequiredField(c *gin.Context, value string, fieldName string) bool {
	if value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fieldName + " é obrigatório",
		})
		return false
	}
	return true
}

func IngestTelemetry(c *gin.Context) {
	var message models.TelemetryMessage

	err := c.ShouldBindJSON(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Payload inválido",
		})
		return
	}

	if message.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "device_id é obrigatório",
		})
		return
	}

	if message.Sensor == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sensor é obrigatório",
		})
		return
	}

	if message.Reading == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "reading é obrigatório",
		})
		return
	}

	if !validateRequiredField(c, message.Sensor.Type, "sensor.type") {
		return
	}

	if !validateRequiredField(c, message.Sensor.Unit, "sensor.unit") {
		return
	}

	if !validateRequiredField(c, message.Reading.Type, "reading.value_type") {
		return
	}

	if message.Timestamp.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "timestamp é obrigatório",
		})
		return
	}

	if message.Reading.Type != "analog" && message.Reading.Type != "discrete" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "value_type deve ser 'analog' ou 'discrete'",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Telemetry recebida com sucesso",
		"data":    message,
	})
}