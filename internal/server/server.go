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
	Prompt  string `json:"prompt"`
	Backend string `json:"backend,omitempty"`
}

type QueryResponse struct {
	SQL     string `json:"sql"`
	Backend string `json:"backend"`
	Elapsed int64  `json:"elapsed_ms"`
	Safe    bool   `json:"safe"`
	Warning string `json:"warning,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var serverWorkDir string

func Start(port int, workDir string) error {
	serverWorkDir = workDir

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /query", handleQuery)
	mux.HandleFunc("GET /backends", handleBackends)

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"workdir": serverWorkDir,
	})
}

func handleBackends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"backends": backend.List()})
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Prompt == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "prompt is required"})
		return
	}

	backendName := req.Backend
	if backendName == "" {
		backendName = viper.GetString("backend")
	}

	b, err := backend.Get(backendName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if !b.Available() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("%s is not available", backendName)})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	start := time.Now()
	sqlPrompt := prompt.BuildSQL(req.Prompt)
	result, err := b.Query(ctx, sqlPrompt, serverWorkDir)
	elapsed := time.Since(start)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	sql := prompt.ExtractSQL(result)
	warning := guardrails.Check(sql)

	resp := QueryResponse{
		SQL:     sql,
		Backend: backendName,
		Elapsed: elapsed.Milliseconds(),
		Safe:    warning == "",
		Warning: warning,
	}

	json.NewEncoder(w).Encode(resp)
}
