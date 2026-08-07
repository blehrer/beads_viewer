package testutil

import (
	"strings"
	"testing"
)

// ContainsANSI reports whether s includes terminal escape sequences.
// Markdown passed to glamour must not contain ANSI — it renders as visible garbage.
func ContainsANSI(s string) bool {
	return strings.Contains(s, "\x1b") || strings.Contains(s, "\033")
}

// AssertNoANSI fails if s is empty or contains ANSI escape sequences.
func AssertNoANSI(t *testing.T, name, s string) {
	t.Helper()
	if s == "" {
		t.Fatalf("%s: empty markdown", name)
	}
	if ContainsANSI(s) {
		idx := strings.Index(s, "\x1b")
		if idx < 0 {
			idx = strings.Index(s, "\033")
		}
		start := idx - 20
		if start < 0 {
			start = 0
		}
		end := idx + 40
		if end > len(s) {
			end = len(s)
		}
		t.Fatalf("%s contains ANSI escapes near: %q", name, s[start:end])
	}
}
