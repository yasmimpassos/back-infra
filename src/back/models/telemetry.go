package main

import "time"

type SensorInfo struct {
	Type string `json:"type"`
	Unit string `json:"unit"`
}

type SensorReading struct {
	Type  string  `json:"value_type"`
	Value float64 `json:"value"`
}

type TelemetryMessage struct {
	Id        int           `json:"device_id"`
	Timestamp time.Time     `json:"timestamp"`
	Sensor    SensorInfo    `json:"sensor"`
	Reading   SensorReading `json:"reading"`
}