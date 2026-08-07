package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// AdaptivePair holds optional light/dark hex overrides for one palette slot.
type AdaptivePair struct {
	Light string `yaml:"light"`
	Dark  string `yaml:"dark"`
}

func (p AdaptivePair) merge(base lipgloss.AdaptiveColor) lipgloss.AdaptiveColor {
	out := base
	if s := strings.TrimSpace(p.Light); s != "" {
		out.Light = s
	}
	if s := strings.TrimSpace(p.Dark); s != "" {
		out.Dark = s
	}
	return out
}

// Palette is the single source of truth for TUI adaptive colors.
// Override keys for ~/.config/bv/config.yaml (under colors:) use snake_case
// names matching the struct fields, e.g. primary, status_open, type_bug.
type Palette struct {
	Bg          lipgloss.AdaptiveColor
	BgDark      lipgloss.AdaptiveColor
	BgSubtle    lipgloss.AdaptiveColor
	BgHighlight lipgloss.AdaptiveColor
	Text        lipgloss.AdaptiveColor
	Subtext     lipgloss.AdaptiveColor
	Muted       lipgloss.AdaptiveColor
	TextBase    lipgloss.AdaptiveColor

	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Info      lipgloss.AdaptiveColor
	Success   lipgloss.AdaptiveColor
	Warning   lipgloss.AdaptiveColor
	Danger    lipgloss.AdaptiveColor

	Border          lipgloss.AdaptiveColor
	Highlight       lipgloss.AdaptiveColor
	HeaderOnPrimary lipgloss.AdaptiveColor

	StatusOpen       lipgloss.AdaptiveColor
	StatusInProgress lipgloss.AdaptiveColor
	StatusBlocked    lipgloss.AdaptiveColor
	StatusDeferred   lipgloss.AdaptiveColor
	StatusPinned     lipgloss.AdaptiveColor
	StatusHooked     lipgloss.AdaptiveColor
	StatusReview     lipgloss.AdaptiveColor
	StatusClosed     lipgloss.AdaptiveColor
	StatusTombstone  lipgloss.AdaptiveColor

	StatusOpenBg       lipgloss.AdaptiveColor
	StatusInProgressBg lipgloss.AdaptiveColor
	StatusBlockedBg    lipgloss.AdaptiveColor
	StatusDeferredBg   lipgloss.AdaptiveColor
	StatusPinnedBg     lipgloss.AdaptiveColor
	StatusHookedBg     lipgloss.AdaptiveColor
	StatusReviewBg     lipgloss.AdaptiveColor
	StatusClosedBg     lipgloss.AdaptiveColor
	StatusTombstoneBg  lipgloss.AdaptiveColor

	PrioCritical lipgloss.AdaptiveColor
	PrioHigh     lipgloss.AdaptiveColor
	PrioMedium   lipgloss.AdaptiveColor
	PrioLow      lipgloss.AdaptiveColor

	PrioCriticalBg lipgloss.AdaptiveColor
	PrioHighBg     lipgloss.AdaptiveColor
	PrioMediumBg   lipgloss.AdaptiveColor
	PrioLowBg      lipgloss.AdaptiveColor

	TypeBug     lipgloss.AdaptiveColor
	TypeFeature lipgloss.AdaptiveColor
	TypeTask    lipgloss.AdaptiveColor
	TypeEpic    lipgloss.AdaptiveColor
	TypeChore   lipgloss.AdaptiveColor

	FooterHint lipgloss.AdaptiveColor
	FooterKey  lipgloss.AdaptiveColor
	FooterSep  lipgloss.AdaptiveColor
	FooterDim  lipgloss.AdaptiveColor

	// Accent hex strings for ThemeFg (triage hints, arrows).
	TriagePriorityUp   string
	TriagePriorityDown string
	TriageStar         string
	TriageUnblocks     string
	TriageUnblocksAlt  string
}

// PaletteOverrides maps config keys to partial light/dark overrides.
type PaletteOverrides map[string]AdaptivePair

var activePalette = draculaPalette()
var activePaletteName = "dracula"

// ActivePaletteName returns the resolved built-in palette name.
func ActivePaletteName() string {
	return activePaletteName
}

// ActivePalette returns the process-wide palette (built-in plus any config overrides).
func ActivePalette() Palette {
	return activePalette
}

// InitPalette loads a built-in palette, applies optional color overrides, and
// refreshes package-level Color* globals. Call once at startup.
func InitPalette(name string, overrides PaletteOverrides) {
	canonical := CanonicalPaletteName(name)
	if canonical == "" {
		canonical = "dracula"
	}
	p := PaletteByName(canonical)
	if len(overrides) > 0 {
		p = overrides.Apply(p)
	}
	activePalette = p
	activePaletteName = canonical
	syncStyleGlobalsFromPalette(activePalette)
}

// DefaultPalette returns the built-in Dracula-inspired adaptive palette.
func DefaultPalette() Palette {
	return draculaPalette()
}

// Apply merges partial overrides onto base. Unknown keys are ignored.
func (o PaletteOverrides) Apply(base Palette) Palette {
	p := base
	for key, pair := range o {
		if !applyPaletteOverride(&p, key, pair) {
			continue
		}
	}
	return p
}

var validPaletteKeys map[string]struct{}

func init() {
	syncStyleGlobalsFromPalette(activePalette)
	validPaletteKeys = make(map[string]struct{}, len(PaletteKeyNames()))
	for _, key := range PaletteKeyNames() {
		validPaletteKeys[key] = struct{}{}
	}
}

// UnknownKeys returns config keys that did not match any palette slot.
func (o PaletteOverrides) UnknownKeys() []string {
	var unknown []string
	for key := range o {
		if _, ok := validPaletteKeys[strings.ToLower(strings.TrimSpace(key))]; !ok {
			unknown = append(unknown, key)
		}
	}
	return unknown
}

func applyPaletteOverride(p *Palette, key string, pair AdaptivePair) bool {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "bg":
		p.Bg = pair.merge(p.Bg)
	case "bg_dark":
		p.BgDark = pair.merge(p.BgDark)
	case "bg_subtle":
		p.BgSubtle = pair.merge(p.BgSubtle)
	case "bg_highlight":
		p.BgHighlight = pair.merge(p.BgHighlight)
	case "text":
		p.Text = pair.merge(p.Text)
	case "subtext":
		p.Subtext = pair.merge(p.Subtext)
	case "muted":
		p.Muted = pair.merge(p.Muted)
	case "text_base":
		p.TextBase = pair.merge(p.TextBase)
	case "primary":
		p.Primary = pair.merge(p.Primary)
	case "secondary":
		p.Secondary = pair.merge(p.Secondary)
	case "info":
		p.Info = pair.merge(p.Info)
	case "success":
		p.Success = pair.merge(p.Success)
	case "warning":
		p.Warning = pair.merge(p.Warning)
	case "danger":
		p.Danger = pair.merge(p.Danger)
	case "border":
		p.Border = pair.merge(p.Border)
	case "highlight":
		p.Highlight = pair.merge(p.Highlight)
	case "header_on_primary":
		p.HeaderOnPrimary = pair.merge(p.HeaderOnPrimary)
	case "status_open":
		p.StatusOpen = pair.merge(p.StatusOpen)
	case "status_in_progress":
		p.StatusInProgress = pair.merge(p.StatusInProgress)
	case "status_blocked":
		p.StatusBlocked = pair.merge(p.StatusBlocked)
	case "status_deferred":
		p.StatusDeferred = pair.merge(p.StatusDeferred)
	case "status_pinned":
		p.StatusPinned = pair.merge(p.StatusPinned)
	case "status_hooked":
		p.StatusHooked = pair.merge(p.StatusHooked)
	case "status_review":
		p.StatusReview = pair.merge(p.StatusReview)
	case "status_closed":
		p.StatusClosed = pair.merge(p.StatusClosed)
	case "status_tombstone":
		p.StatusTombstone = pair.merge(p.StatusTombstone)
	case "status_open_bg":
		p.StatusOpenBg = pair.merge(p.StatusOpenBg)
	case "status_in_progress_bg":
		p.StatusInProgressBg = pair.merge(p.StatusInProgressBg)
	case "status_blocked_bg":
		p.StatusBlockedBg = pair.merge(p.StatusBlockedBg)
	case "status_deferred_bg":
		p.StatusDeferredBg = pair.merge(p.StatusDeferredBg)
	case "status_pinned_bg":
		p.StatusPinnedBg = pair.merge(p.StatusPinnedBg)
	case "status_hooked_bg":
		p.StatusHookedBg = pair.merge(p.StatusHookedBg)
	case "status_review_bg":
		p.StatusReviewBg = pair.merge(p.StatusReviewBg)
	case "status_closed_bg":
		p.StatusClosedBg = pair.merge(p.StatusClosedBg)
	case "status_tombstone_bg":
		p.StatusTombstoneBg = pair.merge(p.StatusTombstoneBg)
	case "prio_critical":
		p.PrioCritical = pair.merge(p.PrioCritical)
	case "prio_high":
		p.PrioHigh = pair.merge(p.PrioHigh)
	case "prio_medium":
		p.PrioMedium = pair.merge(p.PrioMedium)
	case "prio_low":
		p.PrioLow = pair.merge(p.PrioLow)
	case "prio_critical_bg":
		p.PrioCriticalBg = pair.merge(p.PrioCriticalBg)
	case "prio_high_bg":
		p.PrioHighBg = pair.merge(p.PrioHighBg)
	case "prio_medium_bg":
		p.PrioMediumBg = pair.merge(p.PrioMediumBg)
	case "prio_low_bg":
		p.PrioLowBg = pair.merge(p.PrioLowBg)
	case "type_bug":
		p.TypeBug = pair.merge(p.TypeBug)
	case "type_feature":
		p.TypeFeature = pair.merge(p.TypeFeature)
	case "type_task":
		p.TypeTask = pair.merge(p.TypeTask)
	case "type_epic":
		p.TypeEpic = pair.merge(p.TypeEpic)
	case "type_chore":
		p.TypeChore = pair.merge(p.TypeChore)
	case "footer_hint":
		p.FooterHint = pair.merge(p.FooterHint)
	case "footer_key":
		p.FooterKey = pair.merge(p.FooterKey)
	case "footer_sep":
		p.FooterSep = pair.merge(p.FooterSep)
	case "footer_dim":
		p.FooterDim = pair.merge(p.FooterDim)
	case "triage_priority_up":
		if s := strings.TrimSpace(pair.Light); s != "" {
			p.TriagePriorityUp = s
		}
		if s := strings.TrimSpace(pair.Dark); s != "" {
			p.TriagePriorityUp = s
		}
	case "triage_priority_down":
		if s := strings.TrimSpace(pair.Light); s != "" {
			p.TriagePriorityDown = s
		}
		if s := strings.TrimSpace(pair.Dark); s != "" {
			p.TriagePriorityDown = s
		}
	case "triage_star":
		if s := strings.TrimSpace(pair.Light); s != "" {
			p.TriageStar = s
		}
		if s := strings.TrimSpace(pair.Dark); s != "" {
			p.TriageStar = s
		}
	case "triage_unblocks":
		if s := strings.TrimSpace(pair.Light); s != "" {
			p.TriageUnblocks = s
		}
		if s := strings.TrimSpace(pair.Dark); s != "" {
			p.TriageUnblocks = s
		}
	case "triage_unblocks_alt":
		if s := strings.TrimSpace(pair.Light); s != "" {
			p.TriageUnblocksAlt = s
		}
		if s := strings.TrimSpace(pair.Dark); s != "" {
			p.TriageUnblocksAlt = s
		}
	default:
		return false
	}
	return true
}

func syncStyleGlobalsFromPalette(p Palette) {
	ColorBg = p.Bg
	ColorBgDark = p.BgDark
	ColorBgSubtle = p.BgSubtle
	ColorBgHighlight = p.BgHighlight
	ColorText = p.Text
	ColorSubtext = p.Subtext
	ColorMuted = p.Muted

	ColorPrimary = p.Primary
	ColorSecondary = p.Secondary
	ColorInfo = p.Info
	ColorSuccess = p.Success
	ColorWarning = p.Warning
	ColorDanger = p.Danger

	ColorStatusOpen = p.StatusOpen
	ColorStatusInProgress = p.StatusInProgress
	ColorStatusBlocked = p.StatusBlocked
	ColorStatusDeferred = p.StatusDeferred
	ColorStatusPinned = p.StatusPinned
	ColorStatusHooked = p.StatusHooked
	ColorStatusReview = p.StatusReview
	ColorStatusClosed = p.StatusClosed
	ColorStatusTombstone = p.StatusTombstone

	ColorStatusOpenBg = p.StatusOpenBg
	ColorStatusInProgressBg = p.StatusInProgressBg
	ColorStatusBlockedBg = p.StatusBlockedBg
	ColorStatusDeferredBg = p.StatusDeferredBg
	ColorStatusPinnedBg = p.StatusPinnedBg
	ColorStatusHookedBg = p.StatusHookedBg
	ColorStatusReviewBg = p.StatusReviewBg
	ColorStatusClosedBg = p.StatusClosedBg
	ColorStatusTombstoneBg = p.StatusTombstoneBg

	ColorPrioCritical = p.PrioCritical
	ColorPrioHigh = p.PrioHigh
	ColorPrioMedium = p.PrioMedium
	ColorPrioLow = p.PrioLow

	ColorPrioCriticalBg = p.PrioCriticalBg
	ColorPrioHighBg = p.PrioHighBg
	ColorPrioMediumBg = p.PrioMediumBg
	ColorPrioLowBg = p.PrioLowBg

	ColorTypeBug = p.TypeBug
	ColorTypeFeature = p.TypeFeature
	ColorTypeTask = p.TypeTask
	ColorTypeEpic = p.TypeEpic
	ColorTypeChore = p.TypeChore

	ColorFooterHint = p.FooterHint
	ColorFooterKey = p.FooterKey
	ColorFooterSep = p.FooterSep
	ColorFooterDim = p.FooterDim

	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBgHighlight)
	FocusedPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary)
}

// PaletteKeyNames returns documented override keys for config help.
func PaletteKeyNames() []string {
	return []string{
		"bg", "bg_dark", "bg_subtle", "bg_highlight",
		"text", "subtext", "muted", "text_base",
		"primary", "secondary", "info", "success", "warning", "danger",
		"border", "highlight", "header_on_primary",
		"status_open", "status_in_progress", "status_blocked", "status_deferred",
		"status_pinned", "status_hooked", "status_review", "status_closed", "status_tombstone",
		"status_open_bg", "status_in_progress_bg", "status_blocked_bg", "status_deferred_bg",
		"status_pinned_bg", "status_hooked_bg", "status_review_bg", "status_closed_bg", "status_tombstone_bg",
		"prio_critical", "prio_high", "prio_medium", "prio_low",
		"prio_critical_bg", "prio_high_bg", "prio_medium_bg", "prio_low_bg",
		"type_bug", "type_feature", "type_task", "type_epic", "type_chore",
		"footer_hint", "footer_key", "footer_sep", "footer_dim",
		"triage_priority_up", "triage_priority_down", "triage_star", "triage_unblocks", "triage_unblocks_alt",
	}
}

// FormatUnknownPaletteKeys formats a warning message for unrecognized config keys.
func FormatUnknownPaletteKeys(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	return fmt.Sprintf("unknown colors keys: %s", strings.Join(keys, ", "))
}
