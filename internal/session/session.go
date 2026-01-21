package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	sessionDir  = ".qry"
	sessionFile = "session"
)

// Session holds the persistent session state
type Session struct {
	Backend   string    `json:"backend"`
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Path returns the session file path for the given work directory
func Path(workDir string) string {
	return filepath.Join(workDir, sessionDir, sessionFile)
}

// DirPath returns the .qry directory path
func DirPath(workDir string) string {
	return filepath.Join(workDir, sessionDir)
}

// Load reads the session from disk
func Load(workDir string) (*Session, error) {
	data, err := os.ReadFile(Path(workDir))
	if err != nil {
		return nil, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

// Save writes the session to disk
func Save(workDir string, s *Session) error {
	dir := DirPath(workDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(Path(workDir), data, 0644)
}

// Delete removes the session file
func Delete(workDir string) error {
	err := os.Remove(Path(workDir))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// IsValid checks if the session is still valid
// Returns false if:
// - Session is older than TTL
// - Backend doesn't match current backend
func (s *Session) IsValid(currentBackend string, ttl time.Duration) bool {
	if s == nil || s.SessionID == "" {
		return false
	}

	// Check backend match
	if s.Backend != currentBackend {
		return false
	}

	// Check TTL (0 means no expiry)
	if ttl > 0 && time.Since(s.CreatedAt) > ttl {
		return false
	}

	return true
}

// GetOrCreate loads an existing valid session or returns nil if a new one should be created
// If session is invalid (expired or backend changed), it deletes the old session
func GetOrCreate(workDir, backend string, ttl time.Duration) (*Session, error) {
	s, err := Load(workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No session exists, caller should create new
		}
		return nil, err
	}

	if !s.IsValid(backend, ttl) {
		// Session invalid, delete it
		_ = Delete(workDir)
		return nil, nil
	}

	return s, nil
}

// Update saves a new or updated session
func Update(workDir, backend, sessionID string) error {
	// Load existing to preserve created_at if same session
	existing, _ := Load(workDir)

	s := &Session{
		Backend:   backend,
		SessionID: sessionID,
		CreatedAt: time.Now(),
	}

	// If same session ID, preserve original creation time
	if existing != nil && existing.SessionID == sessionID {
		s.CreatedAt = existing.CreatedAt
	}

	return Save(workDir, s)
}
