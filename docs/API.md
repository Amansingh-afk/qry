# API Reference

## Starting the Server

```bash
cd your-project
qry serve
```

Default port: `7133`

Custom port:
```bash
qry serve -p 8080
```

## Endpoints

### POST /query

Generate SQL from natural language.

**Request**

```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "get users who signed up last week"}'
```

**Request Body**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| query | string | yes | Natural language query |
| backend | string | no | Override default backend |
| model | string | no | Model to use |
| dialect | string | no | SQL dialect (postgresql, mysql, sqlite) |
| session_id | string | no | Override server-managed session |

**Response**

```json
{
  "sql": "SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '7 days'",
  "backend": "claude",
  "model": "sonnet",
  "dialect": "postgresql",
  "warning": "",
  "security_warning": "",
  "session_id": "abc123-def456"
}
```

| Field | Type | Description |
|-------|------|-------------|
| sql | string | Generated SQL |
| backend | string | Backend used |
| model | string | Model used |
| dialect | string | SQL dialect |
| warning | string | Safety warning (if any) |
| security_warning | string | Security warning (if in warn mode) |
| session_id | string | Session ID (managed by server) |

**Error Response**

```json
{
  "error": "claude not installed"
}
```

**Security Violation (403)**

If security mode is `strict` and the query references excluded data:

```json
{
  "error": "Security violation: blocked: references table 'api_keys'"
}
```

### GET /session

Get current session info.

**Request**

```bash
curl http://localhost:7133/session
```

**Response**

```json
{
  "backend": "claude",
  "session_id": "abc123-def456",
  "created_at": "2026-01-15T10:30:00Z",
  "age": "6d2h30m"
}
```

**Error Response (no session)**

```json
{
  "error": "no session found"
}
```

### DELETE /session

Delete current session (forces re-index on next query).

**Request**

```bash
curl -X DELETE http://localhost:7133/session
```

**Response**

```json
{
  "status": "deleted"
}
```

### GET /health

Health check endpoint.

**Request**

```bash
curl http://localhost:7133/health
```

**Response**

```json
{
  "status": "ok"
}
```

## Session Management

The server manages sessions automatically. You don't need to track session IDs.

**How it works:**
1. First query → Full prompt sent (role + rules), LLM indexes codebase, server stores session
2. Subsequent queries → Only the query is sent, LLM already has context
3. Session expires after TTL (default: 7 days) or when backend changes

**Session lifecycle:**
- Sessions persist across server restarts (stored in `.qry/session`)
- Sessions auto-invalidate if you switch backends
- Sessions auto-expire after configured TTL
- Full prompt template is only sent on first query (token efficient)

**Force re-index:**
```bash
# Via API
curl -X DELETE http://localhost:7133/session

# Via CLI
qry init --force
```

## Examples

**Basic query:**
```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "count users"}'
```

**With model and dialect:**
```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "count users", "model": "opus", "dialect": "postgresql"}'
```

**Using jq for pretty output:**
```bash
curl -s -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "get active users"}' | jq .sql -r
```

## Multi-Turn Conversations

The server maintains context across queries automatically.

**First query:**
```bash
curl -s -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "get all users"}'
```

Response:
```json
{
  "sql": "SELECT * FROM users;",
  "backend": "claude",
  "session_id": "abc123-def456"
}
```

**Follow-up query (no session_id needed):**
```bash
curl -s -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "filter by active ones"}'
```

Response:
```json
{
  "sql": "SELECT * FROM users WHERE active = true;",
  "backend": "claude",
  "session_id": "abc123-def456"
}
```

**Notes:**
- Server manages sessions — clients don't need to track session IDs
- You can still pass `session_id` to override the server-managed session
- Sessions persist codebase context, making queries faster and more accurate

**Example integration (Python):**
```python
import requests

BASE = "http://localhost:7133"

def ask_qry(question):
    resp = requests.post(f"{BASE}/query", json={"query": question})
    return resp.json()["sql"]

# First query (server creates session, indexes codebase)
print(ask_qry("get all users"))
# SELECT * FROM users;

# Follow-up (server reuses session, LLM has context)
print(ask_qry("now filter by active"))
# SELECT * FROM users WHERE active = true;

# Reset session if needed
requests.delete(f"{BASE}/session")
```

## Safety Warnings

The API returns warnings for potentially destructive operations:

- `DROP TABLE`
- `DROP DATABASE`
- `TRUNCATE`
- `DELETE FROM`
- `ALTER TABLE`
- `UPDATE` without `WHERE`

These are warnings only. The SQL is still returned.

## Security

If security rules are configured in `.qry.yaml`, the API enforces them:

**Warn Mode (default)**

SQL is returned with `security_warning` field:

```json
{
  "sql": "SELECT * FROM api_keys WHERE ...",
  "security_warning": "Security violation: query references excluded data\n  - table: api_keys (matched rule: api_*)\n"
}
```

**Strict Mode**

Returns `403 Forbidden`:

```json
{
  "error": "Security violation: blocked: references table 'api_keys'"
}
```

See the [Setup Guide](./SETUP.md) for configuration details.
