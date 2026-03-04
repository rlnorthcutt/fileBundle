package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Write creates the final bundle file at outputPath. It writes a TABLE OF CONTENTS
// section followed by the bundled file contents provided in bundleBody.
//
// The TABLE OF CONTENTS lists all relative file paths, one per line.
func Write(root string, outputPath string, filePaths []string, bundleBody string) error {
	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}

	// Ensure we are not overwriting one of the input files.
	for _, p := range filePaths {
		absInput, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		if absInput == absOutput {
			return fmt.Errorf("output file %s is one of the input files being bundled", outputPath)
		}
	}

	// Ensure output directory exists.
	if err := os.MkdirAll(filepath.Dir(absOutput), 0o755); err != nil {
		return err
	}

	f, err := os.Create(absOutput)
	if err != nil {
		return err
	}
	defer f.Close()

	// Build TABLE OF CONTENTS header.
	var tocBuilder strings.Builder
	tocBuilder.WriteString("TABLE OF CONTENTS\n")
	for _, p := range filePaths {
		rel, err := filepath.Rel(root, p)
		if err != nil {
			rel = p
		}
		tocBuilder.WriteString(filepath.ToSlash(rel))
		tocBuilder.WriteString("\n")
	}
	tocBuilder.WriteString("\n\n")

	if _, err := f.WriteString(tocBuilder.String()); err != nil {
		return err
	}
	if _, err := f.WriteString(bundleBody); err != nil {
		return err
	}

	return nil
}

