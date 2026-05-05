# Project Discovery: `zshrc-manager` — A Rust CLI Learning Project

---

## Executive Summary

This project is a Rust CLI tool for managing `.zshrc` files — organizing aliases, managing named sections, searching functions, and backing up versions. The primary goal is **learning Rust** through a real, personally useful tool. The `.zshrc` management domain is an excellent choice: it involves file I/O, string parsing, error handling, user interaction, and structured data — all without requiring a network stack or async complexity in the MVP, keeping the learning surface focused.

The document below assumes no prior answer from the user on any of the clarifying questions. It is structured so you can read the full document now and answer the questions as you go, then use the requirements section as your implementation backlog.

---

## Clarifying Questions

### Learning Goals

1. Are you brand new to Rust, or do you have some exposure (e.g., you've read parts of the book, written small programs)?
2. Is there a specific Rust concept you most want to crack — ownership and borrowing, traits and generics, error handling, iterators, or something else?
3. Do you want the project to eventually touch async/await (e.g., syncing to a remote backup), or would you prefer to stay synchronous for now?
4. How important is test coverage as a learning objective? Writing idiomatic Rust tests is a distinct skill worth targeting explicitly.
5. Do you want to publish this as a real crate on crates.io, or is it purely a personal learning artifact?

### zshrc Pain Points

6. What is your current `.zshrc` like — a flat file with no structure, or do you already use comment-based sections (`# --- Git Aliases ---`)?
7. What is the most painful thing about managing it today? (Can't find an alias? Fearful of breaking something? Duplicates accumulating?)
8. Do you source other files from `.zshrc` (e.g., `source ~/.zsh/aliases.sh`)? Should the tool understand that structure?
9. Do you ever share aliases across machines? If so, is syncing between machines a use case you'd want eventually?
10. Do you use Oh My Zsh, Prezto, or a custom setup? This affects what the file looks like and what lines the tool should treat as "managed vs. unmanaged."

### Scope & Complexity

11. Should the MVP be purely read-only first (list/search/view) before adding write operations? This is a safer learning ramp.
12. How do you feel about the tool modifying your `.zshrc` directly versus writing to a staging file you review and apply?
13. Should the tool have a TUI (terminal user interface, like `fzf` or `lazygit`) or is a pure command-line interface sufficient?
14. Do you want shell completions generated from the CLI (a very Rust-native feature worth learning)?

---

## User Personas

**Primary: You (the developer/user)**
A developer who maintains a growing `.zshrc` with aliases, functions, exports, and path modifications accumulated over years. You want to find things faster, stop accidentally breaking the file, and understand what's in it at a glance. You are also learning Rust and want the project to teach real patterns, not toy examples.

**Secondary (future): Another developer on a different machine**
Eventually the tool should be portable enough that a colleague could install it and manage their own zshrc. This shapes naming conventions, config file design, and documentation but is explicitly out of scope for MVP.

---

## User Stories by Epic

### Epic 1: Parsing & Reading

**US-01**
As a user, I want to list all aliases in my `.zshrc` so that I can see what shortcuts I have without reading the whole file.

*Acceptance Criteria:*
- Given a `.zshrc` with `alias gs='git status'`, when I run `zshrc list aliases`, then I see `gs -> git status` printed to stdout.
- Given a file with no aliases, then I see a clear "no aliases found" message, not an error.
- Output supports `--json` flag for piping to other tools.

**US-02**
As a user, I want to list all functions defined in my `.zshrc` so that I can audit what custom logic I've built up.

*Acceptance Criteria:*
- Detects both `function foo() {}` and `foo() {}` syntax.
- Shows function name and the line number where it starts.
- Optionally shows the full function body with `--verbose`.

**US-03**
As a user, I want to search for a keyword across my `.zshrc` so that I can find the alias or function I half-remember.

*Acceptance Criteria:*
- `zshrc search docker` returns all lines containing "docker" with line numbers and context.
- Search is case-insensitive by default, with `--case-sensitive` flag.
- Matches are highlighted in output (using ANSI codes).

**US-04**
As a user, I want to see my `.zshrc` organized by section so that I understand the structure at a glance.

*Acceptance Criteria:*
- Sections are detected by comment headers matching a configurable pattern (default: `# ---`).
- `zshrc outline` prints a table of contents with section names and line ranges.
- Unsectioned lines are grouped under an "Ungrouped" section.

---

### Epic 2: Writing & Editing

**US-05**
As a user, I want to add a new alias so that I don't have to open my editor for a one-liner.

*Acceptance Criteria:*
- `zshrc add alias gs='git status'` appends the alias under the correct section (or a default "Aliases" section).
- Before writing, the tool prints a diff and asks for confirmation.
- A backup of the original file is created before any write operation.
- If the alias key already exists, the tool warns and requires `--force` to overwrite.

**US-06**
As a user, I want to remove an alias by name so that I can clean up without hunting for the exact line.

*Acceptance Criteria:*
- `zshrc remove alias gs` finds and removes the line.
- Shows a diff before confirming.
- Does not remove if the alias name appears in a function body (with a warning).

**US-07**
As a user, I want to move an alias into a named section so that I can organize my file over time.

*Acceptance Criteria:*
- `zshrc move alias gs --section git` moves the alias line to the named section.
- Creates the section header if it doesn't exist.
- Preserves all other lines exactly.

---

### Epic 3: Backup & Safety

**US-08**
As a user, I want every write operation to create a timestamped backup so that I can recover from mistakes.

*Acceptance Criteria:*
- Backups are stored in `~/.zshrc.bak/YYYY-MM-DD_HH-MM-SS.zshrc`.
- `zshrc backup list` shows all backups with timestamps and file sizes.
- `zshrc backup restore <timestamp>` restores a specific backup after confirmation.

**US-09**
As a user, I want the tool to validate my `.zshrc` for obvious syntax errors before and after edits so that I don't break my shell.

*Acceptance Criteria:*
- Runs `zsh -n ~/.zshrc` (zsh syntax check) and reports the result.
- If validation fails post-edit, offers to restore the backup automatically.
- `zshrc validate` runs the check without making any changes.

---

### Epic 4: Configuration

**US-10**
As a user, I want to configure the path to my `.zshrc` so that the tool works if my file is in a non-standard location.

*Acceptance Criteria:*
- Default path is `~/.zshrc`.
- Configurable via `~/.config/zshrc-manager/config.toml`.
- Overridable per-command with `--file <path>`.
- Config file is auto-created with defaults on first run.

**US-11**
As a user, I want to configure the section header pattern so that it matches the comment style I already use.

*Acceptance Criteria:*
- Default pattern is `# ---` prefix.
- Pattern is a regex string in the config file.
- `zshrc config show` prints current config values.

---

## Technical Decision Framework

### 1. CLI Framework

**The decision:** Which crate to use for argument parsing and subcommand structure.

**Options:**

| Option | Pros | Cons | Pedagogical Value |
|---|---|---|---|
| `clap` v4 (derive API) | Industry standard, excellent docs, derive macros feel magical in a good way | Macro-heavy; the magic can obscure what's happening | High — you'll use this pattern everywhere in Rust |
| `clap` v4 (builder API) | Explicit, no macros, easier to understand the underlying model | More verbose | Very high for learning, lower for productivity |
| `argh` | Minimal, fast to learn | Limited ecosystem, less commonly used | Low — not transferable |

**Recommendation:** Start with `clap` v4 **builder API** for the first two epics, then migrate to the derive API. This forces you to understand what the derive macros are generating, which pays off when debugging.

---

### 2. Error Handling

**The decision:** How to represent and propagate errors throughout the application.

**Options:**

| Option | Pros | Cons | Pedagogical Value |
|---|---|---|---|
| `anyhow` | Dead simple, great for application code, excellent error context chaining | Hides the type system; you learn less | Medium — good for "just ship it" mode |
| `thiserror` | Teaches you to define typed errors; compose well with the type system | More boilerplate | Very high — this is how libraries are written |
| `anyhow` + `thiserror` combined | `thiserror` for domain errors, `anyhow` in `main` | Slightly complex two-layer model | Highest — this is the idiomatic production pattern |

**Recommendation:** Use `thiserror` for your domain errors (parsing errors, validation errors) and `anyhow` in `main` and command handlers. This teaches the right mental model: libraries expose typed errors, binaries use dynamic dispatch.

---

### 3. File Parsing Strategy

**The decision:** How to parse `.zshrc` into a structured representation.

**Options:**

| Option | Pros | Cons | Pedagogical Value |
|---|---|---|---|
| Line-by-line with regex (`regex` crate) | Simple, teaches regex in Rust | Fragile for complex zsh syntax | Medium |
| Hand-written recursive descent parser | Maximum control, teaches parser design | Significant complexity, full zsh syntax is vast | High but risky for scope |
| Streaming line parser with state machine | Balanced — models real parsers, handles multiline functions | Requires careful state design | Very high — state machines are fundamental |

**Recommendation:** Build a **streaming line-based state machine parser**. You define an enum for parser state (`InFunction`, `InAlias`, `InSection`, `Freeform`), process the file line by line, and emit structured tokens. This teaches:
- Rust enums with data (ADTs)
- `match` exhaustiveness
- Iterator-based processing
- How to represent a document as a lossless data structure (preserving comments, blank lines)

The key insight is to build a **lossless AST**: every line in the file maps to a node in your structure, so you can serialize back to an identical file and only the changed lines differ. This is critical for safe write operations.

---

### 4. Configuration File

**The decision:** Format and parsing library for tool configuration.

**Options:**

| Option | Pros | Cons |
|---|---|---|
| TOML via `toml` crate | First-class Rust ecosystem citizen, used by Cargo itself | Slightly verbose for simple configs |
| TOML via `config` crate | Layered config (file + env vars + CLI flags) | More abstraction than needed for MVP |
| JSON via `serde_json` | Already know it, universal | Poor config UX (no comments) |

**Recommendation:** `toml` crate with `serde` derives. This teaches derive macros for serialization/deserialization, which transfers directly to working with JSON APIs, databases, and protocol buffers later.

---

### 5. Testing Strategy

**The decision:** How to test a tool that reads and writes real files.

**Recommendations by layer:**

| Layer | Approach | Crate |
|---|---|---|
| Parser unit tests | Feed string literals into parser, assert on output structures | Built-in `#[test]` |
| Write operation tests | Use `tempfile` crate to create temporary `.zshrc` fixtures | `tempfile` |
| CLI integration tests | Use `assert_cmd` to invoke the binary as a subprocess and assert stdout/stderr/exit code | `assert_cmd` + `predicates` |
| Snapshot tests | Capture parser output for a complex `.zshrc` fixture and diff against stored snapshots | `insta` |

`insta` is particularly valuable for learning because it shows you what your parser is actually emitting and lets you "accept" snapshots with a single command, which builds intuition fast.

---

### 6. Output Formatting

**The decision:** How to format terminal output (colors, tables, diffs).

**Recommendations:**

| Concern | Crate | Why |
|---|---|---|
| Colors and ANSI styles | `owo-colors` or `colored` | Lightweight, `owo-colors` is zero-cost |
| Table output | `comfy-table` | Simple API, good defaults |
| Diffs before writes | `similar` | Excellent diff algorithm, used by `cargo` itself |
| Spinner/progress | `indicatif` | Teaches `Display` trait implementations |

---

## Requirements (MoSCoW)

### Functional Requirements

**Must Have (MVP)**
- FR-01: Parse `.zshrc` and extract all aliases with name and value
- FR-02: Parse `.zshrc` and extract all function names and line ranges
- FR-03: Search file content by keyword with line numbers and context
- FR-04: Display section outline (detected from comment headers)
- FR-05: Add a new alias with pre-write diff and confirmation
- FR-06: Remove an alias by name with pre-write diff and confirmation
- FR-07: Automatic timestamped backup before every write
- FR-08: Restore a backup by timestamp
- FR-09: List all backups

**Should Have**
- FR-10: `zshrc validate` running `zsh -n` and reporting result
- FR-11: Move an alias to a named section
- FR-12: Config file (`~/.config/zshrc-manager/config.toml`) with path and section pattern settings
- FR-13: `--json` output flag on all read commands
- FR-14: Add a new function (with multi-line body support)

**Could Have**
- FR-15: Shell completion generation (`zshrc completions zsh > _zshrc`)
- FR-16: `zshrc doctor` — audit for duplicate aliases, unused functions, broken sourced files
- FR-17: `zshrc format` — rewrite file with consistent section ordering
- FR-18: Git integration — auto-commit backups to a local git repo
- FR-19: Fuzzy search via `skim` or `fzf` integration

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
- NFR-03: All error messages must tell the user what went wrong and what to do about it (no raw `unwrap` panics in release builds)
- NFR-04: Binary must compile on macOS and Linux

**Should Have**
- NFR-05: Cold-start time under 50ms (no heavy startup costs — no JVM, no Node.js)
- NFR-06: At least 70% unit test coverage on the parser module
- NFR-07: Integration tests for all write commands using temp files

**Could Have**
- NFR-08: Shell completions for zsh and bash
- NFR-09: Man page generated from clap metadata

---

## Risks & Open Decisions

### Open Decisions (resolve before starting)

| # | Decision | Impact if wrong |
|---|---|---|
| OD-01 | Lossless vs. lossy parser — do you preserve blank lines, comments, and exact whitespace? | If lossy, the tool will "reformat" your file on every write, which is unacceptable |
| OD-02 | Scope of zsh syntax support — multiline strings, heredocs, sourced files? | Determines parser complexity; recommend deferring heredocs and sourced files to v2 |
| OD-03 | clap builder vs. derive API starting point | Learning curve differs; builder teaches more but slows initial momentum |

### Assumptions Made

- You are on macOS or Linux (Windows zsh support is not a goal)
- Your `.zshrc` is a single file, not a directory of sourced files, at least for MVP
- `zsh -n` is available on your system for syntax validation
- You are comfortable with the Rust toolchain already installed (`rustup`, `cargo`)
- "Safe by default" means the tool confirms before writing and always backs up — no `--force` on first use

### Risks

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Parser scope creep — zsh syntax is surprisingly complex | High | High | Define a strict "supported syntax" list up front; emit a warning and skip unsupported constructs rather than failing |
| Lossless round-trip harder than expected | Medium | High | Write a round-trip test on day one: parse your real `.zshrc`, serialize it back, diff against original — must be identical |
| Learning Rust + building a real tool simultaneously causes frustration | Medium | Medium | Keep Epic 1 (read-only) as the first milestone; you get a working tool before touching any write logic |
| Over-engineering the parser before validating the CLI UX | Medium | Medium | Stub the parser with regex for `list` commands first; replace with state machine once the CLI structure is proven |

---

## Recommended Crate List

```toml
[dependencies]
clap = { version = "4", features = ["derive"] }   # CLI framework
anyhow = "1"                                        # Error handling in main
thiserror = "2"                                     # Typed domain errors
serde = { version = "1", features = ["derive"] }   # Serialization base
toml = "0.8"                                        # Config file parsing
serde_json = "1"                                    # --json output
similar = "2"                                       # Diffs before writes
owo-colors = "4"                                    # Terminal colors
comfy-table = "7"                                   # Table output
tempfile = "3"                                      # Tests only
chrono = "0.4"                                      # Backup timestamps

[dev-dependencies]
assert_cmd = "2"                                    # CLI integration tests
predicates = "3"                                    # Assertion helpers
insta = "1"                                         # Snapshot testing
```

---

## Recommended Next Steps

1. **Answer the clarifying questions** above — especially OD-01 (lossless vs. lossy) and your current Rust experience level. These two answers will change the implementation plan meaningfully.

2. **Set up the project skeleton** with `cargo new zshrc-manager --bin`, add the crates above, and commit a `src/main.rs` that only prints help text. Verify the binary compiles and runs.

3. **Write the round-trip test first** (before any feature code). Load your actual `.zshrc`, run it through an identity parser that does nothing, write it back, and assert byte-for-byte equality. This test will fail correctly until your parser is truly lossless, and it protects you from regressions forever.

4. **Implement Epic 1 (read-only) as your first sprint.** No writes, no backups, no config — just parse and display. You will have a useful tool after this sprint and a solid foundation of parser types to build on.

5. **Add write operations in Epic 2 only after Epic 1 tests pass.** The backup system (US-08) must be implemented before the first write command (US-05), not after.

---

This document is intended as a living reference. As you implement and discover new constraints in the zsh syntax or your own workflow, update the requirements backlog and close out open decisions. The parser state machine design and the lossless round-trip constraint are the two architectural pillars — everything else can be changed cheaply once those are solid.
