<p align="center">
  <br>
  <b style="font-size: 48px;">QRY</b>
  <br>
  <b>Ask. Get SQL.</b>
  <br>
  <br>
</p>

<p align="center">
  <a href="#install">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#backends">Backends</a> •
  <a href="docs/API.md">API</a> •
  <a href="docs/CONTRIBUTING.md">Contributing</a>
</p>

---

Generate SQL from natural language. QRY wraps AI CLI tools and runs them in your project directory for full codebase context.

## Install

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/amansingh-afk/qry/main/scripts/install.sh | sh
```

### Homebrew

```bash
brew install amansingh/tap/qry
```

### Go

```bash
go install github.com/amansingh-afk/qry@latest
```

## Usage

```bash
cd your-project
qry "get active users from last week"
```

```sql
SELECT * FROM users
WHERE active = true
  AND created_at >= CURRENT_DATE - INTERVAL '7 days'
```

### Initialize Project

```bash
cd your-project
qry init
```

Creates `.qry.yaml` with default settings.

### Choose Backend

```bash
qry -b claude "get active users"
qry -b gemini "get active users"
qry -b codex "get active users"
qry -b cursor "get active users"
```

### JSON Output

```bash
qry --json "get active users"
```

```json
{
  "sql": "SELECT * FROM users WHERE active = true",
  "backend": "claude",
  "elapsed_ms": 342,
  "safe": true
}
```

### API Server

Run from your project directory:

```bash
cd your-project
qry serve
```

```bash
curl -X POST http://localhost:7133/query \
  -H "Content-Type: application/json" \
  -d '{"prompt": "get active users"}'
```

### Configuration

Config is per-project in `.qry.yaml`:

```bash
qry config set backend gemini
qry config set dialect postgres
qry config list
```

## Backends

QRY wraps existing AI CLI tools. Install at least one:

| Backend | Install |
|---------|---------|
| Claude | `npm i -g @anthropic-ai/claude-code` |
| Gemini | `npm i -g @google/gemini-cli` |
| Codex | `npm i -g @openai/codex` |
| Cursor | `curl https://cursor.com/install -fsS \| bash` |

## How It Works

1. Run `qry` from your project directory
2. QRY calls the AI CLI (Claude/Gemini/etc.) from that directory
3. The AI CLI has full context of your codebase
4. Returns SQL specific to your schema

## Documentation

- [Setup Guide](docs/SETUP.md)
- [API Reference](docs/API.md)
- [Contributing](docs/CONTRIBUTING.md)

## License

[MIT](LICENSE)
