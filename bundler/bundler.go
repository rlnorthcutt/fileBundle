package bundler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// Bundle reads all files in filePaths and returns a single string containing
// their contents formatted according to the bundle format specification.
//
// Each file block will look like:
//
// -----------------------------------------------
// relative/path/to/file.ext
// -----------------------------------------------
//
// <file contents>
//
func Bundle(root string, filePaths []string) (string, error) {
	bar := progressbar.Default(int64(len(filePaths)), "Bundling files")

	var result string

	for _, fullPath := range filePaths {
		if err := bar.Add(1); err != nil {
			// Progress failures should not break the core logic.
		}

		relPath, err := filepath.Rel(root, fullPath)
		if err != nil {
			relPath = fullPath
		}

		content, err := readFileSafely(fullPath)
		if err != nil {
			// Gracefully skip unreadable files; they simply won't appear in the bundle body.
			continue
		}

		block := fmt.Sprintf(
			"-----------------------------------------------\n%s\n-----------------------------------------------\n\n%s\n\n",
			filepath.ToSlash(relPath),
			content,
		)
		result += block
	}

	return result, nil
}

func readFileSafely(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

