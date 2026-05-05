# Project Discovery: `zshrc-manager` — A Rust CLI Learning Project

---

## Executive Summary

This project is a Rust CLI tool for managing `.zshrc` files — organizing aliases, managing named sections, searching functions, and backing up versions. The primary goal is **learning Rust** through a real, personally useful tool. The `.zshrc` management domain is an excellent choice: it involves file I/O, string parsing, error handling, user interaction, and structured data — all without requiring a network stack or async complexity in the MVP, keeping the learning surface focused.

---

## User Clarifications (Answered)

- **Rust experience:** Brand new (no prior experience)
- **Parser strategy:** Lossy parsing acceptable (reformatting is OK)
- **MVP scope:** Include write operations from day one (not read-only first)

These answers significantly simplify the implementation:
1. Start with regex-based parsing instead of state machines (simpler for a beginner, refactor later if needed)
2. Don't worry about preserving exact formatting
3. Build backup system first (critical safety net for writes)

---

## User Personas

**Primary: You (the developer/user)**
A developer learning Rust who maintains a growing `.zshrc` with aliases, functions, exports, and path modifications. You want to find things faster, stop accidentally breaking the file, and understand what's in it at a glance.

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
- Highlight matches in output (ANSI colors)

---

### Epic 2: Backup & Safety (Implement Before Writes!)

**US-04: Create backup before write**
As a user, I want every write operation to create a timestamped backup so that I can recover from mistakes.

*Acceptance Criteria:*
- Backups stored in `~/.zshrc.bak/YYYY-MM-DD_HH-MM-SS.zshrc`
- Backup created atomically before any modification
- Directory created automatically if missing

**US-05: List backups**
As a user, I want to see all available backups with timestamps and file sizes.

*Acceptance Criteria:*
- `zshrc backup list` shows all backups sorted by timestamp
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
- Does not remove if alias name appears in function body (with warning)

**US-09: Add function**
As a user, I want to add a custom function to my `.zshrc`.

*Acceptance Criteria:*
- `zshrc add function myfunction 'echo hello'` (or multiline with stdin)
- Handles both single-line and multi-line functions
- Creates backup before writing

---

### Epic 4: Validation & Safety

**US-10: Validate syntax**
As a user, I want the tool to check my `.zshrc` for syntax errors before and after edits.

*Acceptance Criteria:*
- `zshrc validate` runs `zsh -n ~/.zshrc`
- Shows clear pass/fail status
- If post-edit validation fails, offer automatic rollback to backup

---

## Technical Decisions (Adjusted for Beginner + Lossy Parsing)

### 1. CLI Framework: clap builder API

**Why:** You need to understand what's happening. The builder API is explicit and teachable. Derive macros come later once you see the patterns.

**What you'll learn:**
- How Rust's type system represents command-line structures
- Builder pattern in Rust
- How argument parsing works under the hood

---

### 2. Error Handling: `anyhow` + simple string errors

**Why:** As a beginner, `anyhow::Context` keeps errors simple while still teaching good practice. `thiserror` can come in version 2.

**What you'll learn:**
- How to add context to errors
- The `?` operator and error propagation
- When to panic vs. return errors

---

### 3. File Parsing: Regex + line-by-line

**Why:** Simpler than a state machine. You can refactor to a proper parser once you understand the domain.

**What you'll learn:**
- Rust regex syntax and the `regex` crate
- String manipulation and pattern matching
- Iterator patterns (`.lines()`, `.filter()`, `.map()`)

**Future refactor:** Once this works, replace regex parser with a real parser to handle edge cases.

---

### 4. Backup System: Simple file copy with timestamps

**Why:** Uses only `std::fs` and `chrono`. No external complexity.

**What you'll learn:**
- File I/O in Rust (`std::fs`)
- Path manipulation (`std::path`)
- Creating directories atomically
- Timestamps and formatting (chrono)

---

### 5. Configuration: Not in MVP

**Why:** Keep scope tight. Start with hardcoded `~/.zshrc` path and `# ---` section pattern.

**When to add:** Once you can add/remove/search reliably.

---

## Recommended Crate List (MVP Only)

```toml
[dependencies]
clap = "4"                    # CLI argument parsing
anyhow = "1"                  # Error handling with context
regex = "1"                   # Pattern matching for parsing
chrono = "0.4"                # Timestamps for backups
colored = "2"                 # Colored terminal output (optional, nice UX)

[dev-dependencies]
tempfile = "3"                # Temporary files for testing
```

Later additions (v2):
- `thiserror` — typed domain errors
- `serde` + `toml` — config files
- `comfy-table` — table output
- `similar` — better diffs
- `assert_cmd` — CLI integration tests

---

## Implementation Roadmap

### Phase 1: Setup & Parsing (Week 1)
1. `cargo new zshrc-manager`
2. Set up clap with basic `--help` and subcommands: `list`, `search`, `add`, `remove`, `backup`
3. Implement simple regex parser to extract aliases and functions
4. Implement `list aliases` and `list functions` commands
5. Write unit tests for parser

**Milestone:** Can run `zshrc list aliases` and see your real aliases printed.

---

### Phase 2: Read Operations (Week 2)
1. Implement `search` command with case-insensitive matching
2. Add colored output for better UX
3. Add unit tests for search
4. Test against your actual `.zshrc`

**Milestone:** Can find any alias or function by keyword.

---

### Phase 3: Backup System (Week 3) — Implement BEFORE any writes
1. Create backup infrastructure: `~/.zshrc.bak/` directory
2. Implement `backup list` and `backup restore` commands
3. Add atomic backup creation (create temp file, rename atomically)
4. Test backup/restore with temp files

**Milestone:** Can safely recover from any mistake.

---

### Phase 4: Write Operations (Week 4)
1. Implement `add alias` with pre-write diff and confirmation
2. Implement `remove alias` with safety checks
3. Add post-write validation (`zsh -n`)
4. Test with temporary `.zshrc` copies

**Milestone:** Can modify your `.zshrc` safely with automatic backups.

---

### Phase 5: Polish & Testing (Week 5)
1. Write comprehensive integration tests
2. Add error messages for edge cases
3. Test with your real `.zshrc`
4. Document commands and examples

**Milestone:** Production-ready tool for personal use.

---

## Learning Outcomes

By the end of this project, you'll understand:

1. **File I/O & Paths** — reading, writing, creating directories
2. **String Processing** — regex, pattern matching, working with text
3. **CLI Design** — parsing arguments, user confirmation, exit codes
4. **Error Handling** — context, propagation, recovery
5. **Testing** — unit tests, integration tests with temp files
6. **Iterators** — filtering, mapping, collecting results
7. **Type System** — using Rust's type system to make bad states unrepresentable

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

### Should Have
- Add new function (multi-line support)
- Case-sensitive search flag
- Colored output
- Helpful error messages
- Unit + integration tests

### Could Have
- Move alias to section
- `zshrc doctor` (find duplicates, unused functions)
- Config file for custom paths
- Shell completions

### Won't Have (v1)
- Remote backup sync
- TUI interface
- Support for bash/fish
- Plugin system

---

## Getting Started

1. Create the project: `cargo new zshrc-manager && cd zshrc-manager`
2. Add dependencies to `Cargo.toml`
3. Start with Phase 1: Parse your `.zshrc` and print aliases
4. Commit early and often
5. Write tests as you go, don't leave them for the end

This is a real, useful tool you'll actually use. That's what makes it a great learning project.
