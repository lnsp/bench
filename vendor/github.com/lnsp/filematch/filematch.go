// Package filematch provides file path matching using gitignore-style rules.
package filematch

import (
	"path/filepath"
	"strings"
)

// MatchPath tests if a path matches one of the given patterns.
func MatchPath(path string, ignorePatterns []string) (bool, error) {
	for _, pattern := range ignorePatterns {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if (strings.HasSuffix(pattern, "/") && strings.HasPrefix(path, pattern)) || path == pattern || matched {
			return true, nil
		}
	}
	return false, nil
}
