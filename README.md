# QRY

> Natural language to SQL that actually knows your schema — no sync required.

**Ask. Get SQL.**

![QRY Demo](docs/qry.gif)

## How it works

QRY wraps CLIs (Claude Code, Codex, Cursor) that already understand your codebase. No custom indexing, no embeddings, no schema sync — it leverages their built-in context awareness. It is basically Repo2SQL.

**Why this matters:**
- New table? Just `git pull`. The CLI already sees it.
- Schema change? No manual updates. It's already indexed.
- Complex joins? The backend knows your actual table relationships.

Traditional NL2SQL tools require you to maintain schema definitions and regenerate embeddings. QRY doesn't — the underlying CLI handles all of that.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/amansingh-afk/qry/main/scripts/install.sh | bash
```

Or with Go:

```bash
go install github.com/amansingh-afk/qry@latest
```

## Quick Start

```bash
cd your-project
qry init
qry
```

## Commands

| Command | Description |
|---------|-------------|
| `qry` | Interactive chat (default) |
| `qry q "query"` | One-shot query (for scripting) |
| `qry init` | Setup config |
| `qry init --force` | Reset session (re-index codebase) |
| `qry serve` | Start API server |

## Interactive Mode

Just run `qry`:

```bash
qry
```

Follow-up questions work naturally:

```
> get active users
SELECT * FROM users WHERE status = 'active';

> filter by last 7 days
SELECT * FROM users WHERE status = 'active' AND created_at >= NOW() - INTERVAL '7 days';
```

## One-shot Mode

For scripting and piping:

```bash
qry q "count orders by month"
qry q "get top 10 products" --json
qry q "find users" -m sonnet -d postgresql
qry q "get users" | pbcopy
```

Or if you're feeling brave:

```bash
qry q "get users" | psql
```

## Session Management

QRY maintains a unified session. The LLM indexes your codebase once and remembers context for subsequent queries.

- Sessions persist for 7 days by default (configurable)
- Same session is shared between one-shot and chat modes
- Sessions auto-invalidate if you switch backends
- Full prompt (role + rules) sent only on first query; follow-ups send just the query

**Reset session to re-index codebase:**
```bash
qry init --force
```

## Config

`qry init` creates `.qry.yaml` in your project:

```yaml
backend: claude
dialect: postgresql
db_version: "16"             # Optional: PostgreSQL 16, MySQL 8.0, etc.
timeout: 2m

session:
  ttl: 7d                    # Session lifetime

prompt: |                    # Customizable prompt template
  You are a SQL expert. Based on the codebase context (schemas, migrations, models), generate ONLY the SQL query.
  
  Rules:
  - Output ONLY the SQL, no explanation
  - Use actual table/column names from the codebase
  - Use {{dialect}}{{version}} syntax
  
  Request: {{query}}

defaults:
  claude: haiku
  codex: gpt-4o-mini
  cursor: auto
```

| Field | Description |
|-------|-------------|
| `backend` | LLM CLI to use (claude, codex, cursor) |
| `dialect` | SQL syntax (postgresql, mysql, sqlite) |
| `db_version` | Database version for accurate syntax (e.g., `16`, `8.0`) |
| `timeout` | Request timeout |
| `session.ttl` | Session lifetime (e.g., `7d`, `24h`) |
| `prompt` | Prompt template with `{{dialect}}`, `{{version}}`, `{{query}}` variables |

## API Server

Build Slack bots, admin tools, or n8n workflows on top of QRY.

Run from your project directory:

```bash
cd your-project
qry serve
```

```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "get active users"}'
```

The server is scoped to the directory it's started in. For multiple repos, run separate servers on different ports (`qry serve -p 7134`).

Sessions are managed automatically — no need to track session IDs.

**Session endpoints:**
```bash
# View current session
curl http://localhost:7133/session

# Reset session (re-index)
curl -X DELETE http://localhost:7133/session
```

## Supported Backends

| Backend | Default Model | Install |
|---------|---------------|---------|
| Claude | haiku | `npm i -g @anthropic-ai/claude-code` |
| Codex | gpt-4o-mini | `npm i -g @openai/codex` |
| Cursor | auto | `curl -fsSL https://cursor.com/install \| sh` |
| Gemini | — | *WIP* |

Use `-m` to override: `qry q "query" -m sonnet`

## Docs

- [Setup Guide](docs/SETUP.md)
- [API Reference](docs/API.md)
- [Contributing](docs/CONTRIBUTING.md)

## License

MIT
