# Contributing

## Development Setup

```bash
git clone https://github.com/amansingh-afk/qry.git
cd qry
go mod download
go build -o qry .
./qry --help
```

## Project Structure

```
qry/
├── cmd/
│   ├── root.go      # Main command, flags, config
│   ├── init.go      # qry init
│   ├── query.go     # qry "query"
│   └── serve.go     # qry serve
├── internal/
│   ├── backend/     # LLM CLI integrations
│   │   ├── backend.go   # Interface + registry
│   │   ├── claude.go
│   │   ├── gemini.go
│   │   ├── codex.go
│   │   └── cursor.go
│   ├── guardrails/  # SQL safety checks
│   ├── output/      # JSON + pretty output
│   ├── prompt/      # Prompt building
│   ├── server/      # HTTP server
│   └── ui/          # Terminal colors/messages
├── docs/
├── scripts/
└── main.go
```

## Adding a New Backend

1. Create `internal/backend/mybackend.go`:

```go
package backend

import (
    "context"
    "fmt"
    "os/exec"
    "strings"
)

type MyBackend struct{}

func (m *MyBackend) Name() string {
    return "mybackend"
}

func (m *MyBackend) InstallCmd() string {
    return "npm i -g @example/mybackend-cli"
}

func (m *MyBackend) Available() bool {
    _, err := exec.LookPath("mybackend")
    return err == nil
}

func (m *MyBackend) Query(ctx context.Context, prompt string, workDir string) (string, error) {
    cmd := exec.CommandContext(ctx, "mybackend", "-p", prompt)
    cmd.Dir = workDir

    out, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("mybackend: %w\n%s", err, string(out))
    }
    return strings.TrimSpace(string(out)), nil
}
```

2. Register in `internal/backend/backend.go`:

```go
var registry = map[string]Backend{
    "claude":    &Claude{},
    "gemini":    &Gemini{},
    "codex":     &Codex{},
    "cursor":    &Cursor{},
    "mybackend": &MyBackend{},  // Add here
}

func List() []string {
    return []string{"claude", "gemini", "codex", "cursor", "mybackend"}
}
```

## Linting

```bash
golangci-lint run
```

Fix all issues before submitting.

## Testing

```bash
go test ./...
```

## Releasing

Push a tag to trigger the release workflow:

```bash
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions will:
- Build binaries for all platforms
- Create a GitHub Release
- Upload binaries

## Code Style

- Keep functions small and focused
- Minimal comments (code should be self-explanatory)
- Handle errors explicitly
- Use `internal/ui` for all terminal output
