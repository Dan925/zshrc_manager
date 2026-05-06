package cmd

import (
	"fmt"
	"strings"

	bk "github.com/Dan925/zshrc-manager/internal/backup"
	"github.com/Dan925/zshrc-manager/internal/parser"
	"github.com/spf13/cobra"
)

var forceAdd bool

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an entry to your .zshrc",
}

var addAliasCmd = &cobra.Command{
	Use:   "alias <key>=<value>",
	Short: "Add a new alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("expected format: key=value, got %q", args[0])
		}
		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		newLine := fmt.Sprintf("alias %s=%s", name, value)

		zf, err := parser.NewParser(filePath).Parse()
		if err != nil {
			return err
		}

		for _, a := range zf.Aliases {
			if a.Name == name {
				if !forceAdd {
					return fmt.Errorf("alias %q already exists (use --force to overwrite)", name)
				}
				zf.RawLines = append(zf.RawLines[:a.Line-1], zf.RawLines[a.Line:]...)
				break
			}
		}

		fmt.Println("Changes:")
		showDiff("", newLine)
		fmt.Println()

		if dryRun {
			fmt.Println("Dry run — no changes made.")
			return nil
		}

		if !confirm("Apply this change?") {
			fmt.Println("Aborted.")
			return nil
		}

		backupPath, err := bk.Create(filePath)
		if err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
		fmt.Printf("Backup created: %s\n", backupPath)

		n := len(zf.RawLines)
		if n > 0 && zf.RawLines[n-1] == "" {
			zf.RawLines[n-1] = newLine
			zf.RawLines = append(zf.RawLines, "")
		} else {
			zf.RawLines = append(zf.RawLines, newLine)
		}

		if err := zf.WriteTo(filePath); err != nil {
			return err
		}
		fmt.Printf("Added: %s -> %s\n", name, value)
		return nil
	},
}

func init() {
	addAliasCmd.Flags().BoolVar(&forceAdd, "force", false, "overwrite existing alias")
	addCmd.AddCommand(addAliasCmd)
	rootCmd.AddCommand(addCmd)
}
