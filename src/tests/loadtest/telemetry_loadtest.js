import http from "k6/http";
import { check } from "k6";
import { Counter, Rate, Trend } from "k6/metrics";

export const options = {
  scenarios: {
    telemetry_constant_rate: {
      executor: "constant-arrival-rate",
      rate: 20,
      timeUnit: "1s",
      duration: "2m",
      preAllocatedVUs: 20,
      maxVUs: 50,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
    success_rate: ["rate>0.99"],
  },
};

const baseUrl = __ENV.BASE_URL || "http://localhost:8080";

const telemetryRequests = new Counter("telemetry_requests_total");
const successRate = new Rate("success_rate");
const requestDuration = new Trend("request_duration_ms");
const payloadBytes = new Trend("payload_bytes");

function randomSensor() {
  const sensors = [
    { type: "temperature", unit: "celsius", value_type: "analog", value: (20 + Math.random() * 15).toFixed(2) },
    { type: "humidity", unit: "percent", value_type: "analog", value: (40 + Math.random() * 40).toFixed(2) },
    { type: "presence", unit: "boolean", value_type: "discrete", value: Math.random() > 0.5 ? 1 : 0 },
    { type: "vibration", unit: "m/s2", value_type: "analog", value: (Math.random() * 10).toFixed(2) },
    { type: "luminosity", unit: "lux", value_type: "analog", value: (Math.random() * 1000).toFixed(2) },
    { type: "level", unit: "percent", value_type: "analog", value: (Math.random() * 100).toFixed(2) },
  ];

  return sensors[Math.floor(Math.random() * sensors.length)];
}

export default function () {
  const sensor = randomSensor();

  const payload = JSON.stringify({
    device_id: Math.floor(Math.random() * 1000) + 1,
    timestamp: new Date().toISOString(),
    sensor: {
      type: sensor.type,
      unit: sensor.unit,
    },
    reading: {
      value_type: sensor.value_type,
      value: Number(sensor.value),
    },
  });

  const res = http.post(`${baseUrl}/telemetry`, payload, {
    headers: { "Content-Type": "application/json" },
  });

  telemetryRequests.add(1);
  successRate.add(res.status === 200);
  requestDuration.add(res.timings.duration);
  payloadBytes.add(payload.length);

  check(res, {
    "status is 200": (r) => r.status === 200,
    "response has success message": (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.message === "Telemetry enviada com sucesso";
      } catch {
        return false;
      }
    },
  });
}