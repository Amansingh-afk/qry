package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

var (
	sqlStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4"))
)

type Result struct {
	SQL     string `json:"sql"`
	Backend string `json:"backend"`
	Model   string `json:"model,omitempty"`
	Dialect string `json:"dialect,omitempty"`
}

func JSON(w io.Writer, sql, backend, model, dialect string) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(Result{
		SQL:     sql,
		Backend: backend,
		Model:   model,
		Dialect: dialect,
	})
}

func Pretty(w io.Writer, sql, backend, model string) {
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, sqlStyle.Render(sql))
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, dimStyle.Render("— "+backend+"/"+model))
}
