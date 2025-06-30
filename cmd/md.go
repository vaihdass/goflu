package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaihdass/goflu/internal/parser"
)

var (
	outputFile string
	overwrite  bool
)

var mdCmd = &cobra.Command{
	Use:   "md <file>",
	Short: "Convert Confluence HTML file to Markdown",
	Long: `Parse a Confluence HTML export file and convert it to clean Markdown format.
The command extracts the main content while filtering out navigation elements
and other UI components.`,
	Args: cobra.ExactArgs(1),
	RunE: runMd,
}

func init() {
	rootCmd.AddCommand(mdCmd)
	mdCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: input file with .md extension)")
	mdCmd.Flags().BoolVarP(&overwrite, "force", "f", false, "Overwrite output file if it exists")
}

func runMd(_ *cobra.Command, args []string) error {
	inputFile := args[0]

	// Resolve input file path
	absInputFile, err := filepath.Abs(inputFile)
	if err != nil {
		return fmt.Errorf("failed to resolve input file path: %w", err)
	}

	// Check if input file exists
	if _, err := os.Stat(absInputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", absInputFile)
	}

	// Determine output file path
	if outputFile == "" {
		outputFile = strings.TrimSuffix(absInputFile, filepath.Ext(absInputFile)) + ".md"
	} else {
		outputFile, err = filepath.Abs(outputFile)
		if err != nil {
			return fmt.Errorf("failed to resolve output file path: %w", err)
		}
	}

	// Check if output file exists and overwrite flag
	if _, err := os.Stat(outputFile); err == nil && !overwrite {
		return fmt.Errorf("output file already exists: %s (use -f to overwrite)", outputFile)
	}

	// Read input file
	htmlContent, err := os.ReadFile(absInputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse HTML to Markdown
	fmt.Printf("Parsing %s...\n", absInputFile)
	markdown, err := parser.ParseConfluenceHTML(string(htmlContent))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Write output file
	err = os.WriteFile(outputFile, []byte(markdown), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully converted to %s\n", outputFile)
	return nil
}
