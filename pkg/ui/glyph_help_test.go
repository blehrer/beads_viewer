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

func TestGlyphHelpScrollKeys(t *testing.T) {
	issues := []model.Issue{
		{ID: "bv-1", Title: "Test", Status: model.StatusOpen, IssueType: model.TypeBug, Priority: 0},
	}
	m := NewModel(issues, nil, "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("K")})
	m = updated.(Model)
	if !m.showGlyphHelp {
		t.Fatal("expected glyph help open")
	}

	keys := []struct {
		key  tea.KeyMsg
		want int
	}{
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}, 1},
		{tea.KeyMsg{Type: tea.KeyDown, Runes: []rune("down")}, 2},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}, 1},
		{tea.KeyMsg{Type: tea.KeyUp, Runes: []rune("up")}, 0},
		{tea.KeyMsg{Type: tea.KeyPgDown}, 32},
		{tea.KeyMsg{Type: tea.KeyPgUp}, 0},
	}
	for _, tc := range keys {
		updated, _ = m.Update(tc.key)
		m = updated.(Model)
		if !m.showGlyphHelp {
			t.Fatalf("key %q closed overlay unexpectedly", tc.key.String())
		}
		if m.glyphHelpScroll != tc.want {
			t.Fatalf("key %q: scroll = %d, want %d (focus=%v)", tc.key.String(), m.glyphHelpScroll, tc.want, m.focused)
		}
	}
}

func TestGlyphHelpScrollKeysFromBoardView(t *testing.T) {
	issues := []model.Issue{
		{ID: "bv-1", Title: "Test", Status: model.StatusOpen, IssueType: model.TypeTask, Priority: 1},
	}
	m := NewModel(issues, nil, "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	m.isBoardView = true
	m.focused = focusBoard

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("K")})
	m = updated.(Model)
	if !m.showGlyphHelp || m.focused != focusGlyphHelp {
		t.Fatalf("expected glyph help from board, show=%v focus=%v", m.showGlyphHelp, m.focused)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	m = updated.(Model)
	if !m.showGlyphHelp || m.glyphHelpScroll != 1 {
		t.Fatalf("j from board: show=%v scroll=%d focus=%v", m.showGlyphHelp, m.glyphHelpScroll, m.focused)
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
