package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check .zshrc for syntax errors using zsh -n",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("zsh", "-n", filePath).CombinedOutput()
		if err != nil {
			fmt.Printf("Syntax errors found:\n%s\n", string(out))
			return fmt.Errorf("validation failed")
		}
		fmt.Printf("%s: syntax OK\n", filePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
