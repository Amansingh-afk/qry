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

## 4. Configure

Edit `.qry.yaml`:

```yaml
backend: claude
model: sonnet
dialect: postgresql
timeout: 30s
```

| Field | Values | Description |
|-------|--------|-------------|
| backend | claude, gemini, codex, cursor | LLM CLI to use |
| model | (depends on backend) | Model to use |
| dialect | postgresql, mysql, sqlite | SQL syntax |
| timeout | 30s, 1m, 2m | Request timeout |

## 5. Usage

**One-shot queries:**
```bash
qry q "get active users"
qry q "count by month" --json
qry q "find users" -m opus
```

**Interactive session:**
```bash
qry chat                # New session
qry chat --continue     # Resume last
```

## API Server

```bash
qry serve              # Port 7133
qry serve -p 8080      # Custom port
```

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
