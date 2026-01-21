# Setup Guide

## Prerequisites

QRY requires at least one LLM CLI backend installed.

### Install a Backend

| Backend | Install Command |
|---------|-----------------|
| Claude | `npm i -g @anthropic-ai/claude-code` |
| Gemini | `npm i -g @google/gemini-cli` |
| Codex | `npm i -g @openai/codex` |
| Cursor | `curl https://cursor.com/install -fsS \| bash` |

## Install QRY

### Option 1: Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/amansingh-afk/qry/main/scripts/install.sh | sh
```

### Option 2: Homebrew

```bash
brew install amansingh/tap/qry
```

### Option 3: Go Install

```bash
go install github.com/amansingh-afk/qry@latest
```

### Option 4: Build from Source

```bash
git clone https://github.com/amansingh-afk/qry
cd qry
go build -o qry .
sudo mv qry /usr/local/bin/
```

## Initialize Project

```bash
cd your-project
qry init
```

This creates `.qry.yaml`:

```yaml
backend: claude
dialect: postgres
```

## Configuration

QRY uses a per-project `.qry.yaml` file in your project directory.

### Set Config

```bash
qry config set backend gemini
qry config set dialect mysql
```

### View Config

```bash
qry config list
qry config get backend
```

### Environment Variables

You can also use environment variables:

```bash
export QRY_BACKEND=gemini
export QRY_DIALECT=mysql
```

## Verify Installation

```bash
qry version
qry "get active users"
```

## Backend Authentication

Each backend requires its own authentication:

### Claude

```bash
claude auth
```

### Gemini

```bash
gemini auth
```

### Codex

```bash
export OPENAI_API_KEY=your-key
```

### Cursor

Uses your existing Cursor IDE authentication.

## Running as API Server

Run from your project directory:

```bash
cd your-project
qry serve
```

Default port is 7133. Change with `--port`:

```bash
qry serve --port 9000
```

The server will have context of the directory it was started from.

## Troubleshooting

### "backend not found"

Make sure the CLI is installed and in your PATH:

```bash
which claude  # or gemini, codex, cursor
```

### "command not found: qry"

Add Go bin to PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Wrong schema context

Make sure you're running QRY from your project directory:

```bash
cd /path/to/your/project
qry "get users"
```
