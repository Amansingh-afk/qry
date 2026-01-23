package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amansingh-afk/qry/internal/backend"
	"github.com/amansingh-afk/qry/internal/guardrails"
	"github.com/amansingh-afk/qry/internal/prompt"
	"github.com/amansingh-afk/qry/internal/security"
	"github.com/amansingh-afk/qry/internal/session"
	"github.com/spf13/viper"
)

type QueryRequest struct {
	Query     string `json:"query"`
	Backend   string `json:"backend,omitempty"`
	Model     string `json:"model,omitempty"`
	Dialect   string `json:"dialect,omitempty"`
	SessionID string `json:"session_id,omitempty"` // For multi-turn conversations
}

type QueryResponse struct {
	SQL             string `json:"sql"`
	Backend         string `json:"backend"`
	Model           string `json:"model,omitempty"`
	Dialect         string `json:"dialect,omitempty"`
	Warning         string `json:"warning,omitempty"`
	SecurityWarning string `json:"security_warning,omitempty"`
	SessionID       string `json:"session_id,omitempty"` // For multi-turn conversations
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Start(port int, workDir string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /query", func(w http.ResponseWriter, r *http.Request) {
		handleQuery(w, r, workDir)
	})

	mux.HandleFunc("GET /session", func(w http.ResponseWriter, r *http.Request) {
		handleGetSession(w, r, workDir)
	})

	mux.HandleFunc("DELETE /session", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteSession(w, r, workDir)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// getSessionTTL parses the session TTL from config
func getSessionTTL() time.Duration {
	ttlStr := viper.GetString("session.ttl")
	if ttlStr == "" {
		return 7 * 24 * time.Hour // Default 7 days
	}

	// Handle "Xd" format (days)
	if len(ttlStr) > 1 && ttlStr[len(ttlStr)-1] == 'd' {
		var days int
		if _, err := fmt.Sscanf(ttlStr, "%dd", &days); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}

	// Try standard duration format
	if d, err := time.ParseDuration(ttlStr); err == nil {
		return d
	}

	return 7 * 24 * time.Hour // Fallback
}

func handleQuery(w http.ResponseWriter, r *http.Request, workDir string) {
	w.Header().Set("Content-Type", "application/json")

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
		return
	}

	if req.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "query required"})
		return
	}

	backendName := req.Backend
	if backendName == "" {
		backendName = viper.GetString("backend")
	}

	b, err := backend.Get(backendName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if !b.Available() {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("%s not available", b.Name())})
		return
	}

	// Get model: request > config model > config defaults
	model := req.Model
	if model == "" {
		model = viper.GetString("model")
	}
	if model == "" {
		model = viper.GetString("defaults." + backendName)
	}

	dialect := req.Dialect
	if dialect == "" {
		dialect = viper.GetString("dialect")
	}

	// Server-side session management: use stored session if client didn't provide one
	sessionID := req.SessionID
	if sessionID == "" {
		if s, _ := session.GetOrCreate(workDir, backendName, getSessionTTL()); s != nil {
			sessionID = s.SessionID
		}
	}

	opts := backend.Options{
		Model:     model,
		Dialect:   dialect,
		SessionID: sessionID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	// Use full prompt for new sessions, minimal prompt for existing sessions
	var sqlPrompt string
	if sessionID == "" {
		sqlPrompt = prompt.BuildSQL(req.Query, dialect)
	} else {
		sqlPrompt = prompt.BuildFollowUp(req.Query)
	}

	result, err := b.Query(ctx, sqlPrompt, workDir, opts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Persist session for future requests
	if result.SessionID != "" {
		_ = session.Update(workDir, backendName, result.SessionID)
	}

	sql := prompt.ExtractSQL(result.Response)

	// Security validation
	secResult := security.Validate(sql)
	sec := security.Get()

	if sec.IsBlocked(secResult) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Security violation: " + secResult.Summary(),
		})
		return
	}

	var securityWarning string
	if sec.ShouldWarn(secResult) {
		securityWarning = secResult.Error()
	}

	warning := guardrails.Check(sql)

	_ = json.NewEncoder(w).Encode(QueryResponse{
		SQL:             sql,
		Backend:         b.Name(),
		Model:           model,
		Dialect:         dialect,
		Warning:         warning,
		SecurityWarning: securityWarning,
		SessionID:       result.SessionID,
	})
}

// SessionResponse represents session info
type SessionResponse struct {
	Backend   string `json:"backend"`
	SessionID string `json:"session_id"`
	CreatedAt string `json:"created_at"`
	Age       string `json:"age"`
}

func handleGetSession(w http.ResponseWriter, r *http.Request, workDir string) {
	w.Header().Set("Content-Type", "application/json")

	s, err := session.Load(workDir)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "no session found"})
		return
	}

	age := time.Since(s.CreatedAt).Round(time.Minute)

	_ = json.NewEncoder(w).Encode(SessionResponse{
		Backend:   s.Backend,
		SessionID: s.SessionID,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		Age:       age.String(),
	})
}

func handleDeleteSession(w http.ResponseWriter, r *http.Request, workDir string) {
	w.Header().Set("Content-Type", "application/json")

	if err := session.Delete(workDir); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
