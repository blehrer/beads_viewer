package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Dicklesworthstone/beads_viewer/pkg/icons"
	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) glyphHelpContextIssue() *model.Issue {
	switch {
	case m.isBoardView:
		return m.board.SelectedIssue()
	case m.isGraphView:
		return m.graphView.SelectedIssue()
	case m.focused == focusTree:
		return m.tree.SelectedIssue()
	case m.focused == focusInsights:
		if id := m.insightsPanel.SelectedIssueID(); id != "" {
			if iss, ok := m.issueByID(id); ok {
				return iss
			}
		}
	case m.isActionableView:
		if id := m.actionableView.SelectedIssueID(); id != "" {
			if iss, ok := m.issueByID(id); ok {
				return iss
			}
		}
	}
	if sel := m.list.SelectedItem(); sel != nil {
		if item, ok := sel.(IssueItem); ok {
			iss := item.Issue
			return &iss
		}
	}
	return nil
}

func (m Model) issueByID(id string) (*model.Issue, bool) {
	if m.issueMap != nil {
		if iss, ok := m.issueMap[id]; ok && iss != nil {
			return iss, true
		}
	}
	if m.snapshot != nil && m.snapshot.IssueMap != nil {
		if iss, ok := m.snapshot.IssueMap[id]; ok && iss != nil {
			return iss, true
		}
	}
	return nil, false
}

func (m Model) renderGlyphHelpOverlay() string {
	t := m.theme
	width := m.width - 4
	if width < 40 {
		width = 40
	}
	height := m.height - 4
	if height < 10 {
		height = 10
	}

	titleStyle := t.Renderer.NewStyle().Foreground(t.Primary).Bold(true)
	subStyle := t.Renderer.NewStyle().Foreground(t.Muted).Italic(true)
	sectionStyle := t.Renderer.NewStyle().Foreground(t.Secondary).Bold(true)
	glyphStyle := t.Renderer.NewStyle().Bold(true).Width(3)
	descStyle := t.Renderer.NewStyle().Foreground(t.Base.GetForeground())

	var lines []string
	lines = append(lines, titleStyle.Render("Symbol Reference"))
	lines = append(lines, subStyle.Render("K keyword help — j/k scroll · esc close · experimental.icon_set: nerd in config"))
	lines = append(lines, "")

	if issue := m.glyphHelpContextIssue(); issue != nil {
		lines = append(lines, sectionStyle.Render("Selected issue"))
		ctx := icons.ContextualGlossary(string(issue.Status), string(issue.IssueType), issue.Priority)
		for _, e := range ctx {
			lines = append(lines, formatGlossaryLine(e, glyphStyle, descStyle, width))
		}
		lines = append(lines, "")
	}

	grouped := icons.GlossaryByCategory()
	categories := make([]string, 0, len(grouped))
	for cat := range grouped {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		entries := grouped[cat]
		lines = append(lines, sectionStyle.Render(cat))
		for _, e := range entries {
			lines = append(lines, formatGlossaryLine(e, glyphStyle, descStyle, width))
		}
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	visibleLines := height - 4
	if visibleLines < 1 {
		visibleLines = 1
	}
	contentLines := strings.Split(content, "\n")
	maxScroll := len(contentLines) - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	scroll := m.glyphHelpScroll
	if scroll > maxScroll {
		scroll = maxScroll
	}
	if scroll < 0 {
		scroll = 0
	}
	end := scroll + visibleLines
	if end > len(contentLines) {
		end = len(contentLines)
	}
	viewContent := strings.Join(contentLines[scroll:end], "\n")

	boxStyle := t.Renderer.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(1, 2).
		Width(width).
		Height(height)

	rendered := boxStyle.Render(viewContent)
	return lipgloss.Place(m.width, m.height-1, lipgloss.Center, lipgloss.Center, rendered)
}

func formatGlossaryLine(e icons.GlossaryEntry, glyphStyle, descStyle lipgloss.Style, width int) string {
	glyph := icons.Get(e.Name)
	if st := statusForIconName(e.Name); st != "" {
		glyph = RenderStatusDot(st)
	} else if p, ok := priorityForIconName(e.Name); ok {
		glyph = RenderPriorityIcon(p)
	}
	glyphPart := glyphStyle.Render(glyph)
	descWidth := width - lipgloss.Width(glyphPart) - 4
	if descWidth < 20 {
		descWidth = 20
	}
	desc := descStyle.Width(descWidth).Render(e.Description)
	return fmt.Sprintf("  %s  %s", glyphPart, desc)
}

func statusForIconName(name icons.Name) string {
	switch name {
	case icons.StatusOpen:
		return "open"
	case icons.StatusInProgress:
		return "in_progress"
	case icons.StatusBlocked:
		return "blocked"
	case icons.StatusClosed:
		return "closed"
	case icons.StatusUnknown:
		return "unknown"
	case icons.Pause:
		return "deferred"
	case icons.StatusPinned:
		return "pinned"
	case icons.StatusHooked:
		return "hooked"
	case icons.StatusReview:
		return "review"
	default:
		return ""
	}
}

func priorityForIconName(name icons.Name) (int, bool) {
	switch name {
	case icons.Fire:
		return 0, true
	case icons.Lightning:
		return 1, true
	case icons.PriorityMedium:
		return 2, true
	case icons.PriorityLow:
		return 3, true
	case icons.PriorityBacklog:
		return 4, true
	default:
		return 0, false
	}
}

func (m Model) handleGlyphHelpKeys(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "j", "down":
		m.glyphHelpScroll++
	case "k", "up":
		if m.glyphHelpScroll > 0 {
			m.glyphHelpScroll--
		}
	case "ctrl+d":
		m.glyphHelpScroll += 10
	case "ctrl+u":
		if m.glyphHelpScroll > 10 {
			m.glyphHelpScroll -= 10
		} else {
			m.glyphHelpScroll = 0
		}
	case "home", "g":
		m.glyphHelpScroll = 0
	case "G", "end":
		m.glyphHelpScroll = 9999
	case "q", "esc", "K":
		m.showGlyphHelp = false
		m.glyphHelpScroll = 0
		m.focused = m.restoreFocusFromHelp()
	default:
		m.showGlyphHelp = false
		m.glyphHelpScroll = 0
		m.focused = m.restoreFocusFromHelp()
	}
	return m
}
