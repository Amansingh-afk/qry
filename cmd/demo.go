package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/amansingh-afk/qry/internal/output"
	"github.com/amansingh-afk/qry/internal/ui"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:    "demo",
	Short:  "Run a demo showcasing QRY",
	Hidden: true, // Hidden from help, used for recording
	Run:    runDemo,
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

type demoQuery struct {
	question string
	sql      string
	duration time.Duration
}

func runDemo(cmd *cobra.Command, args []string) {
	// Demo queries - from simple to complex
	queries := []demoQuery{
		{
			question: "get active users",
			sql:      "SELECT * FROM users WHERE status = 'active';",
			duration: 1200 * time.Millisecond,
		},
		{
			question: "only from last 7 days, sort by newest",
			sql: `SELECT * FROM users
WHERE status = 'active'
  AND created_at >= NOW() - INTERVAL '7 days'
ORDER BY created_at DESC;`,
			duration: 900 * time.Millisecond,
		},
		{
			question: "top 10 customers by revenue with order count",
			sql: `SELECT
  c.id,
  c.name,
  c.email,
  SUM(o.total) as total_revenue,
  COUNT(o.id) as order_count
FROM customers c
JOIN orders o ON o.customer_id = c.id
WHERE o.status = 'completed'
GROUP BY c.id, c.name, c.email
ORDER BY total_revenue DESC
LIMIT 10;`,
			duration: 2100 * time.Millisecond,
		},
		{
			question: "monthly sales trend for 2024 with growth rate",
			sql: `WITH monthly_sales AS (
  SELECT
    DATE_TRUNC('month', created_at) as month,
    SUM(total) as revenue
  FROM orders
  WHERE created_at >= '2024-01-01'
    AND status = 'completed'
  GROUP BY DATE_TRUNC('month', created_at)
)
SELECT
  month,
  revenue,
  LAG(revenue) OVER (ORDER BY month) as prev_month,
  ROUND(
    (revenue - LAG(revenue) OVER (ORDER BY month)) * 100.0 /
    NULLIF(LAG(revenue) OVER (ORDER BY month), 0), 2
  ) as growth_pct
FROM monthly_sales
ORDER BY month;`,
			duration: 2800 * time.Millisecond,
		},
	}

	// Show header
	ui.ChatStarting("claude", "sonnet", "~/myapp")

	// Run through demo queries
	for _, q := range queries {
		// Show prompt and question
		fmt.Print(ui.Prompt())
		typeText(q.question, 35*time.Millisecond)
		fmt.Println()

		// Simulate thinking with spinner
		simulateThinking(q.duration)

		// Show duration
		ui.QueryDone(q.duration)

		// Show SQL result
		output.Pretty(os.Stdout, q.sql, "claude", "sonnet")
		fmt.Println()

		// Pause between queries
		time.Sleep(1500 * time.Millisecond)
	}

	// Final prompt
	fmt.Print(ui.Prompt())
	time.Sleep(2 * time.Second)
}

// typeText simulates typing with a delay between characters
func typeText(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(delay)
	}
}

// simulateThinking shows a spinner for the given duration
func simulateThinking(duration time.Duration) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	start := time.Now()
	i := 0

	for time.Since(start) < duration {
		elapsed := time.Since(start).Seconds()
		fmt.Printf("\r  \033[90m%s thinking... %.1fs\033[0m", frames[i%len(frames)], elapsed)
		i++
		time.Sleep(80 * time.Millisecond)
	}
	fmt.Print("\r\033[K") // Clear line
}
