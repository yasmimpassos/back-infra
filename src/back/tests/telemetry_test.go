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
	router := setupTestRouter()

	payload := `{
		"device_id":
	}`

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "Payload inválido"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}

// ❌ device_id inválido
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "device_id é obrigatório"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}

// ❌ sensor.type vazio
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "sensor.type é obrigatório"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}

// ❌ sensor.unit vazio
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "sensor.unit é obrigatório"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}

// ❌ reading.value_type vazio
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "reading.value_type é obrigatório"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}

// ❌ timestamp ausente ou inválido
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

	req, _ := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	expected := "Payload inválido"
	if response["error"] != expected {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expected, response["error"])
	}
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

	req, _ := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	expected := "Payload inválido"
	if response["error"] != expected {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expected, response["error"])
	}
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

	req, _ := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	expected := "timestamp é obrigatório"
	if response["error"] != expected {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expected, response["error"])
	}
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

	req, _ := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	expected := "Payload inválido"
	if response["error"] != expected {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expected, response["error"])
	}
}

// ❌ value_type inválido
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "value_type deve ser 'analog' ou 'discrete'"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}
                          
// ❌ payload incompleto (sem sensor)
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "sensor é obrigatório"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}

// ❌ payload incompleto (sem reading)
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

	req, err := http.NewRequest("POST", "/telemetry", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("erro ao fazer parse do JSON: %v", err)
	}

	expectedError := "reading é obrigatório"
	if response["error"] != expectedError {
		t.Errorf("erro incorreto: esperado=%v, recebido=%v", expectedError, response["error"])
	}
}