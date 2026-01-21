# Setup Guide

## 1. Install a Backend

QRY needs at least one LLM CLI installed:

```bash
# Claude (recommended)
npm i -g @anthropic-ai/claude-code

# Gemini
npm i -g @google/gemini-cli

# Codex
npm i -g @openai/codex

# Cursor
curl -fsSL https://cursor.com/install | sh
```

## 2. Install QRY

```bash
# Quick install
curl -fsSL https://raw.githubusercontent.com/amansingh-afk/qry/main/scripts/install.sh | bash

# With Go
go install github.com/amansingh-afk/qry@latest

# From releases
# https://github.com/amansingh-afk/qry/releases
```

## 3. Initialize

```bash
cd your-project
qry init
```

This creates:
- `.qry.yaml` — configuration file
- `.qry/` — session storage (gitignored automatically)

**Reset and re-index codebase:**
```bash
qry init --force
```

## 4. Configure

Edit `.qry.yaml`:

```yaml
backend: claude
dialect: postgresql
db_version: "16"             # Optional: PostgreSQL 16, MySQL 8.0, etc.
timeout: 2m

session:
  ttl: 7d                    # Session lifetime

prompt: |                    # Customizable prompt
  You are a SQL expert. Based on the codebase context (schemas, migrations, models), generate ONLY the SQL query.
  
  Rules:
  - Output ONLY the SQL, no explanation
  - Use actual table/column names from the codebase
  - Use {{dialect}}{{version}} syntax
  
  Request: {{query}}

defaults:
  claude: haiku
  gemini: gemini-2.0-flash
  codex: gpt-4o-mini
  cursor: gpt-4o-mini
```

| Field | Values | Description |
|-------|--------|-------------|
| backend | claude, gemini, codex, cursor | LLM CLI to use |
| model | (depends on backend) | Model to use |
| dialect | postgresql, mysql, sqlite | SQL syntax |
| db_version | 16, 8.0, 3, etc. | Database version for accurate syntax |
| timeout | 30s, 1m, 2m | Request timeout |
| session.ttl | 7d, 24h, 168h | Session lifetime before re-index |
| prompt | template string | Prompt with `{{dialect}}`, `{{version}}`, `{{query}}` |

## 5. Usage

**Interactive mode (default):**
```bash
qry
```

**One-shot queries (for scripting):**
```bash
qry q "get active users"
qry q "count by month" --json
qry q "find users" -m opus
qry q "get users" | pbcopy
```

Both modes share the same session — the LLM indexes your codebase once and remembers context.

## Session Management

QRY maintains a unified session that:
- Persists across `qry` and `qry q` commands
- Auto-expires after the configured TTL (default: 7 days)
- Auto-invalidates if you switch backends
- Sends full prompt only on first query; follow-ups are minimal (just the query)

**View session info (API):**
```bash
curl http://localhost:7133/session
```

**Reset session:**
```bash
qry init --force
# or via API
curl -X DELETE http://localhost:7133/session
```

## API Server

```bash
qry serve              # Port 7133
qry serve -p 8080      # Custom port
```

The server manages sessions automatically — clients don't need to track session IDs.

## Troubleshooting

**"X not installed"**
```bash
npm i -g @anthropic-ai/claude-code
```

**"No backends found"**
```bash
qry init
```

**Server not finding schema**
```bash
cd your-project
qry serve
```

**Stale context / schema changed**
```bash
qry init --force
```
