package cmd

import (
	"fmt"

	"github.com/Dan925/zshrc-manager/internal/backup"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage .zshrc backups",
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		backups, err := backup.List()
		if err != nil {
			return err
		}
		if len(backups) == 0 {
			fmt.Println("No backups found.")
			return nil
		}
		fmt.Printf("%-25s  %s\n", "TIMESTAMP", "SIZE")
		for _, b := range backups {
			fmt.Printf("%-25s  %d bytes\n", b.Timestamp.Format("2006-01-02 15:04:05"), b.SizeBytes)
		}
		return nil
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore <timestamp>",
	Short: "Restore a backup (timestamp format: YYYY-MM-DD_HH-MM-SS)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		backups, err := backup.List()
		if err != nil {
			return err
		}

		var match *backup.Backup
		for i, b := range backups {
			if b.Timestamp.Format("2006-01-02_15-04-05") == args[0] {
				match = &backups[i]
				break
			}
		}
		if match == nil {
			return fmt.Errorf("no backup found for %q — run 'zshrc backup list' to see available backups", args[0])
		}

		fmt.Printf("Would restore from: %s\n", match.Path)

		if dryRun {
			fmt.Println("Dry run — no changes made.")
			return nil
		}

		if !confirm(fmt.Sprintf("Restore from %s?", match.Timestamp.Format("2006-01-02 15:04:05"))) {
			fmt.Println("Aborted.")
			return nil
		}

		savedPath, err := backup.Create(filePath)
		if err != nil {
			return fmt.Errorf("backing up current state before restore: %w", err)
		}
		fmt.Printf("Current state backed up to: %s\n", savedPath)

		if err := backup.Restore(match.Path, filePath); err != nil {
			return err
		}
		fmt.Println("Restored successfully.")
		return nil
	},
}

func init() {
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	rootCmd.AddCommand(backupCmd)
}
