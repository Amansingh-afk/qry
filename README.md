# QRY

**Ask. Get SQL.**

Natural language to SQL using your existing LLM CLI.

```bash
qry q "get users who signed up last week"
```

```sql
SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '7 days';
```

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
qry q "find inactive users"
```

## Commands

| Command | Description |
|---------|-------------|
| `qry q "query"` | One-shot SQL generation |
| `qry chat` | Interactive session |
| `qry init` | Setup config |
| `qry serve` | Start API server |

## One-shot Mode

```bash
qry q "count orders by month"
qry q "get top 10 products" --json
qry q "find users" -m sonnet -d postgresql
qry q "complex query" --dry-run
```

## Interactive Mode

Start a chat session with context persistence:

```bash
qry chat                    # New session
qry chat --continue         # Resume last session
qry chat -r <session-id>    # Resume specific session
```

Follow-up questions work naturally:

```
> get active users
SELECT * FROM users WHERE status = 'active';

> filter by last 7 days
SELECT * FROM users WHERE status = 'active' AND created_at >= NOW() - INTERVAL '7 days';
```

## Config

Create `.qry.yaml` in your project:

```yaml
backend: claude
model: sonnet
dialect: postgresql
timeout: 30s
```

## API Server

```bash
qry serve
```

```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"query": "get active users"}'
```

## Supported Backends

| Backend | Default Model | Install |
|---------|---------------|---------|
| Claude | haiku | `npm i -g @anthropic-ai/claude-code` |
| Gemini | gemini-2.0-flash | `npm i -g @google/gemini-cli` |
| Codex | gpt-4o-mini | `npm i -g @openai/codex` |
| Cursor | gpt-4o-mini | `curl -fsSL https://cursor.com/install \| sh` |

Use `-m` to override: `qry q "query" -m sonnet`

## Docs

- [Setup Guide](docs/SETUP.md)
- [API Reference](docs/API.md)
- [Contributing](docs/CONTRIBUTING.md)

## License

MIT
