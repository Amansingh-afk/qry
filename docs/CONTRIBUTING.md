# Contributing

## Development Setup

### Requirements

- Go 1.22+
- One LLM backend (Claude, Gemini, Codex, or Cursor)

### Clone and Build

```bash
git clone https://github.com/amansingh-afk/qry
cd qry
go mod download
go build -o qry .
```

### Using Nix (Optional)

```bash
nix develop
go build
```

## Project Structure

```
qry/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # CLI setup, flags
│   ├── query.go           # Main query command
│   ├── serve.go           # API server
│   └── config.go          # Config management
├── internal/
│   ├── backend/           # LLM backend wrappers
│   ├── prompt/            # SQL prompt templates
│   ├── guardrails/        # Safety checks
│   ├── output/            # Formatting
│   ├── ui/                # Terminal UI
│   └── server/            # HTTP API
└── docs/                  # Documentation
```

## Adding a New Backend

1. Create `internal/backend/newbackend.go`:

```go
package backend

type NewBackend struct{}

func (n *NewBackend) Name() string { return "newbackend" }

func (n *NewBackend) Available() bool {
    _, err := exec.LookPath("newbackend-cli")
    return err == nil
}

func (n *NewBackend) Query(ctx context.Context, prompt string, workDir string) (string, error) {
    cmd := exec.CommandContext(ctx, "newbackend-cli", "-p", prompt)
    cmd.Dir = workDir
    out, err := cmd.CombinedOutput()
    return string(out), err
}
```

2. Register in `internal/backend/backend.go`:

```go
var registry = map[string]Backend{
    // ...
    "newbackend": &NewBackend{},
}
```

## Code Style

- No unnecessary comments
- Keep functions small
- Use descriptive names
- Run `gofmt` before committing

## Testing

```bash
go test ./...
```

## Pull Requests

1. Fork the repo
2. Create a branch: `git checkout -b feature/my-feature`
3. Make changes
4. Test: `go build && ./qry "test query"`
5. Commit: `git commit -m "Add feature"`
6. Push: `git push origin feature/my-feature`
7. Open PR

## Reporting Issues

Include:

- QRY version (`qry version`)
- Backend used
- Command that failed
- Error message
- Working directory context