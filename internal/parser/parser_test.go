package parser

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleZshrc = `# My .zshrc
export PATH="$HOME/bin:$PATH"

alias gs='git status'
alias gp='git push'

function greet() {
	echo "Hello, $1"
}
`

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshrc")
	if err := os.WriteFile(path, []byte(sampleZshrc), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewParser(path)
	zf, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	outPath := filepath.Join(dir, ".zshrc.out")
	if err := zf.WriteTo(outPath); err != nil {
		t.Fatalf("WriteTo() error: %v", err)
	}

	original, _ := os.ReadFile(path)
	roundTripped, _ := os.ReadFile(outPath)
	if string(original) != string(roundTripped) {
		t.Errorf("round-trip mismatch\ngot:\n%s\nwant:\n%s", roundTripped, original)
	}
}
