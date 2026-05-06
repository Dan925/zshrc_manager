package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var caseSensitive bool

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search your .zshrc for a keyword",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := args[0]

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", filePath, err)
		}

		flags := "(?i)"
		if caseSensitive {
			flags = ""
		}
		re, err := regexp.Compile(flags + regexp.QuoteMeta(keyword))
		if err != nil {
			return fmt.Errorf("invalid search pattern: %w", err)
		}

		lines := strings.Split(string(content), "\n")
		found := false
		for i, line := range lines {
			if re.MatchString(line) {
				highlighted := re.ReplaceAllStringFunc(line, func(m string) string {
					return "[" + m + "]"
				})
				fmt.Printf("%d: %s\n", i+1, highlighted)
				found = true
			}
		}
		if !found {
			fmt.Printf("No matches for %q\n", keyword)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().BoolVar(&caseSensitive, "case-sensitive", false, "case-sensitive search")
	rootCmd.AddCommand(searchCmd)
}
