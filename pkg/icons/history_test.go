package icons

import "testing"

func TestHistoryBeadStatus_NerdMode(t *testing.T) {
	Use(SetNerd)
	t.Cleanup(func() { Use(SetEmoji) })

	if got := HistoryBeadStatus("open"); got == "○" {
		t.Fatalf("open status should use nerd glyph, got %q", got)
	}
	if got := HistoryBeadStatus("closed"); got == "✓" {
		t.Fatalf("closed status should use nerd glyph, got %q", got)
	}
}

func TestHistoryViewModeIcon_NerdMode(t *testing.T) {
	Use(SetNerd)
	t.Cleanup(func() { Use(SetEmoji) })

	if got := HistoryViewModeIcon(true); got == "◉" {
		t.Fatalf("git mode should use nerd glyph, got %q", got)
	}
	if got := HistoryViewModeIcon(false); got == "◈" {
		t.Fatalf("bead mode should use nerd glyph, got %q", got)
	}
}

func TestLifecycleEvent_UsesRegistry(t *testing.T) {
	Use(SetEmoji)
	if got := LifecycleEvent("created"); got != Get(New) {
		t.Fatalf("created = %q, want %q", got, Get(New))
	}
	if got := LifecycleEvent("claimed"); got != Get(User) {
		t.Fatalf("claimed = %q, want %q", got, Get(User))
	}
}
