package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Dicklesworthstone/beads_viewer/pkg/drift"
)

func maxLineWidth(s string) int {
	mx := 0
	for _, ln := range strings.Split(s, "\n") {
		if w := lipgloss.Width(ln); w > mx {
			mx = w
		}
	}
	return mx
}

// composedBodyWidth reconstructs tier-0/2 body + tier-1 sidebar join before the
// final View() clamp, matching the production layout pipeline (bv-sl44.4).
func composedBodyWidth(m Model) int {
	cw := m.mainContentWidth()
	bodyH := m.height - 1

	var body string
	if overlay, ok := m.renderOverlay(cw, bodyH); ok {
		body = overlay
	} else {
		body = m.renderBaseView(cw, bodyH)
	}
	return maxLineWidth(m.joinShortcutsSidebar(body))
}

func enableSidebar(t *testing.T, m Model) Model {
	t.Helper()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(";")})
	m = updated.(Model)
	if !m.showShortcutsSidebar {
		t.Fatal("`;` did not enable the shortcuts sidebar")
	}
	return m
}

func assertComposedFitsTerminal(t *testing.T, m Model, label string) {
	t.Helper()
	if cw := composedBodyWidth(m); cw > m.width {
		t.Errorf("%s: composed width %d exceeds terminal width %d (#168 overflow)", label, cw, m.width)
	}
}

func keyModel(t *testing.T, m Model, key rune) Model {
	t.Helper()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
	return updated.(Model)
}

// TestShortcutsSidebarComposedWidthFitsTerminal is a stronger companion to
// TestShortcutsSidebarReservesLayoutWidth. That test measures m.View(), whose
// final lipgloss clamp to m.width hides any over-wide composition inside a
// string buffer — exactly the failure mode (#168) where a real terminal wraps
// the overflow back into the panes. This test instead reconstructs the
// body+sidebar JoinHorizontal BEFORE the final clamp and asserts it fits within
// m.width, so it catches reservation-math drift (e.g. forgetting the sidebar's
// rendered border columns, or a body path that ignores mainContentWidth()).
func TestShortcutsSidebarComposedWidthFitsTerminal(t *testing.T) {
	cases := []struct {
		name      string
		w, h      int
		wantSplit bool
	}{
		{"split_narrow_110", 110, 30, true},
		{"split_120", 120, 30, true},
		{"split_wide_200", 200, 40, true},
		{"mobile_80", 80, 30, false},
		{"mobile_60", 60, 24, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := sizedModel(t, mouseTestIssues(40), tc.w, tc.h)
			if m.isSplitView != tc.wantSplit {
				t.Fatalf("w=%d isSplitView=%v want %v", tc.w, m.isSplitView, tc.wantSplit)
			}

			m = enableSidebar(t, m)

			assertComposedFitsTerminal(t, m, "list")
			if !m.isSplitView {
				m.showDetails = true
				assertComposedFitsTerminal(t, m, "detail")
			}
		})
	}
}

// TestShortcutsSidebarComposedWidth_Matrix verifies every major body and overlay
// combination fits the terminal once the shortcuts sidebar is docked (bv-sl44.5).
func TestShortcutsSidebarComposedWidth_Matrix(t *testing.T) {
	testAlert := drift.Alert{
		Type:     drift.AlertStaleIssue,
		Severity: drift.SeverityWarning,
		Message:  "Stale issue detected",
		IssueID:  "bv-test",
	}

	cases := []struct {
		name  string
		w, h  int
		setup func(t *testing.T, m Model) Model
	}{
		{
			name: "list_mobile",
			w:    80, h: 30,
			setup: func(t *testing.T, m Model) Model { return enableSidebar(t, m) },
		},
		{
			name: "split",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model { return enableSidebar(t, m) },
		},
		{
			name: "actionable",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model {
				m = keyModel(t, m, 'a')
				if !m.isActionableView {
					t.Fatal("a did not open actionable view")
				}
				return enableSidebar(t, m)
			},
		},
		{
			name: "board",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model {
				m = keyModel(t, m, 'b')
				if !m.isBoardView {
					t.Fatal("b did not open board view")
				}
				return enableSidebar(t, m)
			},
		},
		{
			name: "graph",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model {
				m = keyModel(t, m, 'g')
				if !m.isGraphView {
					t.Fatal("g did not open graph view")
				}
				return enableSidebar(t, m)
			},
		},
		{
			name: "insights",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model {
				m = keyModel(t, m, 'i')
				if m.focused != focusInsights {
					t.Fatalf("i did not open insights view, focused=%v", m.focused)
				}
				return enableSidebar(t, m)
			},
		},
		{
			name: "alerts",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model {
				m.alerts = []drift.Alert{testAlert}
				m = keyModel(t, m, '!')
				if !m.showAlertsPanel {
					t.Fatal("! did not open alerts panel")
				}
				return enableSidebar(t, m)
			},
		},
		{
			name: "help",
			w:    120, h: 30,
			setup: func(t *testing.T, m Model) Model {
				m = keyModel(t, m, '?')
				if !m.showHelp {
					t.Fatal("? did not open help overlay")
				}
				return enableSidebar(t, m)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := sizedModel(t, mouseTestIssues(40), tc.w, tc.h)
			m = tc.setup(t, m)
			assertComposedFitsTerminal(t, m, tc.name)
		})
	}
}
