# Project Discovery: `zshrc-manager` — A Go CLI Learning Project

---

## Executive Summary

This project is a Go CLI tool for managing `.zshrc` files — organizing aliases, managing named sections, searching functions, and backing up versions. The primary goal is **learning Go** through a real, personally useful tool. The `.zshrc` management domain is an excellent choice: it involves file I/O, string parsing, error handling, user interaction, and structured data — all without requiring concurrency or network complexity in the MVP, keeping the learning surface focused.

---

## Clarifying Questions

### Learning Goals

1. Are you brand new to Go, or do you have some exposure (e.g., you've read parts of the Go Tour, written small programs)?
2. Is there a specific Go concept you most want to understand — interfaces, error handling, goroutines, or the type system?
3. Do you want the project to eventually touch concurrency (e.g., parallel backup operations), or stay single-threaded for now?
4. How important is test coverage as a learning objective? Writing idiomatic Go tests is a distinct skill worth targeting explicitly.
5. Do you want to publish this as a real module on pkg.go.dev, or is it purely a personal learning artifact?

### zshrc Pain Points

6. What is your current `.zshrc` like — a flat file with no structure, or do you already use comment-based sections (`# --- Git Aliases ---`)?
7. What is the most painful thing about managing it today?
8. Do you source other files from `.zshrc` (e.g., `source ~/.zsh/aliases.sh`)? Should the tool understand that structure?
9. Do you ever share aliases across machines?
10. Do you use Oh My Zsh, Prezto, or a custom setup?

### Scope & Complexity

11. Should the MVP be purely read-only first (list/search/view) before adding write operations?
12. Should the tool have a TUI (terminal user interface, like `fzf` or `lazygit`) or is a pure CLI sufficient?
13. Do you want shell completions? (Cobra generates these automatically — a useful feature to learn.)

---

## User Personas

**Primary: You (the developer/user)**
A developer who maintains a growing `.zshrc` with aliases, functions, exports, and path modifications accumulated over years. You want to find things faster, stop accidentally breaking the file, and understand what's in it at a glance. You are also learning Go and want the project to teach real patterns, not toy examples.

**Secondary (future): Another developer on a different machine**
Eventually the tool should be portable enough that a colleague could install it and manage their own zshrc. This shapes module naming, config file design, and documentation but is explicitly out of scope for MVP.

---

## User Stories by Epic

### Epic 1: Parsing & Reading

**US-01**
As a user, I want to list all aliases in my `.zshrc` so that I can see what shortcuts I have without reading the whole file.

*Acceptance Criteria:*
- Given a `.zshrc` with `alias gs='git status'`, when I run `zshrc list aliases`, then I see `gs -> git status` printed to stdout.
- Given a file with no aliases, I see a clear "no aliases found" message, not an error.

**US-02**
As a user, I want to list all functions defined in my `.zshrc` so that I can audit what custom logic I've built up.

*Acceptance Criteria:*
- Detects both `function foo() {}` and `foo() {}` syntax.
- Shows function name and the line number where it starts.
- Shows full function body with `--verbose`.

**US-03**
As a user, I want to search for a keyword across my `.zshrc` so that I can find the alias or function I half-remember.

*Acceptance Criteria:*
- `zshrc search docker` returns all lines containing "docker" with line numbers.
- Search is case-insensitive by default, with `--case-sensitive` flag.
- Matches are highlighted in output.

**US-04**
As a user, I want to see my `.zshrc` organized by section so that I understand the structure at a glance.

*Acceptance Criteria:*
- Sections are detected by comment headers matching `# ---`.
- `zshrc outline` prints a table of contents with section names and line ranges.
- Unsectioned lines are grouped under an "Ungrouped" section.

---

### Epic 2: Writing & Editing

**US-05**
As a user, I want to add a new alias so that I don't have to open my editor for a one-liner.

*Acceptance Criteria:*
- `zshrc add alias gs='git status'` appends the alias.
- Before writing, the tool prints a diff and asks for confirmation.
- A backup of the original file is created before any write operation.
- If the alias key already exists, the tool warns and requires `--force` to overwrite.

**US-06**
As a user, I want to remove an alias by name so that I can clean up without hunting for the exact line.

*Acceptance Criteria:*
- `zshrc remove alias gs` finds and removes the line.
- Shows a diff before confirming.
- Creates a backup before writing.

**US-07**
As a user, I want a `--dry-run` flag on all commands so that I can see what would happen before committing to a change.

*Acceptance Criteria:*
- `--dry-run` is a global flag available on every command.
- Write commands show the diff and exit without creating a backup, prompting, or writing.
- Read commands note that no changes would be made.

---

### Epic 3: Backup & Safety

**US-08**
As a user, I want every write operation to create a timestamped backup so that I can recover from mistakes.

*Acceptance Criteria:*
- Backups are stored in `~/.zshrc.bak/YYYY-MM-DD_HH-MM-SS.zshrc`.
- `zshrc backup list` shows all backups with timestamps and file sizes, sorted newest first.
- `zshrc backup restore <timestamp>` restores a specific backup after confirmation.
- Restore creates a backup of the current state before overwriting.

**US-09**
As a user, I want the tool to validate my `.zshrc` for syntax errors so that I don't break my shell.

*Acceptance Criteria:*
- Runs `zsh -n ~/.zshrc` and reports the result.
- `zshrc validate` runs the check without making any changes.

---

### Epic 4: Configuration (v2)

**US-10**
As a user, I want to configure the path to my `.zshrc` so that the tool works if my file is in a non-standard location.

*Acceptance Criteria:*
- Default path is `~/.zshrc`.
- Overridable per-command with `--file <path>`.
- (v2) Configurable via `~/.config/zshrc-manager/config.toml`.

---

## Technical Decision Framework

### 1. CLI Framework

**Decision:** Cobra v1.8

**Why:** Industry standard for Go CLIs. Used by kubectl, helm, gh, and docker. The subcommand pattern, persistent flags, and `RunE` error handling are patterns you'll encounter everywhere in the Go ecosystem.

**What you'll learn:**
- Cobra's `Command` struct and subcommand tree
- `PersistentFlags()` for global flags (`--file`, `--dry-run`)
- `RunE` vs `Run` (returning errors vs. handling them inline)
- `cobra.ExactArgs(n)` for argument validation

**Alternatives considered:**

| Option | Pros | Cons |
|---|---|---|
| `flag` stdlib | Zero deps, teaches Go's stdlib | Manual subcommand routing, not how real tools are built |
| `cobra` | Industry standard, transferable | One external dependency |
| `urfave/cli` | Simpler API | Less commonly used, smaller ecosystem |

---

### 2. Error Handling

**Decision:** `fmt.Errorf` with `%w` + `errors.Is`

**Why:** Go's stdlib error handling is explicit by design. You will write `if err != nil` constantly — this is intentional, not a flaw. The `%w` verb wraps errors so callers can inspect the chain with `errors.Is`.

**Pattern:**
```go
// internal packages: wrap with context, return
return nil, fmt.Errorf("opening %s: %w", path, err)

// cmd layer: catch and print clean message
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}

// check specific error types
if errors.Is(err, os.ErrNotExist) {
    fmt.Fprintf(os.Stderr, "File not found: %s\n", path)
}
```

**What you'll learn:**
- Go's error-as-value philosophy
- Wrapping vs. creating errors
- The `%w` verb vs. `%v`
- When to use `errors.Is` vs. type assertion

---

### 3. File Parsing Strategy

**Decision:** Regex-based line-by-line parsing using Go's `regexp` package

**Why:** Simpler than a state machine for a first Go project. Sufficient for the well-defined patterns in a `.zshrc` (alias lines, function declarations). Can be refactored to a state machine in v2.

**Pattern:**
```go
var aliasRe = regexp.MustCompile(`^alias\s+([^=\s]+)=(.+)$`)

for i, line := range zf.RawLines {
    if m := aliasRe.FindStringSubmatch(line); m != nil {
        // m[1] = name, m[2] = value
    }
}
```

**Key design rule — lossless storage:** Every line in the file is stored in `ZshrcFile.RawLines`. Write operations reconstruct the file from `RawLines` with only the targeted line changed. This ensures no line is ever silently dropped.

**What you'll learn:**
- `regexp.MustCompile` and compile-time safety
- `FindStringSubmatch` and capture groups
- Iterating slices with index and value
- Brace-counting for multi-line function bodies

---

### 4. Project Structure

**Decision:** `cmd/` for CLI layer, `internal/` for logic layer

```
zshrc-manager/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go       # --file, --dry-run flags; Execute()
│   ├── helpers.go    # confirm(), showDiff()
│   ├── list.go
│   ├── search.go
│   ├── add.go
│   ├── remove.go
│   ├── backup.go
│   └── validate.go
└── internal/
    ├── parser/
    │   ├── parser.go
    │   └── parser_test.go
    └── backup/
        ├── backup.go
        └── backup_test.go
```

**Rule:** `cmd/` knows about Cobra and user interaction. `internal/` knows nothing about Cobra — only files and data. This separation makes `internal/` fully testable without spinning up a CLI.

**Why `internal/`?** Go's `internal` package convention prevents code outside the module from importing these packages. It's a strong signal: this is implementation detail, not a public API.

---

### 5. Data Structures

```go
// internal/parser/parser.go

type Alias struct {
    Name  string // "gs"
    Value string // "git status"
    Line  int    // 1-indexed line number in file
}

type Function struct {
    Name      string
    StartLine int
    EndLine   int
    Body      string // populated only with --verbose
}

type ZshrcFile struct {
    Aliases   []Alias
    Functions []Function
    RawLines  []string // every line, preserved for lossless writes
}
```

```go
// internal/backup/backup.go

type Backup struct {
    Path      string
    Timestamp time.Time
    SizeBytes int64
}
```

---

### 6. Testing Strategy

| Layer | Approach | How |
|---|---|---|
| Parser unit tests | Feed string literals into parser, assert on returned structs | `testing.T`, `t.TempDir()` |
| Backup unit tests | Use temp dirs, override HOME env var | `os.Setenv("HOME", t.TempDir())` |
| CLI smoke tests | Build binary, test manually with `--file /tmp/test.zshrc` | `go build`, shell |

**Day-one test:** Write a round-trip test before any feature code. Parse a sample `.zshrc`, call `WriteTo`, assert byte-for-byte equality with the original. This test will protect against regressions in every write operation forever.

---

### 7. Output Formatting

For MVP, `fmt.Printf` is sufficient. Color output can be added later with `github.com/fatih/color` — a single-line change per print statement.

---

## Requirements (MoSCoW)

### Functional Requirements

**Must Have (MVP)**
- FR-01: Parse `.zshrc` and extract all aliases with name, value, and line number
- FR-02: Parse `.zshrc` and extract all function names and line ranges
- FR-03: Search file content by keyword with line numbers and highlighted matches
- FR-04: Add a new alias with pre-write diff and confirmation
- FR-05: Remove an alias by name with pre-write diff and confirmation
- FR-06: Automatic timestamped backup before every write (`~/.zshrc.bak/`)
- FR-07: Restore a backup by timestamp
- FR-08: List all backups with timestamps and sizes
- FR-09: `--dry-run` global flag (show diff, skip backup/write)
- FR-10: `--file` global flag (override default `~/.zshrc`)
- FR-11: `zshrc validate` running `zsh -n` and reporting result

**Should Have**
- FR-12: `--verbose` on `list functions` (show full body)
- FR-13: `--case-sensitive` on `search`
- FR-14: Color output with `fatih/color`
- FR-15: Add a new function (with multi-line body support)

**Could Have**
- FR-16: `zshrc outline` — section table of contents from comment headers
- FR-17: Shell completion generation (`zshrc completion zsh > _zshrc`)
- FR-18: `zshrc doctor` — audit for duplicate aliases, broken sourced files
- FR-19: Config file (`~/.config/zshrc-manager/config.toml`)

**Won't Have (this phase)**
- Remote sync or cloud backup
- TUI interface
- Support for bash or fish shell
- Plugin system

---

### Non-Functional Requirements

**Must Have**
- NFR-01: The tool must never write to `.zshrc` without first creating a backup
- NFR-02: The tool must never silently truncate or corrupt the file — lossless round-trip is mandatory
- NFR-03: All error messages must tell the user what went wrong and what to do about it
- NFR-04: Binary must compile on macOS and Linux

**Should Have**
- NFR-05: Cold-start time under 50ms
- NFR-06: Unit tests for the parser module
- NFR-07: Tests for all write commands using `t.TempDir()`

**Could Have**
- NFR-08: Shell completions for zsh and bash
- NFR-09: Man page generated from Cobra metadata

---

## Risks & Open Decisions

### Open Decisions

| # | Decision | Impact if wrong |
|---|---|---|
| OD-01 | Lossless vs. lossy round-trip — does the tool preserve blank lines, comments, and exact whitespace? | If lossy, the tool will reformat your file on every write. Write the round-trip test first to validate. |
| OD-02 | Scope of zsh syntax support — multiline strings, heredocs, sourced files? | Determines parser complexity. Recommend deferring heredocs and sourced files to v2. |

### Assumptions Made

- You are on macOS or Linux
- Your `.zshrc` is a single file, not a directory of sourced files, at least for MVP
- `zsh -n` is available on your system
- You have Go 1.21+ and the Go toolchain installed (`go`, `gopls`)

### Risks

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Parser scope creep — zsh syntax is surprisingly complex | High | High | Define a strict "supported syntax" list; emit a warning and skip unsupported constructs rather than failing |
| Lossless round-trip harder than expected | Medium | High | Write the round-trip test on day one; it must pass before any write command is implemented |
| Learning Go + building a real tool simultaneously causes frustration | Medium | Medium | Keep Epic 1 (read-only) as the first milestone; you get a useful tool before touching any write logic |

---

## Recommended Next Steps

1. **Set up the project skeleton:** `go mod init github.com/Dan925/zshrc-manager`, add Cobra, write `main.go` and `cmd/root.go`. Verify `go build` and `--help` work.

2. **Write the round-trip test first.** Load a sample `.zshrc`, run it through the parser, call `WriteTo`, assert byte-for-byte equality. This test protects you forever.

3. **Implement Epic 1 (read-only) as your first sprint.** No writes, no backups — just parse and display. You'll have a useful tool after this sprint.

4. **Add the backup system before any write command.** `internal/backup` must exist and be tested before `cmd/add.go` or `cmd/remove.go` are written.

5. **Use `--file /tmp/test.zshrc` during development.** Always test write commands against a copy, never your real `.zshrc`, until you're confident.

---

This document is a living reference. As you implement and discover new constraints, update the requirements backlog. The lossless round-trip constraint and the `cmd/` vs `internal/` separation are the two architectural pillars — everything else can change cheaply once those are solid.
