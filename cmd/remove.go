package cmd

import (
	"fmt"

	bk "github.com/Dan925/zshrc-manager/internal/backup"
	"github.com/Dan925/zshrc-manager/internal/parser"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an entry from your .zshrc",
}

var removeAliasCmd = &cobra.Command{
	Use:   "alias <name>",
	Short: "Remove an alias by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		zf, err := parser.NewParser(filePath).Parse()
		if err != nil {
			return err
		}

		var found *parser.Alias
		for i, a := range zf.Aliases {
			if a.Name == name {
				found = &zf.Aliases[i]
				break
			}
		}
		if found == nil {
			return fmt.Errorf("alias %q not found", name)
		}

		removedLine := zf.RawLines[found.Line-1]

		fmt.Println("Changes:")
		showDiff(removedLine, "")
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

		zf.RawLines = append(zf.RawLines[:found.Line-1], zf.RawLines[found.Line:]...)
		if err := zf.WriteTo(filePath); err != nil {
			return err
		}
		fmt.Printf("Removed alias: %s\n", name)
		return nil
	},
}

func init() {
	removeCmd.AddCommand(removeAliasCmd)
	rootCmd.AddCommand(removeCmd)
}
