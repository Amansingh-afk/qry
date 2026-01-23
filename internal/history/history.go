package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	maxEntries  = 100
	historyFile = "history.json"
)

// Entry represents a single history item
type Entry struct {
	Query     string    `json:"query"`
	SQL       string    `json:"sql"`
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration"` // seconds
}

// path returns the history file path for a working directory
func path(workDir string) string {
	return filepath.Join(workDir, ".qry", historyFile)
}

// Load reads history from disk
func Load(workDir string) ([]Entry, error) {
	data, err := os.ReadFile(path(workDir))
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		// Corrupted file, start fresh
		return []Entry{}, nil
	}

	return entries, nil
}

// Save writes history to disk
func Save(workDir string, entries []Entry) error {
	// Ensure .qry directory exists
	dir := filepath.Join(workDir, ".qry")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Trim to max entries (keep newest)
	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path(workDir), data, 0644)
}

// Add appends a new entry and saves
func Add(workDir string, query, sql string, duration time.Duration) error {
	entries, err := Load(workDir)
	if err != nil {
		entries = []Entry{}
	}

	entries = append(entries, Entry{
		Query:     query,
		SQL:       sql,
		Timestamp: time.Now(),
		Duration:  duration.Seconds(),
	})

	return Save(workDir, entries)
}

// Clear removes all history
func Clear(workDir string) error {
	p := path(workDir)
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
