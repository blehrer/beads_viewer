package ui

import (
	"fmt"
	"strings"

	"github.com/Dicklesworthstone/beads_viewer/pkg/icons"
	"github.com/charmbracelet/lipgloss"
)

// ══════════════════════════════════════════════════════════════════════════════
// DESIGN TOKENS - Consistent spacing, colors, and visual language
// ══════════════════════════════════════════════════════════════════════════════

// Spacing constants for consistent layout (in characters)
const (
	SpaceXS = 1
	SpaceSM = 2
	SpaceMD = 3
	SpaceLG = 4
	SpaceXL = 6
)

// ══════════════════════════════════════════════════════════════════════════════
// COLOR PALETTE - Populated from ActivePalette() via syncStyleGlobalsFromPalette.
// Defaults and user overrides live in palette.go.
// ══════════════════════════════════════════════════════════════════════════════

var (
	ColorBg, ColorBgDark, ColorBgSubtle, ColorBgHighlight lipgloss.AdaptiveColor
	ColorText, ColorSubtext, ColorMuted                   lipgloss.AdaptiveColor

	ColorPrimary, ColorSecondary                       lipgloss.AdaptiveColor
	ColorInfo, ColorSuccess, ColorWarning, ColorDanger lipgloss.AdaptiveColor

	ColorStatusOpen, ColorStatusInProgress, ColorStatusBlocked lipgloss.AdaptiveColor
	ColorStatusDeferred, ColorStatusPinned, ColorStatusHooked  lipgloss.AdaptiveColor
	ColorStatusReview, ColorStatusClosed, ColorStatusTombstone lipgloss.AdaptiveColor

	ColorStatusOpenBg, ColorStatusInProgressBg, ColorStatusBlockedBg lipgloss.AdaptiveColor
	ColorStatusDeferredBg, ColorStatusPinnedBg, ColorStatusHookedBg  lipgloss.AdaptiveColor
	ColorStatusReviewBg, ColorStatusClosedBg, ColorStatusTombstoneBg lipgloss.AdaptiveColor

	ColorPrioCritical, ColorPrioHigh, ColorPrioMedium, ColorPrioLow         lipgloss.AdaptiveColor
	ColorPrioCriticalBg, ColorPrioHighBg, ColorPrioMediumBg, ColorPrioLowBg lipgloss.AdaptiveColor

	ColorTypeBug, ColorTypeFeature, ColorTypeTask, ColorTypeEpic, ColorTypeChore lipgloss.AdaptiveColor

	ColorFooterHint, ColorFooterKey, ColorFooterSep, ColorFooterDim lipgloss.AdaptiveColor
)

// ══════════════════════════════════════════════════════════════════════════════
// PANEL STYLES - For split view layouts
// ══════════════════════════════════════════════════════════════════════════════

var (
	// PanelStyle is the default style for unfocused panels (rebuilt on palette init).
	PanelStyle lipgloss.Style
	// FocusedPanelStyle is the style for focused panels (rebuilt on palette init).
	FocusedPanelStyle lipgloss.Style
)

// ══════════════════════════════════════════════════════════════════════════════
// BADGE RENDERING - Polished, consistent badge styles
// ══════════════════════════════════════════════════════════════════════════════

// RenderPriorityBadge returns a styled priority badge
// Priority values: 0=Critical, 1=High, 2=Medium, 3=Low, 4=Backlog
func RenderPriorityBadge(priority int) string {
	var fg, bg lipgloss.AdaptiveColor
	var label string

	switch priority {
	case 0:
		fg, bg, label = ColorPrioCritical, ColorPrioCriticalBg, "P0"
	case 1:
		fg, bg, label = ColorPrioHigh, ColorPrioHighBg, "P1"
	case 2:
		fg, bg, label = ColorPrioMedium, ColorPrioMediumBg, "P2"
	case 3:
		fg, bg, label = ColorPrioLow, ColorPrioLowBg, "P3"
	case 4:
		fg, bg, label = ColorMuted, ColorBgSubtle, "P4"
	default:
		fg, bg, label = ColorMuted, ColorBgSubtle, "P?"
	}

	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Bold(true).
		Padding(0, 0).
		Render(label)
}

func lipStyle() lipgloss.Style {
	return defaultRenderer.NewStyle()
}

// renderMarkerGlyph colors circle-replacement markers for the TUI.
// Nerd Font glyphs get a subtle badge background (like OPEN badges) because
// foreground-only coloring is easy to miss on complex NF shapes; plain ●○◉
// keep foreground-only styling.
func renderMarkerGlyph(glyph string, fg, bg lipgloss.AdaptiveColor) string {
	style := lipStyle().Foreground(fg)
	if icons.ActiveSet() == icons.SetNerd {
		style = style.Background(bg)
	}
	return style.Render(glyph)
}

// statusDotColor maps a beads status to the theme foreground used for list/graph dots.
func statusDotColor(status string) lipgloss.AdaptiveColor {
	switch status {
	case "open":
		return ColorStatusOpen
	case "in_progress":
		return ColorStatusInProgress
	case "blocked":
		return ColorStatusBlocked
	case "deferred", "draft":
		return ColorStatusDeferred
	case "pinned":
		return ColorStatusPinned
	case "hooked":
		return ColorStatusHooked
	case "review":
		return ColorStatusReview
	case "closed":
		return ColorStatusOpen // green ✓ — done/success
	case "tombstone":
		return ColorStatusTombstone
	default:
		return ColorMuted
	}
}

// statusDotBgColor pairs with statusDotColor for nerd-mode marker badges.
func statusDotBgColor(status string) lipgloss.AdaptiveColor {
	switch status {
	case "open":
		return ColorStatusOpenBg
	case "in_progress":
		return ColorStatusInProgressBg
	case "blocked":
		return ColorStatusBlockedBg
	case "deferred", "draft":
		return ColorStatusDeferredBg
	case "pinned":
		return ColorStatusPinnedBg
	case "hooked":
		return ColorStatusHookedBg
	case "review":
		return ColorStatusReviewBg
	case "closed":
		return ColorStatusOpenBg
	case "tombstone":
		return ColorStatusTombstoneBg
	default:
		return ColorBgSubtle
	}
}

// statusNerdGlyph returns the registry glyph for a status when BV_ICON_SET=nerd.
func statusNerdGlyph(status string) string {
	switch status {
	case "open":
		return icons.Get(icons.StatusGraphOpen)
	case "in_progress":
		return icons.Get(icons.StatusInProgress)
	case "blocked":
		return icons.Get(icons.StatusBlocked)
	case "deferred", "draft":
		return icons.Get(icons.Pause)
	case "pinned":
		return icons.Get(icons.StatusPinned)
	case "hooked":
		return icons.Get(icons.StatusHooked)
	case "review":
		return icons.Get(icons.StatusReview)
	case "closed":
		return icons.Get(icons.Check)
	case "tombstone":
		return icons.Get(icons.StatusUnknown)
	default:
		return icons.Get(icons.StatusUnknown)
	}
}

// timelineMilestoneColor matches the theme colors used for ○ ● ✓ milestone markers.
func timelineMilestoneColor(kind string) lipgloss.AdaptiveColor {
	switch kind {
	case "created":
		return ColorStatusOpen
	case "claimed":
		return ColorStatusInProgress
	case "closed":
		return ColorStatusOpen
	default:
		return ColorMuted
	}
}

func timelineMilestoneBgColor(kind string) lipgloss.AdaptiveColor {
	switch kind {
	case "created":
		return ColorStatusOpenBg
	case "claimed":
		return ColorStatusInProgressBg
	case "closed":
		return ColorStatusOpenBg
	default:
		return ColorBgSubtle
	}
}

// footerStatColor matches footer count indicator colors (○◉◈● semantics).
func footerStatColor(kind string) lipgloss.AdaptiveColor {
	switch kind {
	case "open":
		return ColorStatusOpen
	case "ready":
		return ColorSuccess
	case "blocked":
		return ColorWarning
	case "closed":
		return ColorStatusOpen
	default:
		return ColorMuted
	}
}

func footerStatBgColor(kind string) lipgloss.AdaptiveColor {
	switch kind {
	case "open":
		return ColorStatusOpenBg
	case "ready":
		return ColorPrioLowBg
	case "blocked":
		return ColorStatusBlockedBg
	case "closed":
		return ColorStatusOpenBg
	default:
		return ColorBgSubtle
	}
}

func lifecycleEventColors(eventType string) (fg, bg lipgloss.AdaptiveColor) {
	switch eventType {
	case "created":
		return ColorPrimary, ColorStatusReviewBg
	case "claimed":
		return ColorStatusInProgress, ColorStatusInProgressBg
	case "closed":
		return ColorStatusOpen, ColorStatusOpenBg
	case "reopened":
		return ColorSecondary, ColorStatusClosedBg
	case "modified":
		return ColorMuted, ColorBgSubtle
	default:
		return ColorMuted, ColorBgSubtle
	}
}

// RenderStatusDot returns a lipgloss-colored status dot for interactive TUI views.
// The glyph is a plain ● with theme foreground color (not emoji status circles).
// Do not embed the result in markdown passed to glamour or static export;
// use GetStatusIconMD for plain glyphs without ANSI escapes.
//
// ponytail: plain ● + theme color replaces emoji circles (🟢🔵🔴) — no VS-16 width issues.
func RenderStatusDot(status string) string {
	if status == "" {
		status = "unknown"
	}
	var glyph string
	switch status {
	case "closed":
		glyph = "✓"
		if icons.ActiveSet() == icons.SetNerd {
			glyph = icons.Get(icons.Check)
		}
	default:
		glyph = "●"
		if icons.ActiveSet() == icons.SetNerd {
			glyph = statusNerdGlyph(status)
		}
	}
	return renderMarkerGlyph(glyph, statusDotColor(status), statusDotBgColor(status))
}

// RenderStatusDotGraph is like RenderStatusDot but uses ✓ for completed issues in graph view.
func RenderStatusDotGraph(status string) string {
	return RenderStatusDot(status)
}

// priorityIconColor maps beads priority 0–4 to theme foregrounds for TUI glyphs.
func priorityIconColor(priority int) lipgloss.AdaptiveColor {
	switch priority {
	case 0:
		return ColorPrioCritical
	case 1:
		return ColorPrioHigh
	case 2:
		return ColorPrioMedium
	case 3:
		return ColorPrioLow
	case 4:
		return ColorMuted
	default:
		return ColorMuted
	}
}

// RenderPriorityIcon returns a lipgloss-colored priority glyph for interactive TUI views.
// The base glyph comes from pkg/icons; lipgloss adds theme foreground color.
// Do not embed the result in markdown passed to glamour or static export;
// use GetPriorityIconMD for plain glyphs without ANSI escapes.
//
// ponytail: emoji/nerd glyph + theme color replaces raw 🔹/⚡ (P2 rhombus is blue).
func RenderPriorityIcon(priority int) string {
	glyph := icons.Priority(priority)
	if strings.TrimSpace(glyph) == "" {
		return glyph
	}
	return lipStyle().
		Foreground(priorityIconColor(priority)).
		Render(glyph)
}

// RenderTriageScoreIcon colors triage score indicators (same nerd glyph as P2, emoji 🟠).
func RenderTriageScoreIcon(name icons.Name) string {
	var color lipgloss.AdaptiveColor
	switch name {
	case icons.StatusBlocked:
		color = ColorPrioCritical
	case icons.TriageScoreMid:
		color = ColorPrioHigh // orange band — not P2 blue
	default:
		color = ColorStatusInProgress
	}
	return lipStyle().Foreground(color).Render(icons.Get(name))
}

// RenderHistoryBeadStatus returns a colored history-list status glyph (○●✓ or nerd equivalent).
func RenderHistoryBeadStatus(status string) string {
	norm := normalizeHistoryStatus(status)
	return renderMarkerGlyph(icons.HistoryBeadStatus(status), statusDotColor(norm), statusDotBgColor(norm))
}

func normalizeHistoryStatus(status string) string {
	switch status {
	case "closed", "tombstone":
		return "closed"
	case "in_progress":
		return "in_progress"
	case "blocked":
		return "blocked"
	default:
		return "open"
	}
}

// RenderTimelineMilestone returns a colored lifecycle milestone marker for TUI timelines.
func RenderTimelineMilestone(kind string) string {
	return renderMarkerGlyph(
		icons.TimelineMilestone(kind),
		timelineMilestoneColor(kind),
		timelineMilestoneBgColor(kind),
	)
}

// RenderFooterStatIcon returns a colored footer stat marker (○◉◈● or nerd equivalent).
func RenderFooterStatIcon(kind string) string {
	return renderMarkerGlyph(
		icons.FooterStatIcon(kind),
		footerStatColor(kind),
		footerStatBgColor(kind),
	)
}

// RenderLifecycleEvent returns a colored lifecycle event icon for TUI views.
func RenderLifecycleEvent(eventType string) string {
	fg, bg := lifecycleEventColors(eventType)
	return renderMarkerGlyph(icons.LifecycleEvent(eventType), fg, bg)
}

// RenderGraphIssueGlyphs returns status + priority + type icons for the dependency graph view.
// Each segment keeps its own lipgloss color — do not wrap the result in another Foreground style.
func RenderGraphIssueGlyphs(status string, priority int, issueType string, t Theme) string {
	return RenderStatusDotGraph(status) + " " +
		RenderPriorityIcon(priority) + " " +
		t.RenderTypeIcon(issueType)
}

// RenderStatusBadge returns a styled status badge
func RenderStatusBadge(status string) string {
	var fg, bg lipgloss.AdaptiveColor
	var label string

	switch status {
	case "open":
		fg, bg, label = ColorStatusOpen, ColorStatusOpenBg, "OPEN"
	case "in_progress":
		fg, bg, label = ColorStatusInProgress, ColorStatusInProgressBg, "PROG"
	case "blocked":
		fg, bg, label = ColorStatusBlocked, ColorStatusBlockedBg, "BLKD"
	case "deferred":
		fg, bg, label = ColorStatusDeferred, ColorStatusDeferredBg, "DEFR"
	case "draft":
		fg, bg, label = ColorStatusDeferred, ColorStatusDeferredBg, "DRFT"
	case "pinned":
		fg, bg, label = ColorStatusPinned, ColorStatusPinnedBg, "PIN"
	case "hooked":
		fg, bg, label = ColorStatusHooked, ColorStatusHookedBg, "HOOK"
	case "review":
		fg, bg, label = ColorStatusReview, ColorStatusReviewBg, "REVW"
	case "closed":
		fg, bg, label = ColorStatusClosed, ColorStatusClosedBg, "DONE"
	case "tombstone":
		fg, bg, label = ColorStatusTombstone, ColorStatusTombstoneBg, "TOMB"
	default:
		fg, bg, label = ColorMuted, ColorBgSubtle, "????"
	}

	return lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Padding(0, 0).
		Render(label)
}

// ══════════════════════════════════════════════════════════════════════════════
// METRIC VISUALIZATION - Mini-bars and rank badges
// ══════════════════════════════════════════════════════════════════════════════

// RenderMiniBar renders a mini horizontal bar for a value between 0 and 1
func RenderMiniBar(value float64, width int, t Theme) string {
	if width <= 0 {
		return ""
	}
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}

	filled := int(value * float64(width))
	if filled > width {
		filled = width
	}

	// Choose color based on value
	var barColor lipgloss.AdaptiveColor
	if value >= 0.75 {
		barColor = t.Open // Green/Success
	} else if value >= 0.5 {
		barColor = t.Feature // Orange/Warning
	} else if value >= 0.25 {
		barColor = t.InProgress // Cyan/Info
	} else {
		barColor = t.Secondary // Muted
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return t.Renderer.NewStyle().Foreground(barColor).Render(bar)
}

// RenderRankBadge renders a rank badge like "#1" with color based on percentile
func RenderRankBadge(rank, total int) string {
	if total == 0 {
		return lipgloss.NewStyle().Foreground(ColorMuted).Render("#?")
	}

	percentile := float64(rank) / float64(total)

	var color lipgloss.AdaptiveColor
	if percentile <= 0.1 {
		color = ColorSuccess // Top 10%
	} else if percentile <= 0.25 {
		color = ColorInfo // Top 25%
	} else if percentile <= 0.5 {
		color = ColorWarning // Top 50%
	} else {
		color = ColorMuted // Bottom 50%
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Render(fmt.Sprintf("#%d", rank))
}

// ══════════════════════════════════════════════════════════════════════════════
// DIVIDERS AND SEPARATORS
// ══════════════════════════════════════════════════════════════════════════════

// RenderDivider renders a horizontal divider line
func RenderDivider(width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Foreground(ColorBgHighlight).
		Render(strings.Repeat("─", width))
}

// RenderSubtleDivider renders a more subtle divider using dots
func RenderSubtleDivider(width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(strings.Repeat("·", width))
}
