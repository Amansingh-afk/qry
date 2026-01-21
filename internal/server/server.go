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
	SQL       string `json:"sql"`
	Backend   string `json:"backend"`
	Model     string `json:"model,omitempty"`
	Dialect   string `json:"dialect,omitempty"`
	Warning   string `json:"warning,omitempty"`
	SessionID string `json:"session_id,omitempty"` // For multi-turn conversations
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

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
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

	opts := backend.Options{
		Model:     model,
		Dialect:   dialect,
		SessionID: req.SessionID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	sqlPrompt := prompt.BuildSQL(req.Query, dialect)
	result, err := b.Query(ctx, sqlPrompt, workDir, opts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	sql := prompt.ExtractSQL(result.Response)
	warning := guardrails.Check(sql)

	_ = json.NewEncoder(w).Encode(QueryResponse{
		SQL:       sql,
		Backend:   b.Name(),
		Model:     model,
		Dialect:   dialect,
		Warning:   warning,
		SessionID: result.SessionID,
	})
}
