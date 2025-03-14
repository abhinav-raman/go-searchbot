package search

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// SearchOptions contains search configuration
type SearchOptions struct {
	Recursive     bool
	ExactMatch    bool
	CaseSensitive bool
}

// SearchResult represents a single file search result
type SearchResult struct {
	Path    string
	Name    string
	Size    int64
	ModTime string
}

var (
	ErrEmptyPattern = errors.New("search pattern cannot be empty")
	ErrInvalidPath  = errors.New("invalid search path")
)

// shouldSkipDirectory returns true if the directory should be skipped
func shouldSkipDirectory(path string) bool {
	// Skip hidden directories and common system paths
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".") || // Hidden directories
		base == "node_modules" || // Skip node_modules
		base == "Library" || // Skip macOS Library
		base == "System" || // Skip System directory
		base == "Applications" // Skip Applications
}

// SearchFiles searches for files containing the given pattern in their names
func SearchFiles(pattern string, root string, opts SearchOptions) ([]SearchResult, error) {
	// Validate inputs
	if pattern == "" {
		return nil, ErrEmptyPattern
	}

	// Check if root exists
	if _, err := os.Stat(root); err != nil {
		return nil, ErrInvalidPath
	}

	var results []SearchResult

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Handle directories
		if info.IsDir() {
			if path != root && !opts.Recursive {
				return filepath.SkipDir
			}
			if shouldSkipDirectory(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Prepare strings for comparison
		fileName := info.Name()
		searchPattern := pattern

		if !opts.CaseSensitive {
			fileName = strings.ToLower(fileName)
			searchPattern = strings.ToLower(searchPattern)
		}

		// Check if file matches
		var matches bool
		if opts.ExactMatch {
			matches = fileName == searchPattern
		} else {
			matches = strings.Contains(fileName, searchPattern)
		}

		if matches {
			results = append(results, SearchResult{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
