package cmd

import (
	"fmt"
	"strings"
)

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}

func showDiff(removed, added string) {
	fmt.Println("Changes:")
	for _, line := range strings.Split(removed, "\n") {
		if line != "" {
			fmt.Printf("- %s\n", line)
		}
	}
	for _, line := range strings.Split(added, "\n") {
		if line != "" {
			fmt.Printf("+ %s\n", line)
		}
	}
	fmt.Println()
}
