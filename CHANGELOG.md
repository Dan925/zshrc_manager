# Changelog

## v0.0.1 — 2026-05-06

Initial release.

### Features

- **`zshrc list aliases`** — list all aliases in `name -> value` format
- **`zshrc list functions`** — list all functions with line numbers; `--verbose` shows full body
- **`zshrc search <keyword>`** — case-insensitive search with highlighted matches; `--case-sensitive` flag available
- **`zshrc add alias <key>=<value>`** — add a new alias with diff preview and confirmation; `--force` to overwrite existing
- **`zshrc remove alias <name>`** — remove an alias by name with diff preview and confirmation
- **`zshrc backup list`** — list all timestamped backups with sizes
- **`zshrc backup restore <timestamp>`** — restore a backup with confirmation; saves current state first
- **`zshrc validate`** — run `zsh -n` syntax check and report pass/fail
- **`--dry-run`** global flag — preview any command without writing changes
- **`--file`** global flag — target a different file instead of `~/.zshrc`

### Safety

- Every write operation creates a timestamped backup in `~/.zshrc.bak/` before touching the file
- All write commands show a diff and require confirmation
- Lossless file handling — no lines are ever silently dropped
