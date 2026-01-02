# debug-server

Generic HTTP debug server that logs every incoming request (as pretty-printed JSON) to stdout and lets you control the response status code and body via query parameters.

## What it does

- Listens on `ADDR:PORT` (defaults: `0.0.0.0:8080`)
- Logs each request with method, path, query, headers, and body
- Responds on any path with a customizable status code/body/content-type
- Exposes a simple readiness endpoint: `GET /ready` returns `200`

## Run with Docker

```bash
docker run --rm -p 8080:8080 jlmbrt/debug-server
```

Configure bind address/port:

```bash
docker run --rm -e ADDR=0.0.0.0 -e PORT=8080 -p 8080:8080 jlmbrt/debug-server
```

## Run locally (Go)

```bash
ADDR=0.0.0.0 PORT=8080 go run .
```

## Endpoints

### `GET /ready`

Always returns `200`.

### `ANY /<anything>`

Logs the request and responds based on these query parameters:

- `response-code`: HTTP status code (defaults to `200`; invalid values fall back to `200`)
- `response-body`: response body (defaults to `{}`)
- `response-content-type`: response `Content-Type` header (defaults to `application/json`)
- `sleep`: number of seconds to wait before responding (defaults to `0`)

Example:

```bash
curl -i 'http://localhost:8080/test?response-code=418&response-content-type=text/plain&response-body=hello&sleep=1'
```

## Request logging format

Each request is printed as JSON to stdout, including:

- `datetime`, `timestamp`
- `from` (client address)
- `method`, `path`, `url`, `query`
- `headers`
- `body`

Body parsing depends on `Content-Type`:

- `application/json`: parsed JSON value
- `application/x-www-form-urlencoded` or `multipart/form-data`: parsed form fields
- anything else: raw body as string
