package testutil

import "testing"

func TestContainsANSI(t *testing.T) {
	if ContainsANSI("plain text") {
		t.Fatal("plain text should not contain ANSI")
	}
	if !ContainsANSI("\x1b[38;2;255;0;0m●\x1b[0m") {
		t.Fatal("expected ESC byte form to be detected")
	}
	if !ContainsANSI("\033[31mred\033[0m") {
		t.Fatal("expected octal ESC form to be detected")
	}
}
