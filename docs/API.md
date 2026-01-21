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
| session_id | string | no | Session ID for multi-turn conversations |

**Response**

```json
{
  "sql": "SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '7 days'",
  "backend": "claude",
  "model": "sonnet",
  "dialect": "postgresql",
  "warning": "",
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
| session_id | string | Session ID for follow-up queries (claude backend only) |

**Error Response**

```json
{
  "error": "claude not installed"
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

The API supports multi-turn conversations where the LLM maintains context across queries. This is useful for refining queries or asking follow-up questions.

**How it works:**
1. Send your first query (no session_id)
2. Response includes a `session_id`
3. Include that `session_id` in follow-up requests
4. The LLM remembers previous context

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

**Follow-up query (with session_id):**
```bash
curl -s -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "filter by active ones", "session_id": "abc123-def456"}'
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
- Session support is currently available for the `claude` backend only
- Other backends will return an empty `session_id`
- Sessions persist codebase context, making follow-up queries faster and more accurate
- To start a fresh conversation, omit `session_id`

**Example integration (Python):**
```python
import requests

session_id = None

def ask_qry(question):
    global session_id
    resp = requests.post("http://localhost:7133/query", json={
        "query": question,
        "session_id": session_id
    })
    data = resp.json()
    session_id = data.get("session_id")
    return data["sql"]

# First query
print(ask_qry("get all users"))
# SELECT * FROM users;

# Follow-up (LLM remembers context)
print(ask_qry("now filter by active"))
# SELECT * FROM users WHERE active = true;
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
