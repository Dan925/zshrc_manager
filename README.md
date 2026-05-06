# zshrc-manager

![demo](demo.gif)

A Go CLI tool for managing your `.zshrc` file — list aliases, search, add, remove, and back up safely.

Built as a Go learning project.

## Install

```bash
git clone https://github.com/Dan925/zshrc_manager.git
cd zshrc_manager
go build -o zshrc .
```

Move the binary somewhere on your PATH:

```bash
mv zshrc ~/.local/bin/
```

## Usage

```
zshrc [command] [flags]

Flags:
  --file string   path to zshrc file (default ~/.zshrc)
  --dry-run       show what would happen without making changes
```

### List

```bash
zshrc list aliases              # gs -> git status
zshrc list functions            # greet (line 12)
zshrc list functions --verbose  # show full function body
```

### Search

```bash
zshrc search docker             # case-insensitive, highlights matches
zshrc search Docker --case-sensitive
```

### Add

```bash
zshrc add alias gs='git status'
```

Shows a diff and prompts for confirmation before writing. Creates a backup first.

```
Changes:
+ alias gs=git status

Apply this change? [y/N]:
```

### Remove

```bash
zshrc remove alias gs
```

Same diff + confirm flow as add.

### Backup

```bash
zshrc backup list               # list all backups with timestamps and sizes
zshrc backup restore 2026-05-05_14-32-01
```

Backups are stored in `~/.zshrc.bak/`. Every write operation creates one automatically before touching your file.

### Validate

```bash
zshrc validate                  # runs zsh -n ~/.zshrc
```

### Dry run

Any command accepts `--dry-run` to preview changes without writing:

```bash
zshrc add alias foo='bar' --dry-run
zshrc remove alias foo --dry-run
```

### Use a different file

```bash
zshrc --file ~/dotfiles/.zshrc list aliases
```

## Project structure

```
├── main.go
├── cmd/
│   ├── root.go       # --file, --dry-run flags
│   ├── list.go
│   ├── search.go
│   ├── add.go
│   ├── remove.go
│   ├── backup.go
│   ├── validate.go
│   └── helpers.go    # confirm(), showDiff()
└── internal/
    ├── parser/       # alias + function extraction, lossless read/write
    └── backup/       # timestamped backup create, list, restore
```

## Requirements

- Go 1.21+
- zsh (for `zshrc validate`)
