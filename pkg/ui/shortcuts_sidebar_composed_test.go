package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Dicklesworthstone/beads_viewer/pkg/drift"
)

func maxComposedLineWidth(body, sidebar string) int {
	mx := 0
	joined := lipgloss.JoinHorizontal(lipgloss.Top, body, sidebar)
	for _, ln := range strings.Split(joined, "\n") {
		if w := lipgloss.Width(ln); w > mx {
			mx = w
		}
	}
	return mx
}

func composedSidebar(m Model) string {
	m.shortcutsSidebar.SetFocus(m.focused)
	m.shortcutsSidebar.SetSize(m.shortcutsSidebar.Width(), m.height-2)
	return m.shortcutsSidebar.View()
}

func composedListWidth(m Model, showDetails bool) int {
	var body string
	if m.isSplitView {
		body = m.renderSplitView()
	} else if showDetails {
		body = m.viewport.View()
	} else {
		body = m.renderListWithHeader()
	}
	return maxComposedLineWidth(body, composedSidebar(m))
}

func composedActionableWidth(m Model) int {
	cw := m.mainContentWidth()
	bodyH := m.height - 1
	actionableH := bodyH - 1
	if actionableH < 3 {
		actionableH = 3
	}
	m.actionableView.SetSize(cw, actionableH)
	body := m.actionableView.Render()
	return maxComposedLineWidth(body, composedSidebar(m))
}

func composedAlertsWidth(m Model) int {
	body := m.renderAlertsPanel()
	return maxComposedLineWidth(body, composedSidebar(m))
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

			// Open the sidebar via the real `;` key path.
			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(";")})
			m = updated.(Model)
			if !m.showShortcutsSidebar {
				t.Fatalf("`;` did not enable the shortcuts sidebar")
			}

			// List/board body and (mobile) detail body must both fit once the
			// sidebar column is appended.
			if cw := composedListWidth(m, false); cw > m.width {
				t.Errorf("list body + sidebar composed width %d exceeds terminal width %d (#168 overflow)", cw, m.width)
			}
			if !m.isSplitView {
				if cw := composedListWidth(m, true); cw > m.width {
					t.Errorf("detail body + sidebar composed width %d exceeds terminal width %d (#168 overflow)", cw, m.width)
				}
			}
		})
	}
}

func TestShortcutsSidebarComposedWidth_ActionableAndAlerts(t *testing.T) {
	cases := []struct {
		name string
		w, h int
	}{
		{"actionable_120", 120, 30},
		{"actionable_80", 80, 24},
		{"alerts_120", 120, 30},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := sizedModel(t, mouseTestIssues(40), tc.w, tc.h)

			if strings.HasPrefix(tc.name, "actionable") {
				updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
				m = updated.(Model)
				if !m.isActionableView {
					t.Fatal("a did not open actionable view")
				}
				updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(";")})
				m = updated.(Model)
				if cw := composedActionableWidth(m); cw > m.width {
					t.Errorf("actionable + sidebar width %d exceeds terminal %d", cw, m.width)
				}
				return
			}

			m.alerts = []drift.Alert{{
				Type:     drift.AlertStaleIssue,
				Severity: drift.SeverityWarning,
				Message:  "Stale issue detected",
				IssueID:  "bv-test",
			}}
			m.showAlertsPanel = true
			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(";")})
			m = updated.(Model)
			if cw := composedAlertsWidth(m); cw > m.width {
				t.Errorf("alerts + sidebar width %d exceeds terminal %d", cw, m.width)
			}
		})
	}
}
