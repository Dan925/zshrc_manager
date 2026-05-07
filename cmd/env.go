package cmd

import (
	"fmt"
	bk "github.com/Dan925/zshrc_manager/internal/backup"
	"github.com/Dan925/zshrc_manager/internal/parser"
	"github.com/spf13/cobra"
	"strings"
)

var forceEnvAdd bool

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables in your .zshrc",
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		zf, err := parser.NewParser(filePath).Parse()
		if err != nil {
			return err
		}
		if len(zf.EnvVars) == 0 {
			fmt.Println("No environment variables found.")
			return nil
		}
		for _, v := range zf.EnvVars {
			fmt.Printf("%s -> %s\n", v.Name, v.Value)
		}
		return nil

	},
}

var envAddCmd = &cobra.Command{
	Use:   "add <KEY>=<value>",
	Short: "Add or overwrite an env var",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("expected format: key=value, got %q", args[0])
		}
		name := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `'"`)
		zf, err := parser.NewParser(filePath).Parse()

		if err != nil {
			return err
		}

		var oldLine string
		for _, ev := range zf.EnvVars {
			if ev.Name == name {
				oldLine = zf.RawLines[ev.Line-1]
				break
			}
		}

		newLine := fmt.Sprintf("export %s=%s", name, value)
		showDiff(oldLine, newLine)

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

		if err := zf.AddEnvVar(name, value, forceEnvAdd); err != nil {
			return err
		}
		if err := zf.WriteTo(filePath); err != nil {
			return err
		}
		fmt.Printf("Added: %s -> %s\n", name, value)

		return nil

	},
}

var envRemoveCmd = &cobra.Command{
	Use:   "remove <KEY>",
	Short: "Remove environment by name/Key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.Trim(strings.TrimSpace(args[0]), `"'`)
		zf, err := parser.NewParser(filePath).Parse()
		if err != nil {
			return err
		}

		var removedLine string
		for _, ev := range zf.EnvVars {
			if ev.Name == name {
				removedLine = zf.RawLines[ev.Line-1]
				break
			}
		}

		showDiff(removedLine, "")

		if dryRun {
			fmt.Println("Dry run — no changes made.")
			return nil
		}

		if !confirm("Apply this change?") {
			fmt.Println("Aborted.")
			return nil
		}

		if err := zf.RemoveEnvVar(name); err != nil {
			return err
		}
		if err := zf.WriteTo(filePath); err != nil {
			return err
		}
		fmt.Printf("Removed: %s\n", name)
		return nil
	},
}

func init() {
	envAddCmd.Flags().BoolVar(&forceEnvAdd, "force", false, "ovewrite existing env variable")
	envCmd.AddCommand(envAddCmd)
	envCmd.AddCommand(envRemoveCmd)
	envCmd.AddCommand(envListCmd)
	rootCmd.AddCommand(envCmd)
}
