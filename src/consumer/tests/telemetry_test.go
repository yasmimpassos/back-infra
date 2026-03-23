package rabbitmq

import (
	"errors"
	"testing"
	"time"
	"consumer/rabbitmq"
	"consumer/models"
)

func validTelemetryPayload() []byte {
	return []byte(`{
		"device_id": 1,
		"timestamp": "2026-03-20T13:15:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 25.3
		}
	}`)
}

func invalidTelemetryPayload() []byte {
	return []byte(`{
		"device_id":
	}`)
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("esperava nil, recebeu erro: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

func assertTelemetryFields(t *testing.T, got models.TelemetryMessage) {
	t.Helper()

	if got.Id != 1 {
		t.Errorf("device_id incorreto: esperado=%d, recebido=%d", 1, got.Id)
	}

	expectedTime := time.Date(2026, 3, 20, 13, 15, 0, 0, time.UTC)
	if !got.Timestamp.Equal(expectedTime) {
		t.Errorf("timestamp incorreto: esperado=%v, recebido=%v", expectedTime, got.Timestamp)
	}

	if got.Sensor == nil {
		t.Fatal("sensor não deveria ser nil")
	}

	if got.Sensor.Type != "temperature" {
		t.Errorf("sensor.type incorreto: esperado=%v, recebido=%v", "temperature", got.Sensor.Type)
	}

	if got.Sensor.Unit != "celsius" {
		t.Errorf("sensor.unit incorreto: esperado=%v, recebido=%v", "celsius", got.Sensor.Unit)
	}

	if got.Reading == nil {
		t.Fatal("reading não deveria ser nil")
	}

	if got.Reading.Type != "analog" {
		t.Errorf("reading.value_type incorreto: esperado=%v, recebido=%v", "analog", got.Reading.Type)
	}

	if got.Reading.Value != 25.3 {
		t.Errorf("reading.value incorreto: esperado=%v, recebido=%v", 25.3, got.Reading.Value)
	}
}

func TestProcessTelemetryMessage_Success(t *testing.T) {
	payload := validTelemetryPayload()

	var saved models.TelemetryMessage
	saveFn := func(msg models.TelemetryMessage) error {
		saved = msg
		return nil
	}

	err := rabbitmq.ProcessTelemetryMessage(payload, saveFn)
	assertNoError(t, err)
	assertTelemetryFields(t, saved)
}

func TestProcessTelemetryMessage_InvalidJSON(t *testing.T) {
	payload := invalidTelemetryPayload()

	called := false
	saveFn := func(msg models.TelemetryMessage) error {
		called = true
		return nil
	}

	err := rabbitmq.ProcessTelemetryMessage(payload, saveFn)
	assertError(t, err)

	if called {
		t.Fatal("saveFn não deveria ter sido chamada para JSON inválido")
	}
}

func TestProcessTelemetryMessage_SaveError(t *testing.T) {
	payload := validTelemetryPayload()

	expectedErr := errors.New("erro ao salvar no banco")
	saveFn := func(msg models.TelemetryMessage) error {
		return expectedErr
	}

	err := rabbitmq.ProcessTelemetryMessage(payload, saveFn)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("erro incorreto: esperado=%v, recebido=%v", expectedErr, err)
	}
}