package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestDefaultPaletteMatchesStyleGlobals(t *testing.T) {
	p := DefaultPalette()
	if ColorPrimary != p.Primary {
		t.Errorf("ColorPrimary = %v, want %v", ColorPrimary, p.Primary)
	}
	if ColorStatusOpen != p.StatusOpen {
		t.Errorf("ColorStatusOpen = %v, want %v", ColorStatusOpen, p.StatusOpen)
	}
}

func TestInitPaletteAppliesOverrides(t *testing.T) {
	prev := activePalette
	prevPrimary := ColorPrimary
	defer func() {
		activePalette = prev
		syncStyleGlobalsFromPalette(prev)
	}()

	InitPalette("dracula", PaletteOverrides{
		"primary": {Light: "#111111", Dark: "#222222"},
	})
	if ColorPrimary.Light != "#111111" || ColorPrimary.Dark != "#222222" {
		t.Fatalf("InitPalette primary override: got %v", ColorPrimary)
	}
	if ActivePalette().Primary != ColorPrimary {
		t.Fatal("ActivePalette out of sync with globals")
	}

	InitPalette("dracula", nil)
	if ColorPrimary != prevPrimary {
		t.Fatal("InitPalette(nil) should restore defaults")
	}
}

func TestPaletteOverridesUnknownKeys(t *testing.T) {
	o := PaletteOverrides{
		"primary":        {Light: "#111111"},
		"not_a_real_key": {Light: "#000000"},
	}
	unknown := o.UnknownKeys()
	if len(unknown) != 1 || unknown[0] != "not_a_real_key" {
		t.Fatalf("UnknownKeys() = %v, want [not_a_real_key]", unknown)
	}
}

func TestDefaultThemeUsesActivePalette(t *testing.T) {
	prev := activePalette
	defer func() {
		activePalette = prev
		syncStyleGlobalsFromPalette(prev)
	}()

	activePalette = DefaultPalette()
	activePalette.Primary = lipgloss.AdaptiveColor{Light: "#ABCDEF", Dark: "#FEDCBA"}
	syncStyleGlobalsFromPalette(activePalette)

	theme := DefaultTheme(lipgloss.NewRenderer(nil))
	if theme.Primary != activePalette.Primary {
		t.Fatalf("DefaultTheme Primary = %v, want %v", theme.Primary, activePalette.Primary)
	}
	if theme.Open != activePalette.StatusOpen {
		t.Fatalf("DefaultTheme Open = %v, want %v", theme.Open, activePalette.StatusOpen)
	}
}

func TestAdaptivePairMerge(t *testing.T) {
	base := lipgloss.AdaptiveColor{Light: "#111111", Dark: "#222222"}
	got := AdaptivePair{Light: "#AAAAAA"}.merge(base)
	if got.Light != "#AAAAAA" || got.Dark != "#222222" {
		t.Fatalf("merge partial light: got %v", got)
	}
}

func TestCanonicalPaletteName(t *testing.T) {
	cases := map[string]string{
		"":          "dracula",
		"dracula":   "dracula",
		"default":   "dracula",
		"kanagawa":  "kanagawa",
		"wave":      "kanagawa",
		"kanso":     "kanso",
		"ink":       "kanso",
		"solarized": "",
	}
	for in, want := range cases {
		if got := CanonicalPaletteName(in); got != want {
			t.Errorf("CanonicalPaletteName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestKanagawaPaletteDistinctFromDracula(t *testing.T) {
	k := kanagawaPalette()
	d := draculaPalette()
	if k.Primary == d.Primary {
		t.Fatal("kanagawa primary should differ from dracula")
	}
	if k.Bg.Dark != "#1F1F28" {
		t.Fatalf("kanagawa dark bg = %q, want #1F1F28", k.Bg.Dark)
	}
}

func TestInitPaletteSelectsBuiltin(t *testing.T) {
	prev := activePalette
	prevName := activePaletteName
	defer func() {
		activePalette = prev
		activePaletteName = prevName
		syncStyleGlobalsFromPalette(prev)
	}()

	InitPalette("kanagawa", nil)
	if ActivePaletteName() != "kanagawa" {
		t.Fatalf("ActivePaletteName() = %q, want kanagawa", ActivePaletteName())
	}
	if ColorPrimary != kanagawaPalette().Primary {
		t.Fatalf("ColorPrimary = %v, want kanagawa primary", ColorPrimary)
	}
}
