# QRY - Jobs To Be Done & Strategy

## The Name

**QRY**
- Pronunciation: "query" or "Q-R-Y" 
- Meaning: Query (simplified/minimal)
- 3 letters - ultra minimal

---

## What QRY Actually Is

QRY is a **Go CLI wrapper** around existing LLM tools (Claude CLI, Cursor, Gemini CLI) that specializes in SQL generation from natural language.

### The Key Insight

Claude, Cursor, and Gemini **already excel at codebase context understanding**. They can read your migrations, models, and schema files natively. 

**We don't need to build custom parsers.** We just need to:
1. Take user's natural language query
2. Pass it to the LLM CLI with a SQL-focused prompt
3. Return clean, structured SQL output

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                           QRY                                │
│                                                              │
│  • SQL-specialized prompts                                   │
│  • Structured JSON output                                    │
│  • Guardrails (block destructive queries)                   │
│  • Multi-backend support                                     │
│  • Self-hosted API server mode                               │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│            Claude CLI / Cursor / Gemini CLI                  │
│            (They handle codebase context)                    │
└─────────────────────────────────────────────────────────────┘
```

### What QRY Does NOT Do

- ❌ Custom migration parsing (LLM handles this)
- ❌ Build internal schema representation (LLM understands raw files)
- ❌ Direct LLM API calls (wraps existing CLIs)
- ❌ Reinvent codebase indexing (leverage existing tools)

### What QRY DOES Do

- ✅ SQL-optimized prompting
- ✅ Clean, structured output (JSON with metadata)
- ✅ Guardrails (block DROP, DELETE, etc.)
- ✅ Multi-backend switching (claude/gemini/cursor)
- ✅ Self-hosted API server mode
- ✅ Multi-repo, multi-DB support
- ✅ Beautiful CLI UX

---

## Why QRY is GENIUS

### ✅ Massive Strengths

1. **ULTRA SHORT** - Only 3 letters! (like npm, git, aws)
2. **Easy to type** - Super fast CLI command
3. **Minimal aesthetic** - Modern, clean, developer-focused
4. **Memorable** - Impossible to forget
5. **Pronounceable** - Just say "query" naturally
6. **Available** - Very likely free (.dev, .io)
7. **Scales globally** - No translation issues

### The "3-Letter Club"

You're joining an elite club of iconic 3-letter dev tools:

- **npm** - Node Package Manager
- **git** - Version control
- **aws** - Amazon Web Services  
- **gcp** - Google Cloud Platform
- **sql** - Structured Query Language

**QRY fits RIGHT IN** 🔥

---

## Technical Implementation

### Built with Go

Why Go?
- Single binary distribution (no runtime deps)
- Cross-platform builds (Linux/Mac/Windows)
- Fast startup (critical for CLI tools)
- Excellent CLI libraries (Cobra, Viper)
- Developer trust (kubectl, docker, terraform are Go)

### Core Commands

```bash
# Basic usage
qry "get active users from last week"

# Structured JSON output
qry --json "get active users"
# Output: {"sql": "SELECT...", "tables": ["users"], "type": "SELECT", "safe": true}

# Choose backend
qry --backend=claude "get orders over $100"
qry --backend=gemini "get orders over $100"

# Force SQL dialect
qry --dialect=postgres "get active users"
qry --dialect=mysql "get active users"

# Guardrail check (dry-run for destructive queries)
qry --dry-run "delete inactive users"

# API server mode (self-hosted)
qry serve --port 8080

# Configure
qry config set backend claude
qry config set guardrails.block-destructive true
```

### API Server Mode

```bash
# Start server
qry serve --port 8080 --dir /path/to/repos

# Query via HTTP
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"prompt": "get active users", "dialect": "postgres"}'

# Response
{
  "sql": "SELECT * FROM users WHERE active = true",
  "tables": ["users"],
  "type": "SELECT",
  "safe": true
}
```

### Multi-Repo / Multi-DB Setup

The most secure self-hosted setup:

```
┌─────────────────── Your VPS ───────────────────┐
│                                                 │
│   /repos                                        │
│     ├── backend-api/        (Postgres)         │
│     ├── analytics-service/  (ClickHouse)       │
│     ├── legacy-monolith/    (MySQL)            │
│     └── mobile-backend/     (SQLite)           │
│                                                 │
│   $ qry serve --port 8080 --dir /repos         │
│                                                 │
│   Expose API (internal network / VPN)          │
│                                                 │
└─────────────────────────────────────────────────┘
```

The LLM understands context from file paths, schema definitions, and naming conventions across all repos.

---

## Value Proposition

### What Problem Does QRY Solve?

**The Pain:**
- Writing SQL queries requires understanding the schema
- Schema understanding requires digging through migrations/models
- Context-switching between code and database tools
- Repetitive SQL for common queries

**The Solution:**
- Ask in plain English → Get SQL
- LLM already understands your codebase
- QRY provides the SQL-specific interface

### Jobs To Be Done

#### Core Job Statement
> **"When I need to write a SQL query, I want to describe what I need in plain English, so I can get working SQL without hunting through migrations or documentation."**

#### Functional Jobs
- Generate correct SQL from natural language
- Handle JOINs and relationships correctly
- Produce syntax for my specific database dialect
- Block dangerous queries (guardrails)

#### Emotional Jobs
- Feel confident the SQL is correct
- Feel productive (not wasting time)
- Feel safe (guardrails protect production)

#### Social Jobs
- Help onboard new devs faster
- Enable non-SQL-experts to query

---

## Branding

### Primary Tagline
**"Ask. Get SQL."**

Simple, direct, powerful.

### Alternative Taglines
1. "Three letters. Infinite queries."
2. "Plain English → Perfect SQL"
3. "Your codebase. Your queries."
4. "The minimal query layer"

### Logo Concepts

**Option 1: Minimal**
```
QRY
```
Clean, bold, all caps. Nothing else needed.

**Option 2: Symbol**
```
Q|RY  (vertical bar representing query)
[Q]RY (brackets like code)
```

### Color Palette

**Primary:**
- Electric Blue: #0066FF (queries, primary actions)
- Deep Purple: #6B4FBB (database, depth)

**Secondary:**
- Charcoal: #2D3748 (text, CLI)
- Silver: #A0AEC0 (secondary text)

**Accent:**
- Neon Green: #00FF88 (success states)
- Warning Orange: #FF6B35 (warnings/guardrails)

### Typography

**Logo/Brand:**
- Monospace font (JetBrains Mono, Fira Code)
- All caps: QRY
- Weight: Bold

---

## Domain & Social

### Primary Domain
- **qry.dev** ← PERFECT for this! (developer-focused)

### Alternatives
- qry.io
- qry.sh (shell/CLI vibe)

### Social Handles
- GitHub: @qry or @qrydev
- Twitter/X: @qry or @qrydev

---

## Open Source Strategy

### Priority: Traction First

```
Phase 1: Ship CLI (NOW)
├── Core CLI functionality
├── Multi-backend support
├── JSON structured output
├── Basic guardrails
└── Great README + demo GIF

Phase 2: Build Community
├── Respond to issues fast
├── Accept contributions
├── Blog posts / tutorials
├── Show HN launch
└── Target: 1k+ stars

Phase 3: Expand Features
├── API server mode
├── More guardrails
├── Query history
├── Web UI (optional)
└── Integrations (Slack, n8n, etc.)
```

### MVP Feature List

| Feature | Priority |
|---------|----------|
| CLI wrapper (claude/gemini backend) | P0 |
| Structured JSON output | P0 |
| Multi-backend switching | P0 |
| `qry serve` API mode | P1 |
| Basic guardrails (block DROP/DELETE) | P1 |
| Dialect selection (postgres/mysql/sqlite) | P1 |
| Config file support | P2 |
| Query history | P2 |

---

## README Structure

```markdown
# QRY

> Ask. Get SQL.

[Demo GIF showing: qry "get active users" → SQL output]

## Install

# Go
go install github.com/username/qry@latest

# Homebrew
brew install qry

# Binary
curl -fsSL https://qry.dev/install.sh | sh

## Quick Start

cd your-project
qry "show me active users from last week"

## Features

- ✅ **Natural Language to SQL** - Ask in English, get SQL
- ✅ **Multi-Backend** - Works with Claude, Gemini, Cursor
- ✅ **Structured Output** - JSON with metadata
- ✅ **Guardrails** - Block destructive queries
- ✅ **Self-Hosted API** - `qry serve` for team use
- ✅ **Multi-DB Support** - Postgres, MySQL, SQLite
- ✅ **Open Source** - MIT licensed

## Usage

# Basic query
qry "get users who signed up this week"

# JSON output
qry --json "get active orders"

# Choose backend
qry --backend=gemini "count posts by user"

# SQL dialect
qry --dialect=mysql "get recent products"

# API server mode
qry serve --port 8080

## Self-Hosted Setup

Best for teams - secure, no data leaves your infra:

# Clone your repos
mkdir /repos && cd /repos
git clone your-backend
git clone your-analytics

# Start QRY server
qry serve --port 8080 --dir /repos

# Query via API
curl -X POST http://localhost:8080/query \
  -d '{"prompt": "get active users"}'

## Contributing

[Contribution guidelines]

## License

MIT
```

---

## Launch Strategy

### Phase 1: Soft Launch (Week 1)

```
🚀 Introducing QRY

Three letters. Infinite SQL queries.

$ qry "show active users from last week"
→ Perfect SQL, instantly.

✓ Wraps Claude/Gemini CLI
✓ Structured JSON output
✓ Guardrails built-in
✓ Self-hostable API

Open source: github.com/username/qry

[Demo GIF]
```

**Where to post:**
- Twitter/X
- LinkedIn
- Dev.to

### Phase 2: HackerNews (Week 2)

```
Show HN: QRY – Natural language to SQL CLI (wraps Claude/Gemini)

Hi HN! I built QRY because I wanted a simple way to generate SQL from natural language.

Instead of building yet another NL2SQL engine, QRY wraps existing LLM CLIs (Claude, Gemini) that already understand your codebase perfectly.

What QRY adds:
- SQL-optimized prompting
- Structured JSON output (tables, query type, safety check)
- Guardrails (block DROP/DELETE by default)
- Multi-backend (switch between Claude/Gemini easily)
- Self-hosted API server mode

Example:
$ qry "get users who signed up last week with more than 5 posts"
→ SELECT u.* FROM users u 
  JOIN posts p ON p.user_id = u.id 
  WHERE u.created_at >= CURRENT_DATE - INTERVAL '7 days'
  GROUP BY u.id HAVING COUNT(p.id) > 5

$ qry --json "get active users"
→ {"sql": "...", "tables": ["users"], "type": "SELECT", "safe": true}

Built in Go. Single binary. MIT licensed.

https://github.com/username/qry

Looking for feedback!
```

### Phase 3: Reddit (Week 3)
- r/golang
- r/programming
- r/commandline
- r/selfhosted

---

## Success Metrics

### GitHub Stars
- Month 1: 500 stars
- Month 3: 2,000 stars
- Month 6: 5,000 stars
- Year 1: 10,000 stars

### Users
- Month 1: 100 CLI installs
- Month 3: 1,000 CLI installs
- Month 6: 5,000 CLI installs

---

## Future Possibilities (Not MVP)

Once traction is established, potential expansions:

- **Web UI** - Browser-based interface
- **Slack Bot** - `@qry get active users`
- **n8n/Zapier nodes** - Automation integrations
- **VS Code extension** - Inline SQL generation
- **Team features** - Shared query library
- **Analytics** - Query history, usage stats

---

## Potential Challenges & Solutions

### Challenge 1: Dependency on LLM CLIs
**Problem:** QRY depends on Claude/Gemini CLI being installed

**Solution:** 
- Clear installation docs
- Support multiple backends
- Graceful fallback/error messages

### Challenge 2: Output Quality
**Problem:** LLM might generate incorrect SQL

**Solution:**
- SQL-optimized prompts (iterate and improve)
- Dry-run mode for verification
- Community feedback loop

### Challenge 3: Pronunciation
**Problem:** Is it "Q-R-Y" or "query"?

**Solution:** 
- Official: "Just say 'query'"
- FAQ on website/README

---

## Brand Voice

### Personality
- **Minimal** - No fluff, just facts
- **Smart** - Technically competent
- **Helpful** - User-focused

### Writing Style
- Short sentences
- Active voice
- Technical but accessible
- No marketing BS

### Examples

❌ Bad:
"QRY revolutionizes the paradigm of database query generation through cutting-edge AI technology..."

✅ Good:
"QRY generates SQL from English. That's it."

---

## Next Steps

1. ✅ Validate idea
2. ⬜ Set up Go project structure
3. ⬜ Implement core CLI (claude backend first)
4. ⬜ Add JSON structured output
5. ⬜ Add guardrails
6. ⬜ Create demo GIF
7. ⬜ Write README
8. ⬜ Launch on GitHub
9. ⬜ Post on HN

---

# QRY

**Three letters.**
**Infinite queries.**

Let's build this! 🚀
