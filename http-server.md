# HTTP Go Language Test

Given only the Go standard library create a server that executes periodic health check for HTTP websites.

It is expected that a **go.mod** exists in the root directory and does not contain any external dependencies.

## Data Model

The following defines a **Health Check**

| Field Name | Type | Description |
| ---        | ---  | ---         |
| ID | string | Unique server side generated id |
| Status | string | Status text describing the last health check result |
| Code | int32 | Status code describing the last health check result |
| Endpoint | string | URL of the health check target |
| Checked | int64 | Unix epoch seconds of the last time the health check executed |
| Duration | string | Duration of the last healch check execution |
| Error | string | Any error message that occurred during a check |

## Command Line Arguments

* Must accept address to listen from the command line
  > `--bind=127.0.0.1:8080`

* Must be able to be configured for SSL via the command line (self signed certificates are adequate for this test)
  > `--ssl --sslcert=certificate.crt --sslkey=certificate.key`

* Must implement background HTTP health checks of provided endpoints on a regular basis configured
  > `--checkfrequency=30s`

* Must execute health checks concurrently.
* Must persist health checks between application restarts.
* Must implement HTTP API that validates methods and data for the following requests:

## API

The following API must be implemented.

### List Health Checks

Returns a list of health checks sorted by endpoint with paging support of 10 items per page.

Appropriate errors must be provided when:

* Health check fails to load from persistence

#### Request

```bash
curl http://127.0.0.1:8080/api/health/checks?page=1
```

#### Response

```json
{
    "items": [
        {
            "id": "3b447fdf-d2e9-42bd-adcf-77d147b8b4dc",
            "status": "Ok",
            "code": 200,
            "endpoint": "https://www.blizzard.com/en-us/",
            "checked": 1564065876,
            "duration": "127ms"
        }, {
            "id": "a85113c0-c5e4-4657-ba88-9df8befdbaa1",
            "status": "Not Found",
            "code": 404,
            "endpoint": "https://www.blizzard.com/en-us/gotest",
            "checked": 1564065876,
            "duration": "127ms"
        }, {
            "id": "1162d149-bd6f-4979-b9e5-7ee6f94f450a",
            "status": "Error",
            "code": 0,
            "endpoint": "https://gotest.blizzard.com",
            "error": "could not resolve host",
            "checked": 1564065876,
            "duration": "127ms"
        }
    ],
    "page": 1,
    "total": 3,
    "size": 10
}
```

### Get Health Check

Return a single health check

Appropriate errors must be provided when:

* Given the health check does not exist
* Health check fails to load from persistence

#### Request

```bash
curl http://127.0.0.1:8080/api/health/checks/3b447fdf-d2e9-42bd-adcf-77d147b8b4dc
```

#### Response

```json
{
    "id": "3b447fdf-d2e9-42bd-adcf-77d147b8b4dc",
    "status": "Ok",
    "code": 200,
    "endpoint": "https://www.blizzard.com/en-us/",
    "checked": 1564065876,
    "duration": "127ms"
}
```

### Create Health Check

Will create a health check, persist server side, and return the details.

Appropriate errors must be provided when:

* Given a health check with the same URL already exists
* Given the endpoint is blank
* Given the endpoint is not a valid URL
* Health Check persistence fails

#### Request

```bash
curl -X POST http://127.0.0.1:8080/api/health/checks -d '{ "endpoint": "http://example.com" }'
```

#### Response

```json
{
    "id": "94a1d1e8-6e44-409e-9cb4-7bfcac2de1ae",
    "endpoint": "http://example.com"
}
```

### Execute a Health Check

This will execute a health check, with a timeout provided in the query string

#### Request

```bash
curl -X POST http://127.0.0.1:8080/api/health/checks/94a1d1e8-6e44-409e-9cb4-7bfcac2de1ae/try?timeout=10s
```

#### Response

```json
{
    "id": "94a1d1e8-6e44-409e-9cb4-7bfcac2de1ae",
    "status": "Ok",
    "code": 200,
    "endpoint": "http://example.com",
    "checked": 1564065975,
    "duration": "127ms"
}
```

### Delete a Health Check

This will delete a health check and remove it from the persisted data.

#### Request

```bash
curl -X DELETE http://127.0.0.1:8080/api/health/checks/94a1d1e8-6e44-409e-9cb4-7bfcac2de1ae
```

#### Response

No response body is expected.
