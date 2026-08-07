package ui

import (
	"os"
	"strings"

	"github.com/Dicklesworthstone/beads_viewer/pkg/icons"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss"
)

// TermProfile holds the detected terminal color profile. Computed once at
// package init so every style helper can branch without re-detecting.
var TermProfile colorprofile.Profile

// defaultRenderer backs package-level Render* helpers in styles.go.
// DefaultTheme assigns the live stdout renderer; tests may override.
var defaultRenderer = lipgloss.NewRenderer(os.Stdout)

// BVThemeOverride holds the user's explicit theme preference.
// Values: "" (auto-detect), "dark", "light". It is seeded from BV_THEME at
// package init for backward compatibility, but the full precedence
// (--theme flag > BV_THEME > config file) is resolved by cmd/bv, which calls
// SetThemeOverride once at startup. Always mutate it through SetThemeOverride
// so the global lipgloss renderer stays in agreement.
var BVThemeOverride string

func init() {
	TermProfile = colorprofile.Detect(os.Stdout, os.Environ())

	// BV_THEME allows users to override auto-detection of light/dark background.
	// This is useful when the terminal reports the wrong background color or
	// when using themes that confuse auto-detection (e.g., Windows Terminal
	// custom schemes, tmux, SSH). (bv-128)
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("BV_THEME"))); v == "light" || v == "dark" {
		BVThemeOverride = v
	}
}

// SetThemeOverride applies an explicit light/dark theme preference
// process-wide. "light" and "dark" record the preference in BVThemeOverride
// AND pin the global (default-renderer) background assumption via
// lipgloss.SetHasDarkBackground. Any other value ("auto", "", unrecognized)
// clears the override and leaves lipgloss's auto-detection in charge.
//
// Pinning the GLOBAL renderer is the load-bearing part: most bv styles are
// built with the package-level lipgloss.NewStyle(), and markdown rendering
// branches on the global lipgloss.HasDarkBackground(), so overriding only a
// per-model renderer (as DefaultTheme does) leaves the majority of adaptive
// colors on auto-detection. Auto-detection assumes a DARK background whenever
// the terminal never answers the background query — typical over SSH and
// inside tmux/screen — which renders near-white text on light terminals.
//
// Call this once, early at startup, before any styles are rendered. (bv-128)
func SetThemeOverride(pref string) {
	switch strings.ToLower(strings.TrimSpace(pref)) {
	case "light":
		BVThemeOverride = "light"
		lipgloss.SetHasDarkBackground(false)
	case "dark":
		BVThemeOverride = "dark"
		lipgloss.SetHasDarkBackground(true)
	default:
		// "auto" / empty / unknown: no pin; adaptive colors follow the
		// terminal's detected background.
		BVThemeOverride = ""
	}
}

// ThemeBg returns the given hex color for TrueColor terminals and
// lipgloss.NoColor{} otherwise, so 16/256-color terminals use the
// terminal's own background instead of a down-converted approximation
// that may clash with palettes like Solarized.
func ThemeBg(hex string) lipgloss.TerminalColor {
	if TermProfile < colorprofile.TrueColor {
		return lipgloss.NoColor{}
	}
	return lipgloss.Color(hex)
}

// ThemeFg returns the given hex color for ANSI256+ terminals and a safe
// ANSI white (color 7) for 16-color or lower terminals.
func ThemeFg(hex string) lipgloss.TerminalColor {
	if TermProfile < colorprofile.ANSI256 {
		return lipgloss.ANSIColor(7)
	}
	return lipgloss.Color(hex)
}

type Theme struct {
	Renderer *lipgloss.Renderer

	// Colors
	Primary         lipgloss.AdaptiveColor
	Secondary       lipgloss.AdaptiveColor
	Subtext         lipgloss.AdaptiveColor
	Text            lipgloss.AdaptiveColor
	HeaderOnPrimary lipgloss.AdaptiveColor

	// Status
	Open       lipgloss.AdaptiveColor
	InProgress lipgloss.AdaptiveColor
	Blocked    lipgloss.AdaptiveColor
	Deferred   lipgloss.AdaptiveColor
	Pinned     lipgloss.AdaptiveColor
	Hooked     lipgloss.AdaptiveColor
	Review     lipgloss.AdaptiveColor
	Closed     lipgloss.AdaptiveColor
	Tombstone  lipgloss.AdaptiveColor

	// Types
	Bug     lipgloss.AdaptiveColor
	Feature lipgloss.AdaptiveColor
	Task    lipgloss.AdaptiveColor
	Epic    lipgloss.AdaptiveColor
	Chore   lipgloss.AdaptiveColor

	// UI Elements
	Border    lipgloss.AdaptiveColor
	Highlight lipgloss.AdaptiveColor
	Muted     lipgloss.AdaptiveColor

	// Styles
	Base     lipgloss.Style
	Selected lipgloss.Style
	Column   lipgloss.Style
	Header   lipgloss.Style

	// Pre-computed delegate styles (bv-o4cj optimization)
	// These are created once at startup instead of per-frame
	MutedText         lipgloss.Style // Age, muted info
	InfoText          lipgloss.Style // Comments
	InfoBold          lipgloss.Style // Search scores
	SecondaryText     lipgloss.Style // ID, assignee
	PrimaryBold       lipgloss.Style // Selection indicator
	PriorityUpArrow   lipgloss.Style // Priority hint ↑
	PriorityDownArrow lipgloss.Style // Priority hint ↓
	TriageStar        lipgloss.Style // Top pick ⭐
	TriageUnblocks    lipgloss.Style // Unblocks indicator 🔓
	TriageUnblocksAlt lipgloss.Style // Secondary unblocks ↪
}

// DefaultTheme returns the standard Dracula-inspired theme (adaptive).
// Respects BV_THEME=light|dark to override background detection. (bv-128)
func DefaultTheme(r *lipgloss.Renderer) Theme {
	if r != nil {
		defaultRenderer = r
	}
	// Apply BV_THEME override so AdaptiveColor picks the right variant
	if r != nil && BVThemeOverride != "" {
		r.SetHasDarkBackground(BVThemeOverride == "dark")
	}
	p := ActivePalette()
	t := Theme{
		Renderer: r,

		Primary:         p.Primary,
		Secondary:       p.Secondary,
		Subtext:         p.Subtext,
		Text:            p.Text,
		HeaderOnPrimary: p.HeaderOnPrimary,

		Open:       p.StatusOpen,
		InProgress: p.StatusInProgress,
		Blocked:    p.StatusBlocked,
		Deferred:   p.StatusDeferred,
		Pinned:     p.StatusPinned,
		Hooked:     p.StatusHooked,
		Review:     p.StatusReview,
		Closed:     p.StatusClosed,
		Tombstone:  p.StatusTombstone,

		Bug:     p.TypeBug,
		Feature: p.TypeFeature,
		Epic:    p.TypeEpic,
		Task:    p.TypeTask,
		Chore:   p.TypeChore,

		Border:    p.Border,
		Highlight: p.Highlight,
		Muted:     p.Muted,
	}

	t.Base = r.NewStyle().Foreground(p.TextBase)

	t.Selected = r.NewStyle().
		Background(t.Highlight).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Primary).
		PaddingLeft(1).
		Bold(true)

	t.Header = r.NewStyle().
		Background(t.Primary).
		Foreground(p.HeaderOnPrimary).
		Bold(true).
		Padding(0, 1)

	// Pre-computed delegate styles (bv-o4cj optimization)
	// Reduces ~16 NewStyle() allocations per visible item per frame
	t.MutedText = r.NewStyle().Foreground(p.Muted)
	t.InfoText = r.NewStyle().Foreground(p.Info)
	t.InfoBold = r.NewStyle().Foreground(p.Info).Bold(true)
	t.SecondaryText = r.NewStyle().Foreground(t.Secondary)
	t.PrimaryBold = r.NewStyle().Foreground(t.Primary).Bold(true)
	t.PriorityUpArrow = r.NewStyle().Foreground(ThemeFg(p.TriagePriorityUp)).Bold(true)
	t.PriorityDownArrow = r.NewStyle().Foreground(ThemeFg(p.TriagePriorityDown)).Bold(true)
	t.TriageStar = r.NewStyle().Foreground(ThemeFg(p.TriageStar))
	t.TriageUnblocks = r.NewStyle().Foreground(ThemeFg(p.TriageUnblocks))
	t.TriageUnblocksAlt = r.NewStyle().Foreground(ThemeFg(p.TriageUnblocksAlt))

	return t
}

func (t Theme) GetStatusColor(s string) lipgloss.AdaptiveColor {
	switch s {
	case "open":
		return t.Open
	case "in_progress":
		return t.InProgress
	case "blocked":
		return t.Blocked
	case "deferred", "draft":
		return t.Deferred
	case "pinned":
		return t.Pinned
	case "hooked":
		return t.Hooked
	case "review":
		return t.Review
	case "closed":
		return t.Closed
	case "tombstone":
		return t.Tombstone
	default:
		return t.Subtext
	}
}

func (t Theme) GetTypeIcon(typ string) (string, lipgloss.AdaptiveColor) {
	switch typ {
	case "bug":
		return icons.IssueType(typ), t.Bug
	case "feature":
		return icons.IssueType(typ), t.Feature
	case "task":
		return icons.IssueType(typ), t.Task
	case "epic":
		return icons.IssueType(typ), t.Epic
	case "chore":
		return icons.IssueType(typ), t.Chore
	default:
		return "•", t.Subtext
	}
}

// RenderTypeIcon returns the issue-type glyph with theme foreground color applied.
func (t Theme) RenderTypeIcon(typ string) string {
	icon, color := t.GetTypeIcon(typ)
	return t.Renderer.NewStyle().Foreground(color).Render(icon)
}
