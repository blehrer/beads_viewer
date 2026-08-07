package beadscli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectFromBeadsDir_BDWorkspace(t *testing.T) {
	t.Setenv(envTool, "")
	t.Cleanup(func() { SetTool("br") })

	dir := t.TempDir()
	beadsDir := filepath.Join(dir, ".beads")
	if err := os.MkdirAll(beadsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(beadsDir, "metadata.json"), []byte(`{"backend":"dolt"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	DetectFromBeadsDir(beadsDir)
	if Tool() != "bd" {
		t.Fatalf("Tool() = %q, want bd", Tool())
	}
	if ExportFlushCommand() != "bd export --no-memories -o .beads/issues.jsonl" {
		t.Fatalf("ExportFlushCommand() = %q", ExportFlushCommand())
	}
}

func TestDetectFromBeadsDir_EnvOverride(t *testing.T) {
	t.Setenv(envTool, "br")
	t.Cleanup(func() { SetTool("br") })

	dir := t.TempDir()
	beadsDir := filepath.Join(dir, ".beads")
	if err := os.MkdirAll(beadsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(beadsDir, "metadata.json"), []byte(`{"backend":"dolt"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	DetectFromBeadsDir(beadsDir)
	if Tool() != "br" {
		t.Fatalf("BV_BEADS_CLI should override workspace detection, got %q", Tool())
	}
}

func TestShellSubstitutesTool(t *testing.T) {
	t.Cleanup(func() { SetTool("br") })
	SetTool("bd")
	got := Shell("{tool} update ISSUE-1 --status=in_progress")
	want := "bd update ISSUE-1 --status=in_progress"
	if got != want {
		t.Fatalf("Shell() = %q, want %q", got, want)
	}
}

func TestTutorialLine_BD(t *testing.T) {
	t.Cleanup(func() { SetTool("br") })
	SetTool("bd")
	got := TutorialLine("br ready\nbr sync  # flush")
	if !strings.Contains(got, "bd ready") {
		t.Fatalf("TutorialLine() = %q, want bd ready", got)
	}
	if !strings.Contains(got, ExportFlushCommand()) {
		t.Fatalf("TutorialLine() should rewrite sync to export, got %q", got)
	}
}
