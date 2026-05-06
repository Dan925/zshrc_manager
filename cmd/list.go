package cmd

import (
	"fmt"

	"github.com/Dan925/zshrc-manager/internal/parser"
	"github.com/spf13/cobra"
)

var verbose bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases or functions in your .zshrc",
}

var listAliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "List all aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		zf, err := parser.NewParser(filePath).Parse()
		if err != nil {
			return err
		}
		if len(zf.Aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}
		for _, a := range zf.Aliases {
			fmt.Printf("%s -> %s\n", a.Name, a.Value)
		}
		return nil
	},
}

var listFunctionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "List all functions",
	RunE: func(cmd *cobra.Command, args []string) error {
		zf, err := parser.NewParser(filePath).Parse()
		if err != nil {
			return err
		}
		if len(zf.Functions) == 0 {
			fmt.Println("No functions found.")
			return nil
		}
		for _, f := range zf.Functions {
			if verbose {
				fmt.Printf("# %s (lines %d-%d)\n%s\n", f.Name, f.StartLine, f.EndLine, f.Body)
			} else {
				fmt.Printf("%s (line %d)\n", f.Name, f.StartLine)
			}
		}
		return nil
	},
}

func init() {
	listFunctionsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show full function body")
	listCmd.AddCommand(listAliasesCmd)
	listCmd.AddCommand(listFunctionsCmd)
	rootCmd.AddCommand(listCmd)
}
