package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ParseIntention extracts process name from natural language intention
func ParseIntention(intention string) string {
	intention = strings.TrimSpace(intention)

	// If single word, return as-is
	if !strings.Contains(intention, " ") {
		return sanitizeProcessName(intention)
	}

	// Pattern: "I am working in X" or "I am working on X"
	re1 := regexp.MustCompile(`(?i)I am working (in|on) (.+)`)
	if matches := re1.FindStringSubmatch(intention); len(matches) > 2 {
		return sanitizeProcessName(matches[2])
	}

	// Pattern: "working in X" or "working on X"
	re2 := regexp.MustCompile(`(?i)working (in|on) (.+)`)
	if matches := re2.FindStringSubmatch(intention); len(matches) > 2 {
		return sanitizeProcessName(matches[2])
	}

	// Pattern: "I want to work on X"
	re3 := regexp.MustCompile(`(?i)I want to work on (.+)`)
	if matches := re3.FindStringSubmatch(intention); len(matches) > 1 {
		return sanitizeProcessName(matches[1])
	}

	// Pattern: "debugging X", "investigating X", etc.
	re4 := regexp.MustCompile(`(?i)(debugging|investigating|building|developing|fixing) (.+)`)
	if matches := re4.FindStringSubmatch(intention); len(matches) > 2 {
		return sanitizeProcessName(matches[2])
	}

	// Pattern: "X session"
	re5 := regexp.MustCompile(`(?i)(.+) session`)
	if matches := re5.FindStringSubmatch(intention); len(matches) > 1 {
		return sanitizeProcessName(matches[1])
	}

	// Fallback: take last word
	parts := strings.Fields(intention)
	if len(parts) > 0 {
		return sanitizeProcessName(parts[len(parts)-1])
	}

	return sanitizeProcessName(intention)
}

// sanitizeProcessName converts a string to a valid process name
func sanitizeProcessName(name string) string {
	// Replace non-alphanumeric characters (except - and _) with _
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	name = re.ReplaceAllString(name, "_")

	// Remove leading/trailing underscores
	name = strings.Trim(name, "_")

	if name == "" {
		return "unnamed"
	}

	return name
}

// FindDirectoryFuzzy finds a directory matching a wildcard pattern
func FindDirectoryFuzzy(pattern string) (string, error) {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Convert pattern to case-insensitive
	pattern = strings.ToLower(pattern)

	var matched string
	maxDepth := 3

	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if !info.IsDir() {
			return nil
		}

		// Check depth
		relPath, _ := filepath.Rel(currentDir, path)
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > maxDepth {
			return filepath.SkipDir
		}

		// Match pattern (case-insensitive)
		dirName := strings.ToLower(info.Name())
		patternClean := strings.Trim(pattern, "*")

		if strings.Contains(dirName, patternClean) {
			matched = path
			return filepath.SkipAll // Found, stop searching
		}

		return nil
	})

	return matched, err
}

// FindDirectoriesInteractive searches for directories matching a pattern
func FindDirectoriesInteractive(pattern string) ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return FindDirectoriesFrom(pattern, homeDir, 5, 20)
}

// FindDirectoriesFrom searches for directories from a starting point
// Much faster than searching from home when you're already in a subdirectory
func FindDirectoriesFrom(pattern, startDir string, maxDepth, maxResults int) ([]string, error) {
	var results []string
	pattern = strings.ToLower(pattern)

	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if !info.IsDir() {
			return nil
		}

		// Check depth
		relPath, _ := filepath.Rel(startDir, path)
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > maxDepth {
			return filepath.SkipDir
		}

		// Match pattern (case-insensitive)
		dirName := strings.ToLower(info.Name())
		if strings.Contains(dirName, pattern) {
			results = append(results, path)

			// Stop after max results
			if len(results) >= maxResults {
				return filepath.SkipAll
			}
		}

		return nil
	})

	return results, err
}

// ExpandPath expands ~ to home directory and handles relative paths
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}

	// Make absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// FormatPath formats a path for display (shortens home directory)
func FormatPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}

	return path
}

// PrintError prints an error message in a consistent format
func PrintError(format string, args ...interface{}) {
	fmt.Printf("✗ "+format+"\n", args...)
}

// PrintSuccess prints a success message in a consistent format
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("✓ "+format+"\n", args...)
}