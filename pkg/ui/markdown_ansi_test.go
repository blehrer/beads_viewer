package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Dicklesworthstone/beads_viewer/pkg/analysis"
	"github.com/Dicklesworthstone/beads_viewer/pkg/correlation"
	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	"github.com/Dicklesworthstone/beads_viewer/pkg/testutil"
)

// Regression: lipgloss-styled icons must never reach glamour markdown sources.
//
// Glamour inputs audited in pkg/ui:
//   - buildListDetailMarkdown (includes renderBeadHistoryMD)
//   - buildDetailMarkdown / renderCalculationProofMD (insights detail panel)
//   - buildBoardDetailMarkdown, board detail help placeholder
//   - insights empty-detail placeholder, metricDescriptions.WhatIs (panel explanations)
//   - RenderDependencyTree (embedded in list detail code blocks)
//   - tutorial page.Content (TutorialModel.markdownRenderer)
//
// Intentional exemptions (not builder output — user/content data only):
//   - Board expanded card description passed to mdRenderer.Render(desc) in board.go
//   - Issue Description/Design/Notes fields echoed verbatim into markdown builders
//
// Checked for ANSI hygiene but not passed through glamour:
//   - copyIssueToClipboard markdown (clipboard export)
const (
	// Mirrors board.go renderDetailPanel help text when no card is selected.
	boardDetailHelpMarkdown = "## No Selection\n\nNavigate to a card with **h/l** and **j/k** to see details here.\n\nPress **Tab** to hide this panel."

	// Mirrors insights.go renderDetailPanel emptyContent when nothing is selected.
	insightsEmptyDetailMarkdown = `
## Select a Bead

Navigate to a metric panel and select an item to view its details here.

**Navigation:**
- ← → to switch panels
- ↑ ↓ to select items
- Ctrl+j/k scroll details
- Enter to view in main view
`
)

func TestMarkdownSources_NoANSI(t *testing.T) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	blocker := model.Issue{
		ID: "blocker-1", Title: "Blocker", Status: model.StatusOpen, Priority: 0,
	}
	child := model.Issue{
		ID: "child-1", Title: "Child task", Status: model.StatusBlocked, Priority: 1,
		IssueType: model.TypeBug, Assignee: "alice", Labels: []string{"backend", "urgent"},
		Description:          "Do the thing",
		Design:               "Use pattern X",
		AcceptanceCriteria:   "Tests pass",
		Notes:                "Watch edge cases",
		CreatedAt:            now,
		UpdatedAt:            now,
		Dependencies: []*model.Dependency{
			{DependsOnID: "blocker-1", Type: model.DepBlocks},
		},
		Comments: []*model.Comment{
			{Author: "bob", Text: "Looks good", CreatedAt: now},
		},
	}
	issues := []model.Issue{blocker, child}

	issueMap := make(map[string]*model.Issue, len(issues))
	for i := range issues {
		issueMap[issues[i].ID] = &issues[i]
	}

	issueItem := IssueItem{
		Issue:           child,
		TriageScore:     0.85,
		TriageReason:    "High impact",
		TriageReasons:   []string{"High impact", "Unblocks work"},
		IsQuickWin:      true,
		IsBlocker:       false,
		UnblocksCount:   3,
		SearchScoreSet:  true,
		SearchScore:     0.72,
		SearchTextScore: 0.55,
		SearchComponents: map[string]float64{
			"pagerank": 0.4, "status": 0.2, "impact": 0.1,
		},
	}

	t.Run("list detail markdown", func(t *testing.T) {
		m := NewModel(issues, nil, "")
		m.semanticSearchEnabled = true
		m.semanticHybridEnabled = true
		md := m.buildListDetailMarkdown(issueItem, child)
		testutil.AssertNoANSI(t, "buildListDetailMarkdown", md)
	})

	t.Run("list detail markdown with history", func(t *testing.T) {
		m := NewModel(issues, nil, "")
		m.historyView = NewHistoryModel(childHistoryReport(now), DefaultTheme(nil))
		md := m.buildListDetailMarkdown(issueItem, child)
		testutil.AssertNoANSI(t, "buildListDetailMarkdown+history", md)
	})

	t.Run("bead history markdown", func(t *testing.T) {
		m := NewModel(issues, nil, "")
		m.historyView = NewHistoryModel(childHistoryReport(now), DefaultTheme(nil))
		md := m.renderBeadHistoryMD("child-1")
		testutil.AssertNoANSI(t, "renderBeadHistoryMD", md)
	})

	t.Run("insights detail markdown", func(t *testing.T) {
		stats := analysis.NewAnalyzer(issues).Analyze()
		ins := analysis.Insights{Stats: &stats}
		m := NewInsightsModel(ins, issueMap, DefaultTheme(nil))
		for _, id := range []string{"child-1", "blocker-1"} {
			md := m.buildDetailMarkdown(id)
			testutil.AssertNoANSI(t, "buildDetailMarkdown("+id+")", md)
		}
	})

	t.Run("insights calculation proof all panels", func(t *testing.T) {
		stats := analysis.NewAnalyzer(issues).Analyze()
		ins := analysis.Insights{
			Stats:  &stats,
			Cycles: [][]string{{"child-1", "blocker-1", "child-1"}},
		}
		m := NewInsightsModel(ins, issueMap, DefaultTheme(nil))
		m.selectedIndex[PanelCycles] = 0
		for panel := PanelBottlenecks; panel < PanelCount; panel++ {
			m.focusedPanel = panel
			md := m.renderCalculationProofMD("child-1")
			testutil.AssertNoANSI(t, fmt.Sprintf("renderCalculationProofMD(%d)", panel), md)
		}
	})

	t.Run("board detail markdown", func(t *testing.T) {
		blocksIndex := map[string][]string{"blocker-1": {"child-1"}}
		for _, prio := range []int{0, 1, 2, 3, 4} {
			iss := child
			iss.Priority = prio
			iss.Status = model.StatusOpen
			md := buildBoardDetailMarkdown(&iss, issueMap, blocksIndex, GetTypeIconMD(string(iss.IssueType)))
			testutil.AssertNoANSI(t, fmt.Sprintf("buildBoardDetailMarkdown P%d", prio), md)
		}
	})

	t.Run("board detail help placeholder", func(t *testing.T) {
		testutil.AssertNoANSI(t, "boardDetailHelpMarkdown", boardDetailHelpMarkdown)
	})

	t.Run("insights empty detail placeholder", func(t *testing.T) {
		testutil.AssertNoANSI(t, "insightsEmptyDetailMarkdown", insightsEmptyDetailMarkdown)
	})

	t.Run("metric panel explanations", func(t *testing.T) {
		for panel := PanelBottlenecks; panel < PanelCount; panel++ {
			info := metricDescriptions[panel]
			testutil.AssertNoANSI(t, fmt.Sprintf("metricDescriptions[%d].WhatIs", panel), info.WhatIs)
		}
	})

	t.Run("copy issue markdown", func(t *testing.T) {
		md := buildCopyIssueMarkdown(child)
		testutil.AssertNoANSI(t, "buildCopyIssueMarkdown", md)
	})

	t.Run("tutorial page content", func(t *testing.T) {
		for _, page := range defaultTutorialPages() {
			testutil.AssertNoANSI(t, "tutorial:"+page.ID, page.Content)
		}
	})

	t.Run("dependency tree in markdown", func(t *testing.T) {
		root := BuildDependencyTree("child-1", issueMap, 3)
		tree := RenderDependencyTree(root)
		testutil.AssertNoANSI(t, "RenderDependencyTree", tree)
	})
}

func TestStyledTUIIcons_DifferFromMarkdownIcons(t *testing.T) {
	// In CI/non-TTY lipgloss may omit ESC sequences; styled vs plain glyph still differs
	// when color is applied, and MD helpers must never carry escapes.
	if GetPriorityIcon(0) == GetPriorityIconMD(0) && GetStatusIcon("open") == GetStatusIconMD("open") {
		t.Fatal("styled icons should differ from markdown-safe icons")
	}
	if testutil.ContainsANSI(GetPriorityIconMD(0)) || testutil.ContainsANSI(GetStatusIconMD("open")) {
		t.Fatal("markdown icon helpers must not contain ANSI")
	}
}

func childHistoryReport(now time.Time) *correlation.HistoryReport {
	return &correlation.HistoryReport{
		GeneratedAt: now,
		Histories: map[string]correlation.BeadHistory{
			"child-1": {
				BeadID: "child-1",
				Title:  "Child task",
				Status: "blocked",
				Events: []correlation.BeadEvent{
					{EventType: correlation.EventCreated, Author: "alice", Timestamp: now},
					{EventType: correlation.EventClaimed, Author: "bob", Timestamp: now.Add(-time.Hour)},
					{EventType: correlation.EventClosed, Author: "bob", Timestamp: now.Add(-2 * time.Hour)},
				},
				Commits: []correlation.CorrelatedCommit{
					{
						ShortSHA:   "abc123d",
						Message:    "fix: child task",
						Author:     "bob",
						Timestamp:  now,
						Confidence: 0.95,
						Files: []correlation.FileChange{
							{Path: "pkg/ui/model.go", Insertions: 3, Deletions: 1},
						},
					},
					{
						ShortSHA:   "def456g",
						Message:    "test: add coverage",
						Author:     "alice",
						Timestamp:  now.Add(-time.Hour),
						Confidence: 0.35,
					},
				},
			},
		},
	}
}

// buildCopyIssueMarkdown mirrors copyIssueToClipboard markdown assembly (clipboard-only path).
func buildCopyIssueMarkdown(issue model.Issue) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s %s\n\n", GetTypeIconMD(string(issue.IssueType)), issue.Title))
	sb.WriteString(fmt.Sprintf("**ID:** %s  \n", issue.ID))
	sb.WriteString(fmt.Sprintf("**Status:** %s  \n", strings.ToUpper(string(issue.Status))))
	sb.WriteString(fmt.Sprintf("**Priority:** P%d  \n", issue.Priority))
	if issue.Assignee != "" {
		sb.WriteString(fmt.Sprintf("**Assignee:** @%s  \n", issue.Assignee))
	}
	sb.WriteString(fmt.Sprintf("**Created:** %s  \n", issue.CreatedAt.Format("2006-01-02")))

	if len(issue.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("**Labels:** %s  \n", strings.Join(issue.Labels, ", ")))
	}

	if issue.Description != "" {
		sb.WriteString(fmt.Sprintf("\n## Description\n\n%s\n", issue.Description))
	}

	if issue.AcceptanceCriteria != "" {
		sb.WriteString(fmt.Sprintf("\n## Acceptance Criteria\n\n%s\n", issue.AcceptanceCriteria))
	}

	if len(issue.Dependencies) > 0 {
		sb.WriteString("\n## Dependencies\n\n")
		for _, dep := range issue.Dependencies {
			if dep == nil {
				continue
			}
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", dep.DependsOnID, dep.Type))
		}
	}

	return sb.String()
}
