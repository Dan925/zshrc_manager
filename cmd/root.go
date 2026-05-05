package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var filePath string
var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "zshrc",
	Short: "Manage your .zshrc file safely",
	Long:  "A CLI tool for listing, searching, and editing your .zshrc with automatic backups.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	rootCmd.PersistentFlags().StringVar(&filePath, "file", home+"/.zshrc", "path to zshrc file")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would happen without making changes")
}
