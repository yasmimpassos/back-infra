package db

import (
	"consumer/models"
)

func InsertTelemetry(t models.TelemetryMessage) error {
	query := `
		INSERT INTO telemetry (
			device_id,
			timestamp,
			sensor_type,
			sensor_unit,
			reading_type,
			value
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := DB.Exec(
		query,
		t.Id,
		t.Timestamp,
		t.Sensor.Type,
		t.Sensor.Unit,
		t.Reading.Type,
		t.Reading.Value,
	)

	return err
}