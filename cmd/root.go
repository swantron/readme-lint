package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/swantron/readme-lint/pkg/linter"
)

var rootCmd = &cobra.Command{
	Use:   "readme-lint [file]",
	Short: "A linter for README.md files",
	Long:  `readme-lint is a fast, standalone command-line tool to enforce quality and completeness standards for README.md files.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := "./README.md"
		if len(args) > 0 {
			filePath = args[0]
		}

		l := linter.NewLinter()
		results, err := l.Run(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(results) > 0 {
			for _, result := range results {
				if result.Line > 0 {
					fmt.Printf("[FAIL] %s:%d: %s\n", filePath, result.Line, result.Message)
				} else {
					fmt.Printf("[FAIL] %s: %s\n", filePath, result.Message)
				}
			}
			os.Exit(1)
		} else {
			fmt.Println("✓ All checks passed!")
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
