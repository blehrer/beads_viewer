// Package icons centralizes UI glyphs for emoji and Nerd Font terminals.
//
// Set BV_ICON_SET=nerd (also nerdfont, nerd-font, nf) or experimental.icon_set: nerd in
// ~/.config/bv/config.yaml for Private Use Area icons. emoji or unset defaults to
// Unicode emoji for markdown export, robot JSON, and non-Nerd terminals. cmd/bv
// calls ApplyPreference at startup; library tests default to emoji.
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
	HistoryGit       Name = "history_git"
	HistoryBead      Name = "history_bead"
	SessionAttach    Name = "session_attach"
	ArrowUp          Name = "arrow_up"
	ArrowDown        Name = "arrow_down"
	Folder           Name = "folder"
	Globe            Name = "globe"
	Comment          Name = "comment"
	Label            Name = "label"
	Navigation       Name = "navigation"
	Lightbulb        Name = "lightbulb"
	Bell             Name = "bell"
	Microscope       Name = "microscope"
	Books            Name = "books"
	Brain            Name = "brain"
	Building         Name = "building"
	Satellite        Name = "satellite"
	CriticalPath     Name = "critical_path"
	Health           Name = "health"
	HistoryScroll    Name = "history_scroll"
	CommitRefactor   Name = "commit_refactor"
	CommitTest       Name = "commit_test"
	CommitChore      Name = "commit_chore"
	CommitStyle      Name = "commit_style"
	CommitRevert     Name = "commit_revert"
	Legend           Name = "legend"
	TrendUp          Name = "trend_up"
	Eye              Name = "eye"
)

var (
	mu         sync.RWMutex
	currentSet = SetEmoji
)

// Default is emoji. Call ApplyPreference from cmd/bv main before rendering.

// ParseSet maps a preference string to a Set. Returns (set, true) when s is
// non-empty and recognized (nerd aliases or emoji).
func ParseSet(s string) (Set, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "nerd", "nerdfont", "nerd-font", "nf":
		return SetNerd, true
	case "emoji":
		return SetEmoji, true
	default:
		return SetEmoji, false
	}
}

// ApplyPreference sets the active icon set with precedence: envVal > configVal > emoji.
func ApplyPreference(envVal, configVal string) {
	if set, ok := ParseSet(envVal); ok {
		Use(set)
		return
	}
	if set, ok := ParseSet(configVal); ok {
		Use(set)
		return
	}
	Use(SetEmoji)
}

// SetFromEnv reads BV_ICON_SET only. Safe to call again after changing the env var in tests.
func SetFromEnv() {
	ApplyPreference(os.Getenv("BV_ICON_SET"), "")
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

// IssueStatus returns export/robot/markdown glyphs for a beads status string.
// For TUI list views use ui.RenderStatusDot; for graph view use ui.RenderStatusDotGraph.
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

// IssueStatusGraph returns export-oriented status glyphs (emoji/nerd) for graph-style labeling.
// The interactive TUI graph uses ui.RenderStatusDotGraph (theme-colored ●/✓) instead.
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
	HistoryGit:       "◉",
	HistoryBead:      "◈",
	SessionAttach:    "📎",
	ArrowUp:          "⬆",
	ArrowDown:        "⬇",
	Folder:           "📁",
	Globe:            "🌐",
	Comment:          "💬",
	Label:            "🏷",
	Navigation:       "🧭",
	Lightbulb:        "💡",
	Bell:             "🔔",
	Microscope:       "🔬",
	Books:            "📚",
	Brain:            "🧠",
	Building:         "🏛️",
	Satellite:        "🛰️",
	CriticalPath:     "🛤️",
	Health:           "🩺",
	HistoryScroll:    "📜",
	CommitRefactor:   "♻",
	CommitTest:       "🧪",
	CommitChore:      "🔧",
	CommitStyle:      "💄",
	CommitRevert:     "↩",
	Legend:           "📖",
	TrendUp:          "📈",
	Eye:              "👁",
}

// ponytail: codepoints from Nerd Fonts 3.4 i_md.sh — not raw MDI webfont names (F0AFA is alpha-m, not bug).
var nerdIcons = map[Name]string{
	Bug:              "\U000f00e4", // md-bug
	Feature:          "\U000f1545", // md-shimmer
	Task:             "\U000f014d", // md-clipboard-text
	Epic:             "\U000f14de", // md-rocket-launch
	Chore:            "\U000f00e2", // md-broom
	Target:           "\U000f04fe", // md-target
	Unlock:           "\U000f0fc6", // md-lock-open-variant
	Warning:          "\U000f0026", // md-alert
	Shuffle:          "\U000f049d", // md-shuffle
	Chart:            "\U000f0128", // md-chart-bar
	Clock:            "\U000f0150", // md-clock-outline
	Calendar:         "\U000f00ed", // md-calendar
	Lightning:        "\U000f0241", // md-flash
	CheckCircle:      "\U000f05e0", // md-check-circle
	Construction:     "\U000f0462", // md-road-variant
	Blocked:          "\U000f073a", // md-cancel
	Pause:            "\U000f03e4", // md-pause
	User:             "\U000f0004", // md-account
	Hourglass:        "\U000f051f", // md-timer-sand
	Siren:            "\U000f078f", // md-alarm-light
	Star:             "\U000f04ce", // md-star
	Fire:             "\U000f0238", // md-fire
	Alarm:            "\U000f0020", // md-alarm
	Link:             "\U000f0337", // md-link
	New:              "\U000f0394", // nf-md-new_box
	Check:            "\U000f012c", // md-check
	Cross:            "\U000f0159", // md-close-circle
	StatusOpen:       "\U000f0ec2", // md-record-circle
	StatusInProgress: "\U000f09de", // md-circle-medium
	StatusBlocked:    "\U000f073a", // md-cancel
	StatusClosed:     "\U000f0159", // md-close-circle
	StatusUnknown:    "\U000f0ec3", // md-record-circle-outline
	StatusDeferred:   "\U000f03e4", // md-pause
	StatusPinned:     "\U000f0403", // md-pin
	StatusHooked:     "\U000f06e2", // md-hook
	StatusReview:     "\U000f0208", // md-eye
	StatusGraphOpen:  "\U000f09de", // md-circle-medium
	StatusGraphWork:  "\U000f0996", // md-progress-clock
	TriageScoreMid:   "\U000f0a10", // md-rhombus-medium
	PriorityMedium:   "\U000f0a10", // md-rhombus-medium
	PriorityLow:      "\U000f0176", // md-coffee
	PriorityBacklog:  "\U000f04b2", // md-sleep
	DepRoot:          "\U000f034e", // md-map-marker
	DepParentChild:   "\U000f03d6", // md-package-variant
	DepDiscovered:    "\U000f0349", // md-magnify
	SwimRefresh:      "\U000f04e6", // md-sync
	SwimProhibited:   "\U000f073a", // md-cancel
	Question:         "\U000f02d7", // md-help-circle
	FileDefault:      "\U000f0219", // md-file-document
	HistoryGit:       "\U000f02a2", // md-git
	HistoryBead:      "\U000f04fe", // md-target
	SessionAttach:    "\U000f0207", // md-paperclip
	ArrowUp:          "\U000f005d", // md-arrow-up
	ArrowDown:        "\U000f0045", // md-arrow-down
	Folder:           "\U000f024b", // nf-md-folder
	Globe:            "\U000f01e7", // nf-md-web
	Comment:          "\U000f017a", // nf-md-comment-text-outline
	Label:            "\U000f0315", // nf-md-label
	Navigation:       "\U000f0142", // nf-md-compass-outline
	Lightbulb:        "\U000f0335", // nf-md-lightbulb
	Bell:             "\U000f009a", // nf-md-bell
	Microscope:       "\U000f0650", // nf-md-microscope
	Books:            "\U000f00f6", // nf-md-book-multiple
	Brain:            "\U000f09e7", // nf-md-brain
	Building:         "\U000f01ad", // nf-md-office-building
	Satellite:        "\U000f0420", // nf-md-satellite-variant
	CriticalPath:     "\U000f0462", // nf-md-road-variant
	Health:           "\U000f0fa4", // nf-md-stethoscope
	HistoryScroll:    "\U000f018f", // nf-md-history
	CommitRefactor:   "\U000f012a", // nf-md-recycle
	CommitTest:       "\U000f0669", // nf-md-test-tube
	CommitChore:      "\U000f05b7", // nf-md-wrench
	CommitStyle:      "\U000f0154", // nf-md-brush
	CommitRevert:     "\U000f045a", // nf-md-undo
	Legend:           "\U000f02f7", // nf-md-book-open-page-variant
	TrendUp:          "\U000f053a", // nf-md-chart-line
	Eye:              "\U000f0208", // nf-md-eye
}
