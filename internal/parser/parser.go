package parser

import (
	"fmt"
	"os"
	"strings"
)

type Alias struct {
	Name  string
	Value string
	Line  int // 1-indexed
}

type Function struct {
	Name      string
	StartLine int
	EndLine   int
	Body      string
}

type ZshrcFile struct {
	Aliases   []Alias
	Functions []Function
	RawLines  []string
}

type Parser struct {
	filePath string
}

func NewParser(path string) *Parser {
	return &Parser{filePath: path}
}

// Parse reads the file into a ZshrcFile. RawLines preserves every line
// so the file can be reconstructed exactly with WriteTo.
func (p *Parser) Parse() (*ZshrcFile, error) {
	content, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", p.filePath, err)
	}
	zf := &ZshrcFile{}
	zf.RawLines = strings.Split(string(content), "\n")
	return zf, nil
}

// WriteTo reconstructs the file from RawLines and writes it to path.
func (zf *ZshrcFile) WriteTo(path string) error {
	content := strings.Join(zf.RawLines, "\n")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
