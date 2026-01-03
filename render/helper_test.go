package render

import "strings"

// contains checks if a string contains a substring (case-sensitive)
// This is a helper function used across all renderer tests
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
