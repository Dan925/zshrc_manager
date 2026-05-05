# Project Discovery: `zshrc-manager` — A Go CLI Learning Project

---

## Executive Summary

This project is a Go CLI tool for managing `.zshrc` files — organizing aliases, managing named sections, searching functions, and backing up versions. The primary goal is **learning Go** through a real, personally useful tool. The `.zshrc` management domain is an excellent choice: it involves file I/O, string parsing, error handling, user interaction, and structured data — all without requiring concurrency or network complexity in the MVP, keeping the learning surface focused.

---

## User Clarifications (Answered)

- **Go experience:** Brand new (no prior experience)
- **Parser strategy:** Regex-based line-by-line parsing (simpler for a beginner, refactor later if needed)
- **MVP scope:** Include write operations from day one (not read-only first)
- **CLI framework:** Cobra + Go stdlib (industry standard, used by kubectl, helm, gh, docker)

These answers significantly simplify the implementation:
1. Use regex-based parsing instead of a state machine (simpler for a beginner)
2. Don't worry about preserving exact formatting
3. Build backup system before write commands (critical safety net)

---

## User Personas

**Primary: You (the developer/user)**
A developer learning Go who maintains a growing `.zshrc` with aliases, functions, exports, and path modifications. You want to find things faster, stop accidentally breaking the file, and understand what's in it at a glance.

---

## User Stories by Epic

### Epic 1: Parsing & Reading

**US-01: List all aliases**
As a user, I want to list all aliases in my `.zshrc` so that I can see what shortcuts I have without reading the whole file.

*Acceptance Criteria:*
- `zshrc list aliases` extracts and displays all `alias key=value` lines
- Format: `key -> value`
- Handle edge cases: quoted values, spaces in names

**US-02: List all functions**
As a user, I want to list all functions defined in my `.zshrc` so that I can audit what custom logic I've built up.

*Acceptance Criteria:*
- Detect both `function foo() {}` and `foo() {}` syntax
- Show function name and starting line number
- Optional `--verbose` flag shows full function body

**US-03: Search**
As a user, I want to search for a keyword across my `.zshrc` so that I can find the alias or function I half-remember.

*Acceptance Criteria:*
- `zshrc search docker` returns all lines containing "docker" with line numbers
- Case-insensitive by default, with `--case-sensitive` flag
- Highlight matches in output

---

### Epic 2: Backup & Safety (Implement Before Writes!)

**US-04: Create backup before write**
As a user, I want every write operation to create a timestamped backup so that I can recover from mistakes.

*Acceptance Criteria:*
- Backups stored in `~/.zshrc.bak/YYYY-MM-DD_HH-MM-SS.zshrc`
- Backup created before any modification
- Directory created automatically if missing

**US-05: List backups**
As a user, I want to see all available backups with timestamps and file sizes.

*Acceptance Criteria:*
- `zshrc backup list` shows all backups sorted by timestamp (newest first)
- Displays file size for each backup
- Shows timestamp in human-readable format

**US-06: Restore backup**
As a user, I want to restore my `.zshrc` from a specific backup.

*Acceptance Criteria:*
- `zshrc backup restore <timestamp>` with confirmation prompt
- Creates a new backup of current state before restoring
- Clear feedback on success

---

### Epic 3: Writing & Editing

**US-07: Add alias**
As a user, I want to add a new alias without opening an editor.

*Acceptance Criteria:*
- `zshrc add alias gs='git status'` appends the alias
- Shows diff before writing with confirmation prompt
- Warns if alias key already exists (requires `--force` to overwrite)
- Creates backup before writing

**US-08: Remove alias**
As a user, I want to remove an alias by name.

*Acceptance Criteria:*
- `zshrc remove alias gs` finds and removes the line
- Shows diff before confirming
- Creates backup before writing

---

### Epic 4: Validation & Safety

**US-09: Validate syntax**
As a user, I want the tool to check my `.zshrc` for syntax errors.

*Acceptance Criteria:*
- `zshrc validate` runs `zsh -n ~/.zshrc`
- Shows clear pass/fail status

---

## Technical Decisions

### 1. CLI Framework: Cobra

**Why:** Industry standard for Go CLIs — used by kubectl, helm, gh, and docker. Learning it transfers directly to real-world Go projects.

**What you'll learn:**
- Cobra's command/subcommand structure
- Persistent flags (available to all subcommands)
- `RunE` for commands that return errors

---

### 2. Error Handling: `fmt.Errorf` + `errors.Is`

**Why:** Go's stdlib error handling is explicit and idiomatic. No external libraries needed.

**What you'll learn:**
- The `if err != nil` pattern
- Wrapping errors with `fmt.Errorf("context: %w", err)`
- Checking error types with `errors.Is`
- When to return errors vs. exit with `os.Exit`

---

### 3. File Parsing: `regexp` + `strings.Split`

**Why:** Go's built-in `regexp` package is sufficient for alias and function detection. No external dependencies.

**What you'll learn:**
- Go's `regexp` package
- `strings.Split` for line-by-line processing
- Struct types and slice operations
- Brace-counting for multi-line function detection

---

### 4. Backup System: `os`, `io`, `time` stdlib

**Why:** Everything needed is in Go's standard library.

**What you'll learn:**
- File I/O (`os.ReadFile`, `os.WriteFile`, `io.Copy`)
- Path manipulation (`path/filepath`)
- Timestamps (`time.Now().Format(layout)`)
- Creating directories (`os.MkdirAll`)

---

### 5. Global Flags: `--file` and `--dry-run`

**Why:** Both flags apply to every command. Cobra's `PersistentFlags()` on the root command makes them available everywhere.

- `--file`: override the default `~/.zshrc` path
- `--dry-run`: show what would happen without making any changes

---

### 6. Configuration: Not in MVP

**Why:** Keep scope tight. Start with hardcoded `~/.zshrc` path.

---

## Recommended Dependencies (MVP Only)

```
github.com/spf13/cobra v1.8   # CLI framework

# Everything else is Go stdlib:
# os, io, bufio, regexp, strings, fmt, path/filepath, time, errors
```

Later additions (v2):
- `github.com/fatih/color` — colored terminal output
- `github.com/olekukonko/tablewriter` — table output
- `github.com/sergi/go-diff` — better diffs

---

## Implementation Roadmap

### Phase 1: Setup & Parsing (Day 1 – Week 1)
1. `go mod init github.com/Dan925/zshrc-manager`
2. Add Cobra, set up `cmd/root.go` with `--file` and `--dry-run` flags
3. Implement `internal/parser` with regex-based alias and function extraction
4. Write round-trip test (parse → reconstruct → assert byte-for-byte equality)
5. Implement `zshrc list aliases` and `zshrc list functions`

**Milestone:** Can run `zshrc list aliases` and see your real aliases printed.

---

### Phase 2: Read Operations (Week 1–2)
1. Implement `zshrc search` with case-insensitive matching
2. Add `--case-sensitive` flag and match highlighting
3. Write unit tests for parser

**Milestone:** Can find any alias or function by keyword.

---

### Phase 3: Backup System (Week 2) — Implement BEFORE any writes
1. Implement `internal/backup` with `Create`, `List`, `Restore`
2. Implement `zshrc backup list` and `zshrc backup restore`
3. Write tests using `t.TempDir()`

**Milestone:** Can safely recover from any mistake.

---

### Phase 4: Write Operations (Week 3)
1. Implement `zshrc add alias` with pre-write diff and confirmation
2. Implement `zshrc remove alias` with safety checks
3. Wire `--dry-run` into all write commands
4. Implement `zshrc validate` (`zsh -n`)

**Milestone:** Can modify your `.zshrc` safely with automatic backups.

---

### Phase 5: Polish & Testing (Week 4)
1. Integration tests for all write commands using temp files
2. Improve error messages for common failure cases
3. Test against your real `.zshrc`

**Milestone:** Production-ready tool for personal use.

---

## Learning Outcomes

By the end of this project, you'll understand:

1. **Go package system** — `internal/` visibility, package organization
2. **Structs and methods** — pointer receivers, constructors by convention
3. **Error handling** — `if err != nil`, `%w` wrapping, `errors.Is`
4. **File I/O** — `os.ReadFile`, `os.WriteFile`, `path/filepath`
5. **Regex in Go** — `regexp.MustCompile`, `FindStringSubmatch`
6. **Testing** — `testing.T`, `t.TempDir()`, table-driven tests
7. **CLI design** — Cobra commands, persistent flags, `RunE`
8. **Slices and maps** — Go's core data structures

---

## MoSCoW Requirements

### Must Have (MVP)
- Extract and list aliases from `.zshrc`
- Extract and list functions from `.zshrc`
- Search `.zshrc` by keyword
- Create automatic backups before writes
- Restore from backup
- Add new alias (with diff & confirmation)
- Remove alias (with diff & confirmation)
- Validate syntax with `zsh -n`
- `--dry-run` global flag
- `--file` global flag

### Should Have
- Add new function (multi-line support)
- `--case-sensitive` search flag
- `--verbose` flag on `list functions`
- Colored output (`fatih/color`)
- Unit + integration tests

### Could Have
- `zshrc outline` (section table of contents)
- Shell completions (`zshrc completion zsh`)
- `zshrc doctor` (find duplicates, unused functions)
- Config file for custom paths

### Won't Have (v1)
- Remote backup sync
- TUI interface
- Support for bash/fish
- Plugin system

---

## Getting Started

1. Initialize: `go mod init github.com/Dan925/zshrc-manager`
2. Add Cobra: `go get github.com/spf13/cobra@v1.8.0`
3. Start with Phase 1: build the skeleton, then parse aliases
4. Write the round-trip test before any write logic
5. Commit early and often — one feature per commit

This is a real, useful tool you'll actually use. That's what makes it a great learning project.
