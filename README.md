<h1 align="center">
  VersityGW Webhook Pulsar Proxy
</h1>
<h2 align="center">
    A lightweight webhook proxy that forwards VersityGW S3 events to Apache Pulsar
</h2>

<div align="center">

[![Alpineworks][alpineworks-badge]][for-the-badge-link] [![Made With Go][made-with-go-badge]][for-the-badge-link]

</div>

---

## Overview

VersityGW Webhook Pulsar Proxy is a high-performance HTTP service that receives S3 event notifications from VersityGW and forwards them to Apache Pulsar topics. It provides observability through OpenTelemetry metrics and tracing, making it suitable for production deployments.

## Usage

<details open>
<summary>Docker Compose (Recommended)</summary>

```bash
# Start the full stack (proxy + Pulsar + monitoring)
docker-compose up -d

# Send a test S3 event
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "eventSource": "aws:s3",
    "eventName": "s3:ObjectCreated:Put",
    "s3": {
      "bucket": {"name": "test-bucket"},
      "object": {"key": "test-object.txt"}
    }
  }'
```

</details>

## Configuration

All configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `error` | Log level (debug, info, warn, error) |
| `SERVER_PORT` | `8080` | HTTP server port |
| `PULSAR_URL` | `pulsar://localhost:6650` | Pulsar broker URL |
| `PULSAR_TOPIC` | `s3-events` | Target Pulsar topic |
| `PULSAR_PRODUCE_TIMEOUT_SECONDS` | `5s` | Message production timeout |
| `METRICS_ENABLED` | `true` | Enable Prometheus metrics |
| `METRICS_PORT` | `8081` | Metrics server port |
| `TRACING_ENABLED` | `false` | Enable OpenTelemetry tracing |
| `TRACING_SAMPLERATE` | `0.01` | Trace sampling rate |
| `TRACING_SERVICE` | `versitygw-webhook-pulsar-proxy` | Service name for tracing |

## API

### POST /webhook

Accepts VersityGW S3 event notifications and forwards them to Pulsar.

**Request Body:** S3 EventRecord JSON
**Response:** 
- `201 Created` - Event successfully forwarded
- `400 Bad Request` - Invalid JSON payload
- `500 Internal Server Error` - Pulsar delivery failure

**Example Response:**
```json
{
  "message": "event forwarded to pulsar",
  "code": 201
}
```

<!-- Reference Variables -->

<!-- Badges -->
[alpineworks-badge]: .github/images/alpine-works.svg
[made-with-go-badge]: .github/images/made-with-go.svg

<!-- Links -->
[for-the-badge-link]: https://forthebadge.com