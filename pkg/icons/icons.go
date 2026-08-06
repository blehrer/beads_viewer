// Package icons centralizes UI glyphs for emoji and Nerd Font terminals.
//
// Set BV_ICON_SET=nerd (also nerdfont, nerd-font, nf) for Private Use Area icons.
// emoji or unset defaults to Unicode emoji for markdown export, robot JSON, and
// non-Nerd terminals. cmd/bv calls SetFromEnv at startup; library tests default to emoji.
package icons

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Set selects which glyph table Get returns.
type Set int

const (
	SetEmoji Set = iota
	SetNerd
)

// Name identifies a logical icon independent of rendering backend.
type Name string

const (
	Bug              Name = "bug"
	Feature          Name = "feature"
	Task             Name = "task"
	Epic             Name = "epic"
	Chore            Name = "chore"
	Target           Name = "target"
	Unlock           Name = "unlock"
	Warning          Name = "warning"
	Shuffle          Name = "shuffle"
	Chart            Name = "chart"
	Clock            Name = "clock"
	Calendar         Name = "calendar"
	Lightning        Name = "lightning"
	CheckCircle      Name = "check_circle"
	Construction     Name = "construction"
	Blocked          Name = "blocked"
	Pause            Name = "pause"
	User             Name = "user"
	Hourglass        Name = "hourglass"
	Siren            Name = "siren"
	Star             Name = "star"
	Fire             Name = "fire"
	Alarm            Name = "alarm"
	Link             Name = "link"
	New              Name = "new"
	Check            Name = "check"
	Cross            Name = "cross"
	StatusOpen       Name = "status_open"
	StatusInProgress Name = "status_in_progress"
	StatusBlocked    Name = "status_blocked"
	StatusClosed     Name = "status_closed"
	StatusUnknown    Name = "status_unknown"
	StatusDeferred   Name = "status_deferred"
	StatusPinned     Name = "status_pinned"
	StatusHooked     Name = "status_hooked"
	StatusReview     Name = "status_review"
	StatusGraphOpen  Name = "status_graph_open"
	StatusGraphWork  Name = "status_graph_work"
	TriageScoreMid   Name = "triage_score_mid"
	PriorityMedium   Name = "priority_medium"
	PriorityLow      Name = "priority_low"
	PriorityBacklog  Name = "priority_backlog"
	DepRoot          Name = "dep_root"
	DepParentChild   Name = "dep_parent_child"
	DepDiscovered    Name = "dep_discovered"
	SwimRefresh      Name = "swim_refresh"
	SwimProhibited   Name = "swim_prohibited"
	Question         Name = "question"
	FileDefault      Name = "file_default"
)

var (
	mu         sync.RWMutex
	currentSet = SetEmoji
)

// Default is emoji. Call SetFromEnv from cmd/bv main before rendering.

// SetFromEnv reads BV_ICON_SET. Safe to call again after changing the env var in tests.
func SetFromEnv() {
	mu.Lock()
	defer mu.Unlock()
	switch strings.ToLower(strings.TrimSpace(os.Getenv("BV_ICON_SET"))) {
	case "nerd", "nerdfont", "nerd-font", "nf":
		currentSet = SetNerd
	default:
		currentSet = SetEmoji
	}
}

// ActiveSet returns the active icon set.
func ActiveSet() Set {
	mu.RLock()
	defer mu.RUnlock()
	return currentSet
}

// Use switches the active set programmatically (tests and early startup overrides).
func Use(set Set) {
	mu.Lock()
	defer mu.Unlock()
	currentSet = set
}

// Get returns the glyph for name in the active set, falling back to emoji then "?".
func Get(name Name) string {
	mu.RLock()
	set := currentSet
	mu.RUnlock()
	if set == SetNerd {
		if s, ok := nerdIcons[name]; ok {
			return s
		}
	}
	if s, ok := emojiIcons[name]; ok {
		return s
	}
	return "?"
}

// IssueType returns the icon for a beads issue type string.
func IssueType(typ string) string {
	switch typ {
	case "bug":
		return Get(Bug)
	case "feature":
		return Get(Feature)
	case "task":
		return Get(Task)
	case "epic":
		return Get(Epic)
	case "chore":
		return Get(Chore)
	default:
		return "•"
	}
}

// IssueStatus returns the icon for a beads status string.
func IssueStatus(status string) string {
	switch status {
	case "open":
		return Get(StatusOpen)
	case "in_progress":
		return Get(StatusInProgress)
	case "blocked":
		return Get(StatusBlocked)
	case "closed", "tombstone":
		return Get(StatusClosed)
	case "deferred", "draft":
		return Get(Pause)
	case "pinned":
		return Get(StatusPinned)
	case "hooked":
		return Get(StatusHooked)
	case "review":
		return Get(StatusReview)
	default:
		return Get(StatusUnknown)
	}
}

// IssueStatusGraph returns status glyphs for the interactive graph view.
// ponytail: graph uses blue/yellow circles and a checkmark for done — not the list-view circles.
func IssueStatusGraph(status string) string {
	switch status {
	case "closed", "tombstone":
		return Get(CheckCircle)
	case "open":
		return Get(StatusGraphOpen)
	case "in_progress":
		return Get(StatusGraphWork)
	case "blocked":
		return Get(StatusBlocked)
	case "deferred", "draft":
		return Get(Pause)
	case "pinned":
		return Get(StatusPinned)
	case "hooked":
		return Get(StatusHooked)
	case "review":
		return Get(StatusReview)
	default:
		return Get(StatusUnknown)
	}
}

// Priority returns the icon for beads priority 0–4 (P0 critical … P4 backlog).
func Priority(level int) string {
	switch level {
	case 0:
		return Get(Fire)
	case 1:
		return Get(Lightning)
	case 2:
		return Get(PriorityMedium)
	case 3:
		return Get(PriorityLow)
	case 4:
		return Get(PriorityBacklog)
	default:
		return "  "
	}
}

// DependencyType returns the icon for a dependency edge type string.
func DependencyType(depType string) string {
	switch depType {
	case "root":
		return Get(DepRoot)
	case "blocks":
		return Get(Blocked)
	case "related":
		return Get(Link)
	case "parent-child":
		return Get(DepParentChild)
	case "discovered-from":
		return Get(DepDiscovered)
	default:
		return "•"
	}
}

// PriorityLabel returns a markdown-friendly priority line (emoji + text).
func PriorityLabel(level int) string {
	switch level {
	case 0:
		return Get(Fire) + " Critical (P0)"
	case 1:
		return Get(Lightning) + " High (P1)"
	case 2:
		return Get(PriorityMedium) + " Medium (P2)"
	case 3:
		return Get(PriorityLow) + " Low (P3)"
	case 4:
		return Get(PriorityBacklog) + " Backlog (P4)"
	default:
		return fmt.Sprintf("P%d", level)
	}
}

var emojiIcons = map[Name]string{
	Bug:              "🐛",
	Feature:          "✨",
	Task:             "📋",
	Epic:             "🚀", // ponytail: avoid U+FE0F variation selector — breaks terminal width
	Chore:            "🧹",
	Target:           "🎯",
	Unlock:           "🔓",
	Warning:          "⚠️",
	Shuffle:          "🔀",
	Chart:            "📊",
	Clock:            "🕐",
	Calendar:         "📅",
	Lightning:        "⚡",
	CheckCircle:      "✅",
	Construction:     "🚧",
	Blocked:          "⛔",
	Pause:            "⏸️",
	User:             "👤",
	Hourglass:        "⏳",
	Siren:            "🚨",
	Star:             "⭐",
	Fire:             "🔥",
	Alarm:            "⏰",
	Link:             "🔗",
	New:              "🆕",
	Check:            "✓",
	Cross:            "❌",
	StatusOpen:       "🟢",
	StatusInProgress: "🔵",
	StatusBlocked:    "🔴",
	StatusClosed:     "⚫",
	StatusUnknown:    "⚪",
	StatusDeferred:   "⏸️",
	StatusPinned:     "📌",
	StatusHooked:     "🪝",
	StatusReview:     "👁️",
	StatusGraphOpen:  "🔵",
	StatusGraphWork:  "🟡",
	TriageScoreMid:   "🟠",
	PriorityMedium:   "🔹",
	PriorityLow:      "☕",
	PriorityBacklog:  "💤",
	DepRoot:          "📍",
	DepParentChild:   "📦",
	DepDiscovered:    "🔍",
	SwimRefresh:      "🔄",
	SwimProhibited:   "🚫",
	Question:         "❓",
	FileDefault:      "📄",
}

// ponytail: NF codepoints are Material Design Icons from the Nerd Fonts 3.x PUA block.
var nerdIcons = map[Name]string{
	Bug:              "\U000f0afa", // md-bug
	Feature:          "\U000f04a4", // md-shimmer
	Task:             "\U000f014e", // md-clipboard-text
	Epic:             "\U000f14de", // md-rocket-launch
	Chore:            "\U000f0823", // md-broom
	Target:           "\U000f0f44", // md-target
	Unlock:           "\U000f08c5", // md-lock-open-variant
	Warning:          "\U000f0026", // md-alert
	Shuffle:          "\U000f0430", // md-shuffle
	Chart:            "\U000f0126", // md-chart-bar
	Clock:            "\U000f0189", // md-clock-outline
	Calendar:         "\U000f00ed", // md-calendar
	Lightning:        "\U000f024b", // md-flash
	CheckCircle:      "\U000f0133", // md-check-circle
	Construction:     "\U000f0461", // md-road-variant
	Blocked:          "\U000f073a", // md-cancel
	Pause:            "\U000f03e4", // md-pause
	User:             "\U000f0004", // md-account
	Hourglass:        "\U000f05ad", // md-timer-sand
	Siren:            "\U000f0078", // md-alarm-light
	Star:             "\U000f04ce", // md-star
	Fire:             "\U000f0238", // md-fire
	Alarm:            "\U000f0120", // md-alarm
	Link:             "\U000f0337", // md-link
	New:              "\U000f0054", // md-new-box
	Check:            "\U000f012c", // md-check
	Cross:            "\U000f0159", // md-close-circle
	StatusOpen:       "\U000f0765", // md-circle
	StatusInProgress: "\U000f0130", // md-progress-clock
	StatusBlocked:    "\U000f073a", // md-cancel
	StatusClosed:     "\U000f0133", // md-check-circle
	StatusUnknown:    "\U000f0766", // md-circle-outline
	StatusDeferred:   "\U000f03e4", // md-pause
	StatusPinned:     "\U000f0403", // md-pin
	StatusHooked:     "\U000f06e2", // md-hook
	StatusReview:     "\U000f0208", // md-eye
	StatusGraphOpen:  "\U000f0765", // md-circle
	StatusGraphWork:  "\U000f0130", // md-progress-clock
	TriageScoreMid:   "\U000f0425", // md-circle-medium
	PriorityMedium:   "\U000f04a0", // md-rhombus-medium
	PriorityLow:      "\U000f0176", // md-coffee
	PriorityBacklog:  "\U000f03d8", // md-sleep
	DepRoot:          "\U000f034e", // md-map-marker
	DepParentChild:   "\U000f0493", // md-package-variant
	DepDiscovered:    "\U000f0349", // md-magnify
	SwimRefresh:      "\U000f0450", // md-sync
	SwimProhibited:   "\U000f073a", // md-cancel
	Question:         "\U000f02d7", // md-help-circle
	FileDefault:      "\U000f0214", // md-file-document
}
