package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"backend/routes"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter()
}

func TestIngestTelemetry_Success(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedMessage := "Telemetry recebida com sucesso"
	if response["message"] != expectedMessage {
		t.Errorf("mensagem incorreta: esperado=%v, recebido=%v", expectedMessage, response["message"])
	}

	if response["data"] == nil {
		t.Errorf("campo data não encontrado na resposta")
	}
}

// ❌ JSON inválido
func TestIngestTelemetry_InvalidJSON(t *testing.T) {
	// TODO: enviar JSON quebrado
	// TODO: esperar status 400
}

// ❌ device_id inválido
func TestIngestTelemetry_InvalidDeviceID(t *testing.T) {
	// TODO: device_id = 0
	// TODO: esperar status 400
}

// ❌ sensor.type vazio
func TestIngestTelemetry_EmptySensorType(t *testing.T) {
	// TODO: sensor.type = ""
	// TODO: esperar status 400
}

// ❌ sensor.unit vazio
func TestIngestTelemetry_EmptySensorUnit(t *testing.T) {
	// TODO: sensor.unit = ""
	// TODO: esperar status 400
}

// ❌ reading.value_type vazio
func TestIngestTelemetry_EmptyValueType(t *testing.T) {
	// TODO: value_type = ""
	// TODO: esperar status 400
}

// ❌ timestamp ausente ou inválido
func TestIngestTelemetry_InvalidTimestamp(t *testing.T) {
	// TODO: timestamp inválido ou omitido
	// TODO: esperar status 400
}

// ❌ value_type inválido
func TestIngestTelemetry_InvalidValueType(t *testing.T) {
	// TODO: value_type = "banana"
	// TODO: esperar status 400
}

// ❌ payload incompleto (sem sensor)
func TestIngestTelemetry_MissingSensor(t *testing.T) {
	// TODO: remover objeto sensor
	// TODO: esperar status 400
}

// ❌ payload incompleto (sem reading)
func TestIngestTelemetry_MissingReading(t *testing.T) {
	// TODO: remover objeto reading
	// TODO: esperar status 400
}