// Package icons centralizes UI glyphs for emoji and Nerd Font terminals.
//
// Set BV_ICON_SET=nerd (also nerdfont, nerd-font) for Private Use Area icons.
// Default is emoji for markdown export, robot JSON, and non-Nerd terminals.
package icons

import (
	"os"
	"strings"
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
	Bug          Name = "bug"
	Feature      Name = "feature"
	Task         Name = "task"
	Epic         Name = "epic"
	Chore        Name = "chore"
	Target       Name = "target"
	Unlock       Name = "unlock"
	Warning      Name = "warning"
	Shuffle      Name = "shuffle"
	Chart        Name = "chart"
	Clock        Name = "clock"
	Calendar     Name = "calendar"
	Lightning    Name = "lightning"
	CheckCircle  Name = "check_circle"
	Construction Name = "construction"
	Blocked      Name = "blocked"
	Pause        Name = "pause"
	User         Name = "user"
	Hourglass    Name = "hourglass"
	Siren        Name = "siren"
	Star         Name = "star"
	Fire         Name = "fire"
	Alarm        Name = "alarm"
	Link         Name = "link"
	New          Name = "new"
	Check        Name = "check"
	Cross        Name = "cross"
)

var currentSet = SetEmoji

func init() {
	SetFromEnv()
}

// SetFromEnv reads BV_ICON_SET. Safe to call again after changing the env var in tests.
func SetFromEnv() {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("BV_ICON_SET"))) {
	case "nerd", "nerdfont", "nerd-font", "nf":
		currentSet = SetNerd
	default:
		currentSet = SetEmoji
	}
}

// ActiveSet returns the active icon set.
func ActiveSet() Set { return currentSet }

// Use switches the active set programmatically (tests and early startup overrides).
func Use(set Set) { currentSet = set }

// Get returns the glyph for name in the active set, falling back to emoji.
func Get(name Name) string {
	if currentSet == SetNerd {
		if s, ok := nerdIcons[name]; ok {
			return s
		}
	}
	return emojiIcons[name]
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
		return statusGlyph("🟢", "\U000f0765") // md-circle
	case "in_progress":
		return statusGlyph("🔵", "\U000f0765")
	case "blocked":
		return statusGlyph("🔴", "\U000f0765")
	case "closed", "tombstone":
		return statusGlyph("⚫", "\U000f0765")
	default:
		return statusGlyph("⚪", "\U000f0766") // md-circle-outline
	}
}

func statusGlyph(emoji, nerd string) string {
	if currentSet == SetNerd {
		return nerd
	}
	return emoji
}

var emojiIcons = map[Name]string{
	Bug:          "🐛",
	Feature:      "✨",
	Task:         "📋",
	Epic:         "🚀",
	Chore:        "🧹",
	Target:       "🎯",
	Unlock:       "🔓",
	Warning:      "⚠️",
	Shuffle:      "🔀",
	Chart:        "📊",
	Clock:        "🕐",
	Calendar:     "📅",
	Lightning:    "⚡",
	CheckCircle:  "✅",
	Construction: "🚧",
	Blocked:      "⛔",
	Pause:        "⏸️",
	User:         "👤",
	Hourglass:    "⏳",
	Siren:        "🚨",
	Star:         "⭐",
	Fire:         "🔥",
	Alarm:        "⏰",
	Link:         "🔗",
	New:          "🆕",
	Check:        "✓",
	Cross:        "❌",
}

// ponytail: NF codepoints are Material Design Icons from the Nerd Fonts 3.x PUA block.
var nerdIcons = map[Name]string{
	Bug:          "\U000f0afa", // md-bug
	Feature:      "\U000f04a4", // md-shimmer
	Task:         "\U000f014e", // md-clipboard-text
	Epic:         "\U000f14de", // md-rocket-launch
	Chore:        "\U000f0823", // md-broom
	Target:       "\U000f0f44", // md-target
	Unlock:       "\U000f08c5", // md-lock-open-variant
	Warning:      "\U000f0026", // md-alert
	Shuffle:      "\U000f0430", // md-shuffle
	Chart:        "\U000f0126", // md-chart-bar
	Clock:        "\U000f0189", // md-clock-outline
	Calendar:     "\U000f00ed", // md-calendar
	Lightning:    "\U000f024b", // md-flash
	CheckCircle:  "\U000f0133", // md-check-circle
	Construction: "\U000f0461", // md-road-variant
	Blocked:      "\U000f073a", // md-cancel
	Pause:        "\U000f03e4", // md-pause
	User:         "\U000f0004", // md-account
	Hourglass:    "\U000f05ad", // md-timer-sand
	Siren:        "\U000f0078", // md-alarm-light
	Star:         "\U000f04ce", // md-star
	Fire:         "\U000f0238", // md-fire
	Alarm:        "\U000f0120", // md-alarm
	Link:         "\U000f0337", // md-link
	New:          "\U000f0054", // md-new-box
	Check:        "\U000f012c", // md-check
	Cross:        "\U000f0159", // md-close-circle
}
