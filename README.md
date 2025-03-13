# SearchBot

A fast and simple terminal-based file search engine written in Go.

## Features

- Recursive file search starting from current directory
- Case-sensitive search by default
- Matches filenames containing the search pattern as substring
- Beautiful terminal output with colors
- Display file size in human-readable format
- Show file modification time

## Installation

1. Clone the repository
```bash
git clone <repository-url>
cd searchbot
```

2. Build the project:
```bash
go build -o s cmd/searchbot/main.go
```

3. Make it globally available (choose one):
```bash
# Option 1 - User specific (recommended)
mkdir -p ~/bin
mv s ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc for bash
source ~/.zshrc  # or source ~/.bashrc for bash

# Option 2 - System wide
sudo mv s /usr/local/bin/
```

## Usage

Basic usage:
```bash
s <search_pattern> [directory]
```

If directory is not specified, searches in the current directory.

### Examples
```bash
s document           # Search for files containing "document" in current directory
s report ~/Documents # Search in specific directory
s .pdf              # Search for PDF files
```

### Options

- `-nr`: Non-recursive search (by default search is recursive)
- `-i`: Case insensitive search (by default search is case sensitive)
- `-e`: Match exact filename (by default matches substrings)

### Examples with Options
```bash
s -nr report        # Non-recursive search
s -i Document       # Case insensitive search
s -e report.pdf     # Match exact filename
s -i -e Report.txt  # Case insensitive and exact match
```

## Output Format

The search results are displayed in a table format with the following columns:
- Name: File name (truncated if too long)
- Size: File size in human-readable format (B, KB, MB, GB)
- Modified: Last modification date and time
- Path: Full path to the file

## Performance

- Automatically skips hidden files and directories
- Skips system directories and common large directories (node_modules, Library, etc.)
- Efficient file traversal using Go's filepath.Walk

## Notes

- The search starts from the current directory by default
- Hidden files and directories (starting with '.') are skipped
- System directories and node_modules are excluded from search 