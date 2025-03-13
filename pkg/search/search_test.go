package search

import (
	"os"
	"path/filepath"
	"testing"
)

// createTestFiles creates a temporary directory with realistic test files
func createTestFiles(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "searchbot_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a realistic file structure
	testFiles := map[string]struct {
		path    string
		content string
		mode    os.FileMode
	}{
		// Documents
		"Documents/report-2024.pdf":          {content: "pdf content"},
		"Documents/meeting-notes.txt":        {content: "meeting notes"},
		"Documents/presentation slides.pptx": {content: "presentation"},
		"Documents/budget 2024.xlsx":         {content: "budget data"},

		// Source code files
		"Projects/go/main.go":               {content: "package main"},
		"Projects/go/test.go":               {content: "package test"},
		"Projects/python/script.py":         {content: "python code"},
		"Projects/node/package.json":        {content: "{}"},
		"Projects/node/node_modules/lib.js": {content: "library"},

		// Media files
		"Media/vacation2024.jpg": {content: "image"},
		"Media/profile-pic.png":  {content: "profile"},
		"Media/video-2024.mp4":   {content: "video"},

		// Hidden files and directories
		".git/config":           {content: "git config"},
		".vscode/settings.json": {content: "settings"},

		// Files with special characters
		"Documents/résumé.pdf":                                           {content: "resume"},
		"Documents/документ.txt":                                         {content: "document"},
		"Projects/test-file-with-very-very-very-very-very-long-name.txt": {content: "test"},

		// Files with different permissions
		"Documents/readonly.txt":  {content: "readonly", mode: 0444},
		"Documents/executable.sh": {content: "#!/bin/bash", mode: 0755},

		// Empty directories
		"Empty Directory/": {mode: os.ModeDir | 0755},

		// Symlinks (if supported)
		"Documents/link-to-report.pdf": {content: "report-2024.pdf", mode: os.ModeSymlink},
	}

	// Create the files
	for path, info := range testFiles {
		fullPath := filepath.Join(tempDir, path)

		// Create parent directories
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		// Handle different file types
		if info.mode&os.ModeDir != 0 {
			err = os.MkdirAll(fullPath, info.mode)
		} else if info.mode&os.ModeSymlink != 0 {
			err = os.Symlink(filepath.Join(tempDir, info.content), fullPath)
		} else {
			// Regular file
			mode := info.mode
			if mode == 0 {
				mode = 0644
			}
			err = os.WriteFile(fullPath, []byte(info.content), mode)
		}

		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestSearchFiles(t *testing.T) {
	tempDir, cleanup := createTestFiles(t)
	defer cleanup()

	tests := []struct {
		name             string
		pattern          string
		opts             SearchOptions
		expectedCount    int
		expectedFiles    []string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:    "Search PDF files",
			pattern: ".pdf",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: true,
			},
			shouldContain:    []string{"report-2024.pdf", "résumé.pdf"},
			shouldNotContain: []string{"document.txt", "script.py"},
		},
		{
			name:    "Search by year 2024",
			pattern: "2024",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: true,
			},
			shouldContain: []string{"report-2024.pdf", "vacation2024.jpg", "budget 2024.xlsx"},
		},
		{
			name:    "Search with spaces in filename",
			pattern: "meeting",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: false,
			},
			shouldContain: []string{"meeting-notes.txt"},
		},
		{
			name:    "Search source code files",
			pattern: ".go",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: true,
			},
			shouldContain:    []string{"main.go", "test.go"},
			shouldNotContain: []string{"script.py", "package.json"},
		},
		{
			name:    "Search with Unicode characters",
			pattern: "документ",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: true,
			},
			shouldContain: []string{"документ.txt"},
		},
		{
			name:    "Search executable files",
			pattern: ".sh",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: true,
			},
			shouldContain: []string{"executable.sh"},
		},
		{
			name:    "Search in specific directory",
			pattern: "script",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: false,
			},
			shouldContain: []string{"script.py"},
		},
		{
			name:    "Search media files",
			pattern: "vacation",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: false,
			},
			shouldContain: []string{"vacation2024.jpg"},
		},
		{
			name:    "Skip hidden directories",
			pattern: "config",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: false,
			},
			shouldNotContain: []string{".git/config"},
		},
		{
			name:    "Skip node_modules",
			pattern: "lib.js",
			opts: SearchOptions{
				Recursive:     true,
				ExactMatch:    false,
				CaseSensitive: false,
			},
			shouldNotContain: []string{"node_modules/lib.js"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := SearchFiles(tt.pattern, tempDir, tt.opts)
			if err != nil {
				t.Errorf("SearchFiles() error = %v", err)
				return
			}

			// Check for required files
			foundFiles := make(map[string]bool)
			for _, result := range results {
				foundFiles[filepath.Base(result.Path)] = true
			}

			// Check files that should be found
			for _, expected := range tt.shouldContain {
				if !foundFiles[expected] {
					t.Errorf("SearchFiles() should contain %s, but didn't", expected)
				}
			}

			// Check files that should not be found
			for _, notExpected := range tt.shouldNotContain {
				if foundFiles[notExpected] {
					t.Errorf("SearchFiles() should not contain %s, but did", notExpected)
				}
			}
		})
	}
}

func TestSearchFilesErrors(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		root    string
		opts    SearchOptions
	}{
		{
			name:    "Non-existent directory",
			pattern: "test",
			root:    "/path/that/does/not/exist",
			opts:    SearchOptions{Recursive: true},
		},
		{
			name:    "Empty pattern",
			pattern: "",
			root:    ".",
			opts:    SearchOptions{Recursive: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SearchFiles(tt.pattern, tt.root, tt.opts)
			if err == nil {
				t.Errorf("SearchFiles() expected error for %s", tt.name)
			}
		})
	}
}
