package crawler

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Crawl walks the inputDir and returns a list of files that match the
// provided includeSubdirs, extensions, and exclude patterns.
//
// - includeSubdirs: "*" means all, otherwise comma-separated list of top-level subdirs under inputDir.
// - extensions: comma-separated list with or without leading dots (e.g. "md,txt" or ".md,.txt").
// - exclude: comma-separated list of substrings; any path containing one of them is skipped.
func Crawl(inputDir, includeSubdirs, extensions, exclude string) ([]string, error) {
	if inputDir == "" {
		return nil, errors.New("input directory is required")
	}

	normalizedRoot, err := filepath.Abs(inputDir)
	if err != nil {
		return nil, err
	}

	includeSet := parseList(includeSubdirs)
	extSet := parseExtensions(extensions)
	excludeList := parseSlice(exclude)

	var files []string

	err = filepath.WalkDir(normalizedRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			// If we can't read a directory entry, skip it but continue walking.
			return nil
		}

		// Always normalize path to use forward slashes for comparisons.
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(normalizedRoot, absPath)
		if err != nil {
			return nil
		}
		relPath = filepath.ToSlash(relPath)

		// Skip the root itself.
		if relPath == "." {
			return nil
		}

		// Exclude patterns: if any substring matches, skip.
		for _, pat := range excludeList {
			if pat != "" && strings.Contains(relPath, pat) {
				if d.IsDir() {
					// Skip entire directory subtree.
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Include filtering is applied to the first path segment under root.
		if len(includeSet) > 0 && !includeSet["*"] {
			firstSegment := relPath
			if idx := strings.Index(firstSegment, "/"); idx != -1 {
				firstSegment = firstSegment[:idx]
			}

			if d.IsDir() {
				// If this is a top-level directory and not in includeSet, skip its subtree.
				parent := filepath.ToSlash(filepath.Dir(relPath))
				if parent == "." && !includeSet[firstSegment] {
					return filepath.SkipDir
				}
			} else {
				// For files, only accept if their top-level directory is included.
				if parent := firstSegment; !includeSet[parent] && parent != "." {
					return nil
				}
			}
		}

		// Only collect files.
		if d.IsDir() {
			return nil
		}

		// Extension filtering.
		if len(extSet) > 0 {
			ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(relPath), "."))
			if ext == "" || !extSet[ext] {
				return nil
			}
		}

		files = append(files, absPath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func parseList(s string) map[string]bool {
	result := make(map[string]bool)
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			result[p] = true
		}
	}
	return result
}

func parseSlice(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func parseExtensions(s string) map[string]bool {
	result := make(map[string]bool)
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		p = strings.TrimPrefix(p, ".")
		p = strings.ToLower(p)
		result[p] = true
	}
	return result
}

