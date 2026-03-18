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

func performRequest(router *gin.Engine, method, path, payload string) (*httptest.ResponseRecorder, map[string]interface{}) {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return w, response
}

func assertError(t *testing.T, w *httptest.ResponseRecorder, response map[string]interface{}, expectedStatus int, expectedMessage string) {
	if w.Code != expectedStatus {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", expectedStatus, w.Code)
	}

	if response["error"] != expectedMessage {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedMessage, response["error"])
	}
}

func assertSuccess(t *testing.T, w *httptest.ResponseRecorder, response map[string]interface{}) {
	if w.Code != http.StatusOK {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	expected := "Telemetry recebida com sucesso"
	if response["message"] != expected {
		t.Errorf("mensagem incorreta: esperado=%v, recebido=%v", expected, response["message"])
	}

	if response["data"] == nil {
		t.Errorf("campo data não encontrado")
	}
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

	w, response := performRequest(router, "POST", "/telemetry", payload)

	assertSuccess(t, w, response)
}

func TestIngestTelemetry_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id":
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)

	assertError(t, w, response, http.StatusBadRequest, "Payload inválido")
}

func TestIngestTelemetry_InvalidDeviceID(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 0,
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

	w, response := performRequest(router, "POST", "/telemetry", payload)

	assertError(t, w, response, http.StatusBadRequest, "device_id é obrigatório")
}

func TestIngestTelemetry_EmptySensorType(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)

	assertError(t, w, response, http.StatusBadRequest, "sensor.type é obrigatório")
}

func TestIngestTelemetry_EmptySensorUnit(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": ""
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)

	assertError(t, w, response, http.StatusBadRequest, "sensor.unit é obrigatório")
}

func TestIngestTelemetry_EmptyValueType(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)

	assertError(t, w, response, http.StatusBadRequest, "reading.value_type é obrigatório")
}

func TestIngestTelemetry_InvalidTimestampFormat(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "data-errada",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "Payload inválido")
}

func TestIngestTelemetry_InvalidTimestampType(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": 123456,
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "Payload inválido")
}

func TestIngestTelemetry_MissingTimestamp(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "timestamp é obrigatório")
}

func TestIngestTelemetry_EmptyTimestamp(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "Payload inválido")
}

func TestIngestTelemetry_InvalidValueType(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analogo",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "value_type deve ser 'analog' ou 'discrete'")
}

func TestIngestTelemetry_MissingSensor(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "sensor é obrigatório")
}

func TestIngestTelemetry_MissingReading(t *testing.T) {
	router := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		}
	}`

	w, response := performRequest(router, "POST", "/telemetry", payload)
	assertError(t, w, response, http.StatusBadRequest, "reading é obrigatório")
}