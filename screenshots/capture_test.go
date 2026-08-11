package screenshots_test

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Dicklesworthstone/beads_viewer/pkg/loader"
	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	"github.com/Dicklesworthstone/beads_viewer/pkg/ui"
)

func loadScreenshotIssues(t *testing.T) ([]model.Issue, string) {
	t.Helper()
	path := filepath.Join("..", "tests", "testdata", "synthetic_complex.jsonl")
	issues, err := loader.LoadIssuesFromFile(path)
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}
	return issues, path
}

func screenshotModel(t *testing.T) ui.Model {
	t.Helper()
	t.Setenv("BV_BACKGROUND_MODE", "0")
	issues, path := loadScreenshotIssues(t)
	m := ui.NewModel(issues, nil, path)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return updated.(ui.Model)
}

func TestScreenshotViewsRender(t *testing.T) {
	m := screenshotModel(t)
	views := []string{"list", "insights", "board", "graph"}
	for _, view := range views {
		out := m.RenderDebugView(view, 120, 40)
		if strings.TrimSpace(out) == "" {
			t.Fatalf("view %q rendered empty output", view)
		}
		if strings.Contains(out, "Unknown view:") {
			t.Fatalf("view %q failed: %s", view, out)
		}
		if strings.Contains(out, "Loading beads...") {
			t.Fatalf("view %q still on loading screen", view)
		}
	}
}

func TestScreenshotViewsNonEmptyLines(t *testing.T) {
	m := screenshotModel(t)
	out := m.RenderDebugView("list", 120, 40)
	if strings.Contains(out, "Loading beads...") {
		t.Fatal("list view still on loading screen")
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 5 {
		t.Fatalf("expected multi-line list view, got %d lines", len(lines))
	}
}
