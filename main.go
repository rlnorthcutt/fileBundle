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

// Config holds all user-configurable options for a single run.
type Config struct {
	InputDir       string
	IncludeDirs    string
	Extensions     string
	Exclude        string
	OutputFilename string
	NonInteractive bool
}

var cfg Config
var interactive bool

var rootCmd = &cobra.Command{
	Use:   "filebundle",
	Short: "Flatten a directory structure into a single TOC-indexed file for AI context.",
	Run:   executeBundle,
}

func main() {
	// Detect whether stdin is attached to a TTY so we can decide
	// if interactive prompts are appropriate.
	stat, _ := os.Stdin.Stat()
	interactive = (stat.Mode() & os.ModeCharDevice) != 0

	if err := rootCmd.Execute(); err != nil {
		handleError("executing command", err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.InputDir, "input", "i", ".", "Root directory to crawl")
	rootCmd.Flags().StringVarP(&cfg.IncludeDirs, "include", "d", "*", "Subdirectories to include")
	rootCmd.Flags().StringVarP(&cfg.Extensions, "extensions", "e", ".md,.txt", "File extensions to include")
	rootCmd.Flags().StringVarP(&cfg.Exclude, "exclude", "x", ".git,node_modules,bin", "Patterns or folders to exclude")
	rootCmd.Flags().StringVarP(&cfg.OutputFilename, "output", "o", "bundle.txt", "Output filename")

	rootCmd.Flags().BoolVar(&cfg.NonInteractive, "non-interactive", false, "Disable prompts and use defaults")
}

func executeBundle(cmd *cobra.Command, args []string) {
	resolveConfig(cmd)

	// Validate input directory
	if info, err := os.Stat(cfg.InputDir); os.IsNotExist(err) || !info.IsDir() {
		handleError("validating input directory", fmt.Errorf("path does not exist or is not a directory: %s", cfg.InputDir))
	}

	// Summary
	fmt.Printf("\nPreparing to bundle files with the following settings:\n")
	fmt.Printf("  Input Root:  %s\n", cfg.InputDir)
	fmt.Printf("  Include:     %s\n", cfg.IncludeDirs)
	fmt.Printf("  Extensions:  %s\n", cfg.Extensions)
	fmt.Printf("  Exclude:     %s\n", cfg.Exclude)
	fmt.Printf("  Output File: %s\n", cfg.OutputFilename)

	if shouldPrompt() {
		confirmation := promptUser("\nDo you want to proceed? (y/n): ", "y")
		if strings.ToLower(confirmation) != "y" {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	fmt.Println()

	// Crawl directory
	filePaths, err := crawler.Crawl(cfg.InputDir, cfg.IncludeDirs, cfg.Extensions, cfg.Exclude)
	handleError("crawling directory", err)

	if len(filePaths) == 0 {
		fmt.Println("No matching files found.")
		return
	}

	// Bundle files
	rootAbs, err := filepath.Abs(cfg.InputDir)
	handleError("resolving input directory", err)

	bundleData, err := bundler.Bundle(rootAbs, filePaths)
	handleError("bundling file contents", err)

	// Write output
	err = writer.Write(rootAbs, cfg.OutputFilename, filePaths, bundleData)
	handleError("writing bundle to disk", err)

	fmt.Printf("\nSuccessfully bundled %d files into %s\n", len(filePaths), cfg.OutputFilename)
}

func resolveConfig(cmd *cobra.Command) {
	// Configuration precedence:
	// 1. CLI flags (highest)
	// 2. Environment variables
	// 3. Interactive prompts (if allowed)
	// 4. Built-in defaults (lowest)
	promptIfMissing(cmd, "input", &cfg.InputDir, "FILEBUNDLE_INPUT", "Enter root directory")
	promptIfMissing(cmd, "include", &cfg.IncludeDirs, "FILEBUNDLE_INCLUDE", "Enter directories to include")
	promptIfMissing(cmd, "extensions", &cfg.Extensions, "FILEBUNDLE_EXTENSIONS", "Enter file extensions")
	promptIfMissing(cmd, "exclude", &cfg.Exclude, "FILEBUNDLE_EXCLUDE", "Enter exclude patterns")
	promptIfMissing(cmd, "output", &cfg.OutputFilename, "FILEBUNDLE_OUTPUT", "Enter output filename")
}

func promptIfMissing(cmd *cobra.Command, flag string, value *string, envKey string, message string) {
	// CLI flag wins.
	if cmd.Flags().Changed(flag) {
		return
	}

	// Environment variable override.
	if env := os.Getenv(envKey); env != "" {
		*value = env
		return
	}

	// Prompt only if allowed.
	if !shouldPrompt() {
		return
	}

	*value = promptUser(
		fmt.Sprintf("%s (default: %s): ", message, *value),
		*value,
	)
}

func shouldPrompt() bool {
	return interactive && !cfg.NonInteractive
}

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

/*
ERROR HANDLING
*/

func handleError(step string, err error) {
	if err != nil {
		log.Fatalf("Error %s: %v", step, err)
	}
}
