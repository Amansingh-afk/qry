# API Reference

QRY can run as an HTTP API server for integrations.

## Start Server

Run from your project directory:

```bash
cd your-project
qry serve
```

Default port is **7133**. Change with `--port`:

```bash
qry serve --port 9000
```

The server uses the context of the directory it was started from.

## Endpoints

### Health Check

```
GET /health
```

Response:

```json
{
  "status": "ok",
  "workdir": "/path/to/your/project"
}
```

### List Backends

```
GET /backends
```

Response:

```json
{
  "backends": ["cursor", "claude", "gemini", "codex"]
}
```

### Query

```
POST /query
```

Request:

```json
{
  "prompt": "get active users from last week",
  "backend": "claude"
}
```

Response:

```json
{
  "sql": "SELECT * FROM users WHERE active = true AND created_at >= CURRENT_DATE - INTERVAL '7 days'",
  "backend": "claude",
  "elapsed_ms": 1234,
  "safe": true,
  "warning": ""
}
```

Error Response:

```json
{
  "error": "backend not available"
}
```

## Examples

### cURL

```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"prompt": "count users by country"}'
```

### JavaScript

```javascript
const response = await fetch('http://localhost:7133/query', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ prompt: 'get active users' })
});
const data = await response.json();
console.log(data.sql);
```

### Python

```python
import requests

response = requests.post('http://localhost:7133/query', json={
    'prompt': 'get orders over $100'
})
print(response.json()['sql'])
```

## Per-Project Setup

Each project should have its own QRY server running from its directory:

```bash
# Project A
cd /repos/backend-api
qry serve

# Project B (different port)
cd /repos/analytics
qry serve --port 7134
```

Each server will have context specific to its project.
