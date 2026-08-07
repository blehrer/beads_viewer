package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func adapt(light, dark string) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: light, Dark: dark}
}

// CanonicalPaletteName maps user input to a built-in palette name.
// Returns "" for unrecognized values.
func CanonicalPaletteName(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "dracula", "default":
		return "dracula"
	case "kanagawa", "kanagawa-wave", "wave":
		return "kanagawa"
	case "kanso", "kanso-ink", "ink":
		return "kanso"
	default:
		return ""
	}
}

// BuiltinPaletteNames lists selectable built-in palettes.
func BuiltinPaletteNames() []string {
	return []string{"dracula", "kanagawa", "kanso"}
}

// PaletteByName returns a built-in palette. Unknown names fall back to dracula.
func PaletteByName(name string) Palette {
	switch CanonicalPaletteName(name) {
	case "kanagawa":
		return kanagawaPalette()
	case "kanso":
		return kansoPalette()
	default:
		return draculaPalette()
	}
}

// draculaPalette is the original bv default (Dracula-inspired, WCAG-tuned light mode).
func draculaPalette() Palette {
	return Palette{
		Bg:          adapt("#FFFFFF", "#282A36"),
		BgDark:      adapt("#F5F5F5", "#1E1F29"),
		BgSubtle:    adapt("#E8E8E8", "#363949"),
		BgHighlight: adapt("#D0D0D0", "#44475A"),
		Text:        adapt("#1A1A1A", "#F8F8F2"),
		Subtext:     adapt("#555555", "#BFBFBF"),
		Muted:       adapt("#666666", "#6272A4"),
		TextBase:    adapt("#000000", "#F8F8F2"),

		Primary:   adapt("#6B47D9", "#BD93F9"),
		Secondary: adapt("#555555", "#6272A4"),
		Info:      adapt("#006080", "#8BE9FD"),
		Success:   adapt("#007700", "#50FA7B"),
		Warning:   adapt("#B06800", "#FFB86C"),
		Danger:    adapt("#CC0000", "#FF5555"),

		Border:          adapt("#AAAAAA", "#44475A"),
		Highlight:       adapt("#E0E0E0", "#44475A"),
		HeaderOnPrimary: adapt("#FFFFFF", "#282A36"),

		StatusOpen:       adapt("#007700", "#50FA7B"),
		StatusInProgress: adapt("#006080", "#8BE9FD"),
		StatusBlocked:    adapt("#CC0000", "#FF5555"),
		StatusDeferred:   adapt("#B06800", "#FFB86C"),
		StatusPinned:     adapt("#0066CC", "#6699FF"),
		StatusHooked:     adapt("#008080", "#00CED1"),
		StatusReview:     adapt("#6B47D9", "#BD93F9"),
		StatusClosed:     adapt("#555555", "#6272A4"),
		StatusTombstone:  adapt("#888888", "#44475A"),

		StatusOpenBg:       adapt("#D4EDDA", "#1A3D2A"),
		StatusInProgressBg: adapt("#D1ECF1", "#1A3344"),
		StatusBlockedBg:    adapt("#F8D7DA", "#3D1A1A"),
		StatusDeferredBg:   adapt("#FFE8CC", "#3D2A1A"),
		StatusPinnedBg:     adapt("#CCE5FF", "#1A2A44"),
		StatusHookedBg:     adapt("#CCFFFF", "#1A3D3D"),
		StatusReviewBg:     adapt("#E8DDFF", "#2A1A44"),
		StatusClosedBg:     adapt("#E2E3E5", "#2A2A3D"),
		StatusTombstoneBg:  adapt("#D0D0D0", "#1E1F29"),

		PrioCritical: adapt("#CC0000", "#FF5555"),
		PrioHigh:     adapt("#B06800", "#FFB86C"),
		PrioMedium:   adapt("#0066CC", "#6699FF"),
		PrioLow:      adapt("#007700", "#50FA7B"),

		PrioCriticalBg: adapt("#F8D7DA", "#3D1A1A"),
		PrioHighBg:     adapt("#FFE8CC", "#3D2A1A"),
		PrioMediumBg:   adapt("#CCE5FF", "#1A2A44"),
		PrioLowBg:      adapt("#D4EDDA", "#1A3D2A"),

		TypeBug:     adapt("#CC0000", "#FF5555"),
		TypeFeature: adapt("#B06800", "#FFB86C"),
		TypeTask:    adapt("#808000", "#F1FA8C"),
		TypeEpic:    adapt("#6B47D9", "#BD93F9"),
		TypeChore:   adapt("#006080", "#8BE9FD"),

		FooterHint: adapt("#444444", "#C8C8D0"),
		FooterKey:  adapt("#333333", "#E0E0E8"),
		FooterSep:  adapt("#888888", "#8888A0"),
		FooterDim:  adapt("#555555", "#A0A0B8"),

		TriagePriorityUp:   "#FF6B6B",
		TriagePriorityDown: "#4ECDC4",
		TriageStar:         "#FFD700",
		TriageUnblocks:     "#50FA7B",
		TriageUnblocksAlt:  "#6272A4",
	}
}

// kanagawaPalette maps rebelot/kanagawa.nvim wave (dark) + lotus (light) to bv slots.
func kanagawaPalette() Palette {
	return Palette{
		Bg:          adapt("#f2ecbc", "#1F1F28"),
		BgDark:      adapt("#e5ddb0", "#16161D"),
		BgSubtle:    adapt("#e7dba0", "#2A2A37"),
		BgHighlight: adapt("#e4d794", "#363646"),
		Text:        adapt("#545464", "#DCD7BA"),
		Subtext:     adapt("#43436c", "#C8C093"),
		Muted:       adapt("#8a8980", "#727169"),
		TextBase:    adapt("#545464", "#DCD7BA"),

		Primary:   adapt("#624c83", "#957FB8"),
		Secondary: adapt("#766b90", "#938AA9"),
		Info:      adapt("#4e8ca2", "#7FB4CA"),
		Success:   adapt("#6f894e", "#98BB6C"),
		Warning:   adapt("#e98a00", "#FF9E3B"),
		Danger:    adapt("#e82424", "#E82424"),

		Border:          adapt("#716e61", "#54546D"),
		Highlight:       adapt("#c9cbd1", "#363646"),
		HeaderOnPrimary: adapt("#f2ecbc", "#16161D"),

		StatusOpen:       adapt("#6f894e", "#98BB6C"),
		StatusInProgress: adapt("#4e8ca2", "#7FB4CA"),
		StatusBlocked:    adapt("#e82424", "#E82424"),
		StatusDeferred:   adapt("#e98a00", "#FFA066"),
		StatusPinned:     adapt("#4d699b", "#7E9CD8"),
		StatusHooked:     adapt("#597b75", "#7AA89F"),
		StatusReview:     adapt("#624c83", "#957FB8"),
		StatusClosed:     adapt("#8a8980", "#727169"),
		StatusTombstone:  adapt("#716e61", "#717C7C"),

		StatusOpenBg:       adapt("#b7d0ae", "#2B3328"),
		StatusInProgressBg: adapt("#c7d7e0", "#252535"),
		StatusBlockedBg:    adapt("#d9a594", "#43242B"),
		StatusDeferredBg:   adapt("#f9d791", "#49443C"),
		StatusPinnedBg:     adapt("#b5cbd2", "#223249"),
		StatusHookedBg:     adapt("#d7e3d8", "#223249"),
		StatusReviewBg:     adapt("#c9cbd1", "#2A2A37"),
		StatusClosedBg:     adapt("#e7dba0", "#363646"),
		StatusTombstoneBg:  adapt("#dcd7ba", "#16161D"),

		PrioCritical: adapt("#e82424", "#FF5D62"),
		PrioHigh:     adapt("#de9800", "#FFA066"),
		PrioMedium:   adapt("#4d699b", "#7E9CD8"),
		PrioLow:      adapt("#6f894e", "#98BB6C"),

		PrioCriticalBg: adapt("#d9a594", "#43242B"),
		PrioHighBg:     adapt("#f9d791", "#49443C"),
		PrioMediumBg:   adapt("#b5cbd2", "#223249"),
		PrioLowBg:      adapt("#b7d0ae", "#2B3328"),

		TypeBug:     adapt("#c84053", "#FF5D62"),
		TypeFeature: adapt("#cc6d00", "#FFA066"),
		TypeTask:    adapt("#77713f", "#C0A36E"),
		TypeEpic:    adapt("#624c83", "#957FB8"),
		TypeChore:   adapt("#4e8ca2", "#7FB4CA"),

		FooterHint: adapt("#43436c", "#C8C093"),
		FooterKey:  adapt("#545464", "#DCD7BA"),
		FooterSep:  adapt("#8a8980", "#727169"),
		FooterDim:  adapt("#766b90", "#938AA9"),

		TriagePriorityUp:   "#FF5D62",
		TriagePriorityDown: "#7AA89F",
		TriageStar:         "#E6C384",
		TriageUnblocks:     "#98BB6C",
		TriageUnblocksAlt:  "#938AA9",
	}
}

// kansoPalette maps webhooked/kanso.nvim ink (dark) + pearl (light) to bv slots.
func kansoPalette() Palette {
	return Palette{
		Bg:          adapt("#f2f1ef", "#14171d"),
		BgDark:      adapt("#e2e1df", "#14171d"),
		BgSubtle:    adapt("#dddddb", "#1f1f26"),
		BgHighlight: adapt("#cacac7", "#22262D"),
		Text:        adapt("#545464", "#C5C9C7"),
		Subtext:     adapt("#5C6068", "#A4A7A4"),
		Muted:       adapt("#6D6D69", "#717C7C"),
		TextBase:    adapt("#545464", "#C5C9C7"),

		Primary:   adapt("#624c83", "#8992a7"),
		Secondary: adapt("#766b90", "#938AA9"),
		Info:      adapt("#4e8ca2", "#7FB4CA"),
		Success:   adapt("#6f894e", "#98BB6C"),
		Warning:   adapt("#de9800", "#DCA561"),
		Danger:    adapt("#e82424", "#C34043"),

		Border:          adapt("#9F9F99", "#393B44"),
		Highlight:       adapt("#cacac7", "#393B44"),
		HeaderOnPrimary: adapt("#f2f1ef", "#14171d"),

		StatusOpen:       adapt("#6f894e", "#98BB6C"),
		StatusInProgress: adapt("#4e8ca2", "#7FB4CA"),
		StatusBlocked:    adapt("#e82424", "#C34043"),
		StatusDeferred:   adapt("#de9800", "#DCA561"),
		StatusPinned:     adapt("#4d699b", "#658594"),
		StatusHooked:     adapt("#597b75", "#6A9589"),
		StatusReview:     adapt("#624c83", "#8992a7"),
		StatusClosed:     adapt("#6D6D69", "#717C7C"),
		StatusTombstone:  adapt("#9F9F99", "#5C6066"),

		StatusOpenBg:       adapt("#b7d0ae", "#2B3328"),
		StatusInProgressBg: adapt("#c7d7e0", "#252535"),
		StatusBlockedBg:    adapt("#d9a594", "#43242B"),
		StatusDeferredBg:   adapt("#f9d791", "#49443C"),
		StatusPinnedBg:     adapt("#b5cbd2", "#223249"),
		StatusHookedBg:     adapt("#d7e3d8", "#223249"),
		StatusReviewBg:     adapt("#c9cbd1", "#393B44"),
		StatusClosedBg:     adapt("#e2e1df", "#393B44"),
		StatusTombstoneBg:  adapt("#cacac7", "#14171d"),

		PrioCritical: adapt("#e82424", "#C34043"),
		PrioHigh:     adapt("#de9800", "#DCA561"),
		PrioMedium:   adapt("#4d699b", "#658594"),
		PrioLow:      adapt("#6f894e", "#98BB6C"),

		PrioCriticalBg: adapt("#d9a594", "#43242B"),
		PrioHighBg:     adapt("#f9d791", "#49443C"),
		PrioMediumBg:   adapt("#b5cbd2", "#223249"),
		PrioLowBg:      adapt("#b7d0ae", "#2B3328"),

		TypeBug:     adapt("#c84053", "#C34043"),
		TypeFeature: adapt("#cc6d00", "#b6927b"),
		TypeTask:    adapt("#77713f", "#8a9a7b"),
		TypeEpic:    adapt("#624c83", "#8992a7"),
		TypeChore:   adapt("#4e8ca2", "#7FB4CA"),

		FooterHint: adapt("#5C6068", "#A4A7A4"),
		FooterKey:  adapt("#545464", "#C5C9C7"),
		FooterSep:  adapt("#6D6D69", "#717C7C"),
		FooterDim:  adapt("#766b90", "#909398"),

		TriagePriorityUp:   "#E46876",
		TriagePriorityDown: "#7AA89F",
		TriageStar:         "#E6C384",
		TriageUnblocks:     "#98BB6C",
		TriageUnblocksAlt:  "#8992a7",
	}
}
