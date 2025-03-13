package display

import (
	"bytes"
	"io"
	"os"
	"searchbot/pkg/search"
	"strings"
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		// Basic cases
		{name: "Zero bytes", size: 0, expected: "0 B"},
		{name: "One byte", size: 1, expected: "1 B"},
		{name: "999 bytes", size: 999, expected: "999 B"},

		// Kilobytes
		{name: "1 KB", size: 1024, expected: "1.0 KB"},
		{name: "1.5 KB", size: 1536, expected: "1.5 KB"},
		{name: "999 KB", size: 1024 * 999, expected: "999.0 KB"},

		// Megabytes
		{name: "1 MB", size: 1024 * 1024, expected: "1.0 MB"},
		{name: "1.5 MB", size: 1024 * 1024 * 3 / 2, expected: "1.5 MB"},
		{name: "999 MB", size: 1024 * 1024 * 999, expected: "999.0 MB"},

		// Gigabytes
		{name: "1 GB", size: 1024 * 1024 * 1024, expected: "1.0 GB"},
		{name: "1.5 GB", size: 1024 * 1024 * 1024 * 3 / 2, expected: "1.5 GB"},

		// Edge cases
		{name: "Max int64", size: 9223372036854775807, expected: "8.0 EB"},
		{name: "Negative size", size: -1024, expected: "-1.0 KB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSize(tt.size)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %v, want %v", tt.size, result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		maxLen   int
		expected string
	}{
		// Basic cases
		{name: "No truncation needed", str: "short.txt", maxLen: 20, expected: "short.txt"},
		{name: "Exact length", str: "exactly.txt", maxLen: 11, expected: "exactly.txt"},
		{name: "Need truncation", str: "very_long_filename.txt", maxLen: 10, expected: "very_lo..."},

		// Edge cases
		{name: "Empty string", str: "", maxLen: 5, expected: ""},
		{name: "MaxLen = 0", str: "test.txt", maxLen: 0, expected: ""},
		{name: "MaxLen = 3", str: "test.txt", maxLen: 3, expected: "..."},
		{name: "MaxLen = 4", str: "test.txt", maxLen: 4, expected: "t..."},

		// Unicode strings
		{name: "Unicode string", str: "файл.txt", maxLen: 5, expected: "фа..."},
		{name: "Mixed Unicode", str: "file文件.txt", maxLen: 7, expected: "file..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.str, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.str, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestPrintResults(t *testing.T) {
	// Save the original stdout
	oldStdout := os.Stdout

	testCases := []struct {
		name     string
		results  []search.SearchResult
		contains []string
		excludes []string
	}{
		{
			name:    "Empty results",
			results: []search.SearchResult{},
			contains: []string{
				"No files found",
			},
		},
		{
			name: "Single result",
			results: []search.SearchResult{
				{
					Path:    "/path/to/file.txt",
					Name:    "file.txt",
					Size:    1024,
					ModTime: "2024-03-20 10:00:00",
				},
			},
			contains: []string{
				"Found 1 files (Total size: 1.0 KB)",
				"NAME",
				"SIZE",
				"MODIFIED",
				"PATH",
				"file.txt",
				"1.0 KB",
				"2024-03-20 10:00:00",
				"/path/to/file.txt",
			},
		},
		{
			name: "Multiple results",
			results: []search.SearchResult{
				{
					Path:    "/path/to/file1.txt",
					Name:    "file1.txt",
					Size:    1024,
					ModTime: "2024-03-20 10:00:00",
				},
				{
					Path:    "/path/to/file2.txt",
					Name:    "file2.txt",
					Size:    2048,
					ModTime: "2024-03-20 11:00:00",
				},
			},
			contains: []string{
				"Found 2 files (Total size: 3.0 KB)",
				"NAME",
				"SIZE",
				"MODIFIED",
				"PATH",
				"file1.txt",
				"file2.txt",
				"1.0 KB",
				"2.0 KB",
				"/path/to/file1.txt",
				"/path/to/file2.txt",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a pipe to capture stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function
			PrintResults(tc.results)

			// Close the write end of the pipe
			w.Close()

			// Read the output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Check if all expected strings are in the output
			for _, expected := range tc.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Output should contain %q but didn't.\nFull output:\n%s", expected, output)
				}
			}

			// Check if excluded strings are not in the output
			for _, excluded := range tc.excludes {
				if strings.Contains(output, excluded) {
					t.Errorf("Output should not contain %q but did.\nFull output:\n%s", excluded, output)
				}
			}
		})
	}

	// Restore the original stdout
	os.Stdout = oldStdout
}
