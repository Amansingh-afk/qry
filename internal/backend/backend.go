package backend

import (
	"context"
	"fmt"
)

type Backend interface {
	Name() string
	Available() bool
	Query(ctx context.Context, prompt string, workDir string) (string, error)
}

var registry = map[string]Backend{
	"cursor": &Cursor{},
	"claude": &Claude{},
	"gemini": &Gemini{},
	"codex":  &Codex{},
}

func Get(name string) (Backend, error) {
	if b, ok := registry[name]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("unknown backend: %s (available: cursor, claude, gemini, codex)", name)
}

func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
