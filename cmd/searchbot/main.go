package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"searchbot/pkg/display"
	"searchbot/pkg/search"
)

func main() {
	// Define flags with new defaults
	recursive := flag.Bool("nr", false, "Non-recursive search (by default search is recursive)")
	exactMatch := flag.Bool("e", false, "Match exact filename (by default matches substrings)")
	caseSensitive := flag.Bool("i", false, "Case insensitive search (by default search is case sensitive)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: s [options] <search_pattern> [directory]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -nr   Non-recursive search (by default search is recursive)\n")
		fmt.Fprintf(os.Stderr, "  -e    Match exact filename (by default matches substrings)\n")
		fmt.Fprintf(os.Stderr, "  -i    Case insensitive search (by default search is case sensitive)\n")
		fmt.Fprintf(os.Stderr, "\nIf directory is not specified, searches in current directory\n")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	pattern := flag.Arg(0)

	// Use current directory by default, or user-specified directory
	root := "."
	if flag.NArg() > 1 {
		root = flag.Arg(1)
	}

	// Convert to absolute path
	root, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Verify directory exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Directory does not exist: %s\n", root)
		os.Exit(1)
	}

	// Create search options with new defaults
	opts := search.SearchOptions{
		Recursive:     !*recursive,     // Default true, -nr flag makes it false
		ExactMatch:    *exactMatch,     // Default false
		CaseSensitive: !*caseSensitive, // Default true, -i flag makes it false
	}

	// Start the search with a status message
	fmt.Printf("Searching for '%s' in %s...\n", pattern, root)

	results, err := search.SearchFiles(pattern, root, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching files: %v\n", err)
		os.Exit(1)
	}

	display.PrintResults(results)
}
