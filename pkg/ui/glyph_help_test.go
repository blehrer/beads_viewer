package ui

import (
	"strings"
	"testing"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	tea "github.com/charmbracelet/bubbletea"
)

func TestGlyphHelpToggle(t *testing.T) {
	issues := []model.Issue{
		{ID: "bv-1", Title: "Test", Status: model.StatusOpen, IssueType: model.TypeBug, Priority: 0},
	}
	m := NewModel(issues, nil, "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("K")})
	m = updated.(Model)
	if !m.showGlyphHelp || m.focused != focusGlyphHelp {
		t.Fatalf("expected glyph help open, show=%v focus=%v", m.showGlyphHelp, m.focused)
	}

	out := m.renderGlyphHelpOverlay()
	if !strings.Contains(out, "Symbol Reference") {
		t.Fatal("glyph help should include title")
	}
	if !strings.Contains(out, "Selected issue") {
		t.Fatal("glyph help should show contextual section for selected issue")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")})
	m = updated.(Model)
	if m.showGlyphHelp || m.focused != focusList {
		t.Fatalf("expected glyph help closed, show=%v focus=%v", m.showGlyphHelp, m.focused)
	}
}

func TestGlyphHelpNotInHistory(t *testing.T) {
	issues := []model.Issue{{ID: "1", Title: "One", Status: model.StatusOpen}}
	m := NewModel(issues, nil, "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	m.isHistoryView = true
	m.focused = focusHistory

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("K")})
	m = updated.(Model)
	if m.showGlyphHelp {
		t.Fatal("K should not open glyph help in history view (J/K reserved)")
	}
}
