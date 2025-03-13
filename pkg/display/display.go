package display

import (
	"fmt"
	"searchbot/pkg/search"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// FormatSize formats file size in human-readable format
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// PrintResults displays the search results in a formatted way
func PrintResults(results []search.SearchResult) {
	if len(results) == 0 {
		color.Yellow("\nNo files found")
		return
	}

	// Sort results by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	// Calculate total size
	var totalSize int64
	for _, result := range results {
		totalSize += result.Size
	}

	color.Green("\nFound %d files (Total size: %s)\n", len(results), FormatSize(totalSize))
	fmt.Println(strings.Repeat("-", 100))

	// Print header
	color.Blue("%-50s %-20s %-15s %s\n", "NAME", "SIZE", "MODIFIED", "PATH")
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		fmt.Printf("%-50s %-20s %-15s %s\n",
			truncateString(result.Name, 47),
			FormatSize(result.Size),
			result.ModTime,
			result.Path,
		)
	}
	fmt.Println(strings.Repeat("-", 100))
}

// truncateString truncates a string if it's longer than maxLen
func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}
