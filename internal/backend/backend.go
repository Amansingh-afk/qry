package backend

import (
	"context"
	"fmt"
)

type Options struct {
	Model     string
	Dialect   string
	SessionID string // For multi-turn conversations (claude only)
}

type Result struct {
	Response  string
	SessionID string // Returned for multi-turn conversations
}

type Backend interface {
	Name() string
	Available() bool
	InstallCmd() string
	Query(ctx context.Context, prompt string, workDir string, opts Options) (Result, error)
}

var registry = map[string]Backend{
	"claude": &Claude{},
	// "gemini": &Gemini{}, // WIP: needs account testing
	"codex":  &Codex{},
	"cursor": &Cursor{},
}

func Get(name string) (Backend, error) {
	b, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown backend: %s", name)
	}
	return b, nil
}

func List() []string {
	return []string{"claude", "codex", "cursor"}
}

func Available() []Backend {
	var available []Backend
	for _, name := range List() {
		if b, _ := Get(name); b != nil && b.Available() {
			available = append(available, b)
		}
	}
	return available
}
