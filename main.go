package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"filebundle/bundler"
	"filebundle/crawler"
	"filebundle/writer"
)

var (
	inputDir       string
	includeSubdirs string
	extensions     string
	exclude        string
	outputFilename string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		handleError("executing command", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "filebundle",
	Short: "Flatten a directory structure into a single TOC-indexed file for AI context.",
	Run:   executeBundle,
}

func init() {
	rootCmd.Flags().StringVarP(&inputDir, "input", "i", ".", "Root directory to crawl")
	rootCmd.Flags().StringVarP(&includeSubdirs, "include", "d", "*", "Subdirectories to include (comma-separated or *)")
	rootCmd.Flags().StringVarP(&extensions, "extensions", "e", ".md,.txt", "File extensions to bundle (comma-separated)")
	rootCmd.Flags().StringVarP(&exclude, "exclude", "x", ".git,node_modules,bin", "Patterns or folders to exclude")
	rootCmd.Flags().StringVarP(&outputFilename, "output", "o", "bundle.txt", "Output filename")
}

// executeBundle orchestrates the crawl, bundle, and write process.
func executeBundle(cmd *cobra.Command, args []string) {
	// 1. Prompt for missing or default values
	inputDir = promptUser("Enter the root directory to crawl (default: '.'): ", inputDir)
	includeSubdirs = promptUser("Enter subdirectories to include (comma-separated, '*' for all): ", includeSubdirs)
	extensions = promptUser("Enter file extensions to include (e.g., 'md,txt'): ", extensions)
	exclude = promptUser("Enter patterns to exclude (e.g., '.git,node_modules'): ", exclude)
	outputFilename = promptUser("Enter the output filename (default: 'bundle.txt'): ", outputFilename)

	// 2. Validate Input
	if info, err := os.Stat(inputDir); os.IsNotExist(err) || !info.IsDir() {
		handleError("validating input directory", fmt.Errorf("path does not exist or is not a directory: %s", inputDir))
	}

	// 3. Summary and Confirmation
	fmt.Printf("\nPreparing to bundle files with the following settings:\n")
	fmt.Printf("  Input Root:  %s\n", inputDir)
	fmt.Printf("  Include:     %s\n", includeSubdirs)
	fmt.Printf("  Extensions:  %s\n", extensions)
	fmt.Printf("  Exclude:     %s\n", exclude)
	fmt.Printf("  Output File: %s\n", outputFilename)

	confirmation := promptUser("\nDo you want to proceed? (y/n): ", "y")
	if strings.ToLower(confirmation) != "y" {
		fmt.Println("Operation cancelled.")
		return
	}
	fmt.Println()

	// Step 1: Crawl the directory to find matching files
	// Passing strings to crawler to handle the filtering logic internally
	filePaths, err := crawler.Crawl(inputDir, includeSubdirs, extensions, exclude)
	handleError("crawling directory", err)

	if len(filePaths) == 0 {
		fmt.Println("No matching files found. Check your filters and try again.")
		return
	}

	// Step 2: Generate the bundled file contents
	// bundler.Bundle reads the file contents and formats them with headers/dividers
	rootAbs, err := filepath.Abs(inputDir)
	handleError("resolving input directory", err)

	bundleData, err := bundler.Bundle(rootAbs, filePaths)
	handleError("bundling file contents", err)

	// Step 3: Write the final product to disk
	// writer.Write handles the TOC header + the bundleData
	err = writer.Write(rootAbs, outputFilename, filePaths, bundleData)
	handleError("writing bundle to disk", err)

	fmt.Printf("\nSuccessfully bundled %d files into %s\n", len(filePaths), outputFilename)
}

// promptUser displays a message and returns user input, falling back to a default value.
func promptUser(message string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

// handleError provides a centralized way to exit on non-recoverable errors.
func handleError(step string, err error) {
	if err != nil {
		log.Fatalf("Error %s: %v", step, err)
	}
}
