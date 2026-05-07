# zshrc-manager

![demo](demo.gif)

A Go CLI tool for managing your `.zshrc` file — list aliases, env vars, and functions; search, add, remove, and back up safely.

Built as a Go learning project.

## Install

### Option 1: Go install (requires Go 1.21+)

```bash
go install github.com/Dan925/zshrc_manager@latest
```

The binary is placed in `~/go/bin/`. Make sure that's on your PATH:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Option 2: Download a release binary (no Go required)

1. Go to [Releases](https://github.com/Dan925/zshrc_manager/releases/latest)
2. Download the archive for your OS and architecture (e.g. `zshrc-manager_linux_amd64.tar.gz`)
3. Extract and move to your PATH:

```bash
tar -xzf zshrc-manager_linux_amd64.tar.gz
mv zshrc ~/.local/bin/
```

### Option 3: Build from source

```bash
git clone https://github.com/Dan925/zshrc_manager.git
cd zshrc_manager
go build -o zshrc .
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

### Env vars

```bash
zshrc env list                      # EDITOR -> nvim
zshrc env add EDITOR=nvim           # add or overwrite, diff + confirm
zshrc env remove EDITOR             # diff + confirm
zshrc env add EDITOR=nvim --dry-run
```

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
│   ├── env.go
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
