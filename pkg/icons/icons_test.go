package icons

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	Use(SetEmoji)
	os.Exit(code)
}

func useSet(t *testing.T, s Set) {
	t.Helper()
	Use(s)
	t.Cleanup(func() { Use(SetEmoji) })
}

func TestIssueType_EmojiDefault(t *testing.T) {
	useSet(t, SetEmoji)
	tests := []struct {
		typ  string
		want string
	}{
		{"bug", "🐛"},
		{"feature", "✨"},
		{"task", "📋"},
		{"epic", "🚀"},
		{"chore", "🧹"},
		{"unknown", "•"},
	}
	for _, tt := range tests {
		if got := IssueType(tt.typ); got != tt.want {
			t.Errorf("IssueType(%q) = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestIssueType_Nerd(t *testing.T) {
	useSet(t, SetNerd)
	if got := IssueType("bug"); got != "\U000f0afa" {
		t.Fatalf("IssueType(bug) = %q, want nerd bug", got)
	}
	if got := IssueType("unknown"); got != "•" {
		t.Fatalf("IssueType(unknown) = %q, want bullet", got)
	}
}

func TestIssueStatus_Emoji(t *testing.T) {
	useSet(t, SetEmoji)
	tests := []struct {
		status string
		want   string
	}{
		{"open", "🟢"},
		{"in_progress", "🔵"},
		{"blocked", "🔴"},
		{"closed", "⚫"},
		{"tombstone", "⚫"},
		{"unknown", "⚪"},
	}
	for _, tt := range tests {
		if got := IssueStatus(tt.status); got != tt.want {
			t.Errorf("IssueStatus(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestIssueStatus_Nerd(t *testing.T) {
	useSet(t, SetNerd)
	tests := []struct {
		status string
		want   string
	}{
		{"open", "\U000f0765"},
		{"in_progress", "\U000f0130"},
		{"blocked", "\U000f073a"},
		{"closed", "\U000f0133"},
		{"tombstone", "\U000f0133"},
		{"unknown", "\U000f0766"},
	}
	for _, tt := range tests {
		if got := IssueStatus(tt.status); got != tt.want {
			t.Errorf("IssueStatus(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestSetFromEnv(t *testing.T) {
	t.Setenv("BV_ICON_SET", "nerd-font")
	SetFromEnv()
	t.Cleanup(func() { Use(SetEmoji) })
	if ActiveSet() != SetNerd {
		t.Fatalf("ActiveSet() = %v, want SetNerd", ActiveSet())
	}
	t.Setenv("BV_ICON_SET", "emoji")
	SetFromEnv()
	if ActiveSet() != SetEmoji {
		t.Fatalf("ActiveSet() = %v, want SetEmoji", ActiveSet())
	}
}

func TestGet_Emoji(t *testing.T) {
	useSet(t, SetEmoji)
	if Get(Target) != "🎯" {
		t.Fatalf("Get(Target) = %q", Get(Target))
	}
}

func TestGet_Nerd(t *testing.T) {
	useSet(t, SetNerd)
	if got := Get(Bug); got != "\U000f0afa" {
		t.Fatalf("Get(Bug) = %q, want nerd bug", got)
	}
}

func TestGet_UnknownName(t *testing.T) {
	useSet(t, SetEmoji)
	if got := Get(Name("nonexistent")); got != "?" {
		t.Fatalf("Get(unknown) = %q, want ?", got)
	}
}

func TestGet_NerdFallbackToEmoji(t *testing.T) {
	useSet(t, SetNerd)
	// StatusBlocked shares md-cancel with Blocked; pick a name only in emoji map by
	// using a synthetic name that nerdIcons lacks but we can verify via deletion pattern:
	// Temporarily verify fallback path: unknown nerd entry returns emoji if present.
	const fake Name = "fake_only_emoji"
	emojiIcons[fake] = "X"
	t.Cleanup(func() { delete(emojiIcons, fake) })
	if got := Get(fake); got != "X" {
		t.Fatalf("Get(fake) nerd fallback = %q, want X", got)
	}
}

func TestSetFromEnv_Aliases(t *testing.T) {
	t.Cleanup(func() { Use(SetEmoji) })
	for _, val := range []string{"nerdfont", "nf", "  NERD  "} {
		t.Setenv("BV_ICON_SET", val)
		SetFromEnv()
		if ActiveSet() != SetNerd {
			t.Fatalf("BV_ICON_SET=%q: ActiveSet() = %v, want SetNerd", val, ActiveSet())
		}
	}
}

func TestAllNamesHaveGlyphs(t *testing.T) {
	names := []Name{
		Bug, Feature, Task, Epic, Chore, Target, Unlock, Warning, Shuffle, Chart,
		Clock, Calendar, Lightning, CheckCircle, Construction, Blocked, Pause, User,
		Hourglass, Siren, Star, Fire, Alarm, Link, New, Check, Cross,
		StatusOpen, StatusInProgress, StatusBlocked, StatusClosed, StatusUnknown,
	}
	for _, name := range names {
		if _, ok := emojiIcons[name]; !ok {
			t.Errorf("emojiIcons missing %q", name)
		}
		if _, ok := nerdIcons[name]; !ok {
			t.Errorf("nerdIcons missing %q", name)
		}
	}
}
