package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Alias struct {
	Name  string
	Value string
	Line  int
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

var aliasRe = regexp.MustCompile(`^alias\s+([^=\s]+)=(.+)$`)

func NewParser(path string) *Parser {
	return &Parser{filePath: path}
}

func (p *Parser) Parse() (*ZshrcFile, error) {
	content, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", p.filePath, err)
	}
	zf := &ZshrcFile{}
	zf.RawLines = strings.Split(string(content), "\n")

	for i, line := range zf.RawLines {
		lineNum := i + 1
		if m := aliasRe.FindStringSubmatch(line); m != nil {
			zf.Aliases = append(zf.Aliases, Alias{
				Name:  strings.TrimSpace(m[1]),
				Value: strings.Trim(strings.TrimSpace(m[2]), `'"`),
				Line:  lineNum,
			})
		}
	}

	return zf, nil
}

func (zf *ZshrcFile) WriteTo(path string) error {
	content := strings.Join(zf.RawLines, "\n")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
