# zshrc-manager demo
# Run with: vhs demo.tape

Output demo.gif

Set Shell zsh
Set FontSize 16
Set Width 800
Set Height 500
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 60ms

# Start with a clean prompt
Type "cd ~/repos/zshrc_manager"
Enter
Sleep 500ms

# List aliases
Type "./zshrc list aliases"
Enter
Sleep 1.5s

# List functions
Type "./zshrc list functions"
Enter
Sleep 1.5s

# Search
Type "./zshrc search wifi"
Enter
Sleep 1.5s

# Dry-run add
Type "./zshrc add alias dev='cd ~/repos' --dry-run"
Enter
Sleep 2s

# Validate
Type "./zshrc validate"
Enter
Sleep 1.5s

# Backup list
Type "./zshrc backup list"
Enter
Sleep 1.5s
