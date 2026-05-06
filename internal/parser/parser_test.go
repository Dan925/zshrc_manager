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

func TestParseFunctions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshrc")
	content := "function greet() {\n\techo \"Hello, $1\"\n}\n\nmkcd() {\n\tmkdir -p \"$1\" && cd \"$1\"\n}\n"
	os.WriteFile(path, []byte(content), 0644)

	zf, err := NewParser(path).Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(zf.Functions) != 2 {
		t.Fatalf("expected 2 functions, got %d: %+v", len(zf.Functions), zf.Functions)
	}
	if zf.Functions[0].Name != "greet" {
		t.Errorf("Functions[0].Name = %q, want %q", zf.Functions[0].Name, "greet")
	}
	if zf.Functions[1].Name != "mkcd" {
		t.Errorf("Functions[1].Name = %q, want %q", zf.Functions[1].Name, "mkcd")
	}
}

func TestParseAliases(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshrc")
	content := "alias gs='git status'\nalias gp=\"git push\"\nalias ll=ls -la\n# not an alias\nexport FOO=bar\n"
	os.WriteFile(path, []byte(content), 0644)

	zf, err := NewParser(path).Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(zf.Aliases) != 3 {
		t.Fatalf("expected 3 aliases, got %d", len(zf.Aliases))
	}

	cases := []struct{ name, value string }{
		{"gs", "git status"},
		{"gp", "git push"},
		{"ll", "ls -la"},
	}
	for i, c := range cases {
		if zf.Aliases[i].Name != c.name {
			t.Errorf("Aliases[%d].Name = %q, want %q", i, zf.Aliases[i].Name, c.name)
		}
		if zf.Aliases[i].Value != c.value {
			t.Errorf("Aliases[%d].Value = %q, want %q", i, zf.Aliases[i].Value, c.value)
		}
	}
}
