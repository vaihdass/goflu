package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goflu",
	Short: "A CLI tool to parse Confluence HTML files to Markdown",
	Long: `goflu is a command-line tool that converts Confluence HTML export files
to clean Markdown format, extracting the main content while filtering out
navigation elements and other UI components.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
