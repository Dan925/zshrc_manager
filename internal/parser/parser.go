package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)
type EnvVar struct {
	Name string
	Value string
	Line  int
}

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
	EnvVars   []EnvVar
	RawLines  []string
}

type Parser struct {
	filePath string
}

var envRe = regexp.MustCompile(`^export\s+([^=\s]+)=(.+)$`)
var aliasRe = regexp.MustCompile(`^alias\s+([^=\s]+)=(.+)$`)
var funcRe   = regexp.MustCompile(`^(?:function\s+)?([a-zA-Z_][a-zA-Z0-9_:-]*)\s*(?:\(\s*\))?\s*\{`)

var shellKeywords = map[string]bool{
	"if": true, "while": true, "for": true, "case": true,
	"do": true, "then": true, "else": true, "elif": true,
}

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

	var currentFunc *Function
	braceDepth := 0

	for i, line := range zf.RawLines {
		lineNum := i + 1

		if currentFunc != nil {
			currentFunc.Body += line + "\n"
			braceDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if braceDepth <= 0 {
				currentFunc.EndLine = lineNum
				zf.Functions = append(zf.Functions, *currentFunc)
				currentFunc = nil
				braceDepth = 0
			}
			continue
		}

		if m := aliasRe.FindStringSubmatch(line); m != nil {
			zf.Aliases = append(zf.Aliases, Alias{
				Name:  strings.TrimSpace(m[1]),
				Value: strings.Trim(strings.TrimSpace(m[2]), `'"`),
				Line:  lineNum,
			})
			continue
		}

		if m := funcRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			if shellKeywords[name] {
				continue
			}
			currentFunc = &Function{
				Name:      name,
				StartLine: lineNum,
				Body:      line + "\n",
			}
			braceDepth = strings.Count(line, "{") - strings.Count(line, "}")
			if braceDepth <= 0 {
				braceDepth = 1
			}
		}
		
		if m:= envRe.FindStringSubmatch(line); m!= nil {
			zf.EnvVars = append(zf.EnvVars, EnvVar{
				Name: strings.TrimSpace(m[1]),
				Value: strings.Trim(strings.TrimSpace(m[2]),`'"`),
				Line: lineNum,
			})
			continue
		}
	}

	return zf, nil
}

func (zf *ZshrcFile) AddEnvVar(name, value string, force bool) error{
	newRawLine := fmt.Sprintf("export %s=%s", name, value)
	for _, v:= range zf.EnvVars {
		if v.Name == name {
			if force {
				zf.RawLines[v.Line-1] = newRawLine //overwrite line if env variable exists
				return nil
			} else {
				return fmt.Errorf("Env variable %q already exist, use the --force flag to overwrite", name)
			}
		}
	}

	zf.RawLines = append(zf.RawLines,newRawLine)
	return nil
}

func (zf *ZshrcFile) RemoveEnvVar(name string) error {
	for _, v:= range zf.EnvVars {
		if v.Name == name {
			zf.RawLines = append(zf.RawLines[:v.Line-1],zf.RawLines[v.Line:]...)
			return nil
		}

	}
	return fmt.Errorf("env var %q not found", name)

}

func (zf *ZshrcFile) WriteTo(path string) error {
	content := strings.Join(zf.RawLines, "\n")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
