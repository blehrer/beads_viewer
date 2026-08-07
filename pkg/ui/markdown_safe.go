package ui

import "strings"

// ContainsANSI reports whether s includes terminal escape sequences.
// Markdown passed to glamour must not contain ANSI — it renders as visible garbage.
func ContainsANSI(s string) bool {
	// ponytail: duplicated from pkg/testutil — ui must not import testutil in prod.
	return strings.Contains(s, "\x1b") || strings.Contains(s, "\033")
}
