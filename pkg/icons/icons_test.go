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
	if got := IssueType("bug"); got != "\U000f00e4" {
		t.Fatalf("IssueType(bug) = %q, want nerd md-bug U+F00E4", got)
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
		{"open", "\U000f0ec2"},
		{"in_progress", "\U000f09de"},
		{"blocked", "\U000f073a"},
		{"closed", "\U000f0159"},
		{"tombstone", "\U000f0159"},
		{"unknown", "\U000f0ec3"},
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

func TestApplyPreference(t *testing.T) {
	t.Cleanup(func() { Use(SetEmoji) })

	ApplyPreference("nerd", "emoji")
	if ActiveSet() != SetNerd {
		t.Fatalf("env over config: ActiveSet() = %v, want SetNerd", ActiveSet())
	}

	ApplyPreference("", "nf")
	if ActiveSet() != SetNerd {
		t.Fatalf("config nerd: ActiveSet() = %v, want SetNerd", ActiveSet())
	}

	ApplyPreference("", "emoji")
	if ActiveSet() != SetEmoji {
		t.Fatalf("config emoji: ActiveSet() = %v, want SetEmoji", ActiveSet())
	}

	ApplyPreference("", "")
	if ActiveSet() != SetEmoji {
		t.Fatalf("default: ActiveSet() = %v, want SetEmoji", ActiveSet())
	}

	ApplyPreference("banana", "nerd")
	if ActiveSet() != SetNerd {
		t.Fatalf("invalid env falls through to config: ActiveSet() = %v, want SetNerd", ActiveSet())
	}
}

func TestParseSet(t *testing.T) {
	cases := []struct {
		in      string
		wantSet Set
		wantOK  bool
	}{
		{"nerd", SetNerd, true},
		{"NERD-FONT", SetNerd, true},
		{" nf ", SetNerd, true},
		{"emoji", SetEmoji, true},
		{"", SetEmoji, false},
		{"banana", SetEmoji, false},
	}
	for _, tc := range cases {
		gotSet, gotOK := ParseSet(tc.in)
		if gotSet != tc.wantSet || gotOK != tc.wantOK {
			t.Errorf("ParseSet(%q) = (%v, %v), want (%v, %v)", tc.in, gotSet, gotOK, tc.wantSet, tc.wantOK)
		}
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
	if got := Get(Bug); got != "\U000f00e4" {
		t.Fatalf("Get(Bug) = %q, want md-bug U+F00E4", got)
	}
}

func TestBugNerd_NotAlphaM(t *testing.T) {
	useSet(t, SetNerd)
	const alphaM = "\U000f0afa" // MDI alpha-m — was wrongly used as md-bug
	if got := Get(Bug); got == alphaM {
		t.Fatal("Bug nerd glyph must not be alpha-m (F0AFA)")
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

func TestPriority(t *testing.T) {
	useSet(t, SetEmoji)
	if got := Priority(0); got != "🔥" {
		t.Fatalf("Priority(0) = %q", got)
	}
	if got := Priority(4); got != "💤" {
		t.Fatalf("Priority(4) = %q", got)
	}
}

func TestDependencyType(t *testing.T) {
	useSet(t, SetEmoji)
	if got := DependencyType("blocks"); got != "⛔" {
		t.Fatalf("DependencyType(blocks) = %q", got)
	}
}

func TestIssueStatusGraph_Nerd(t *testing.T) {
	useSet(t, SetNerd)
	if got := IssueStatusGraph("open"); got != "\U000f09de" {
		t.Fatalf("IssueStatusGraph(open) = %q", got)
	}
	if got := IssueStatusGraph("closed"); got != "\U000f05e0" {
		t.Fatalf("IssueStatusGraph(closed) = %q", got)
	}
}

func TestIssueStatusGraph(t *testing.T) {
	useSet(t, SetEmoji)
	if got := IssueStatusGraph("open"); got != "🔵" {
		t.Fatalf("IssueStatusGraph(open) = %q", got)
	}
	if got := IssueStatusGraph("closed"); got != "✅" {
		t.Fatalf("IssueStatusGraph(closed) = %q", got)
	}
}

func TestIssueStatus_Extended(t *testing.T) {
	useSet(t, SetEmoji)
	if got := IssueStatus("deferred"); got != "⏸️" {
		t.Fatalf("IssueStatus(deferred) = %q", got)
	}
}

func TestPriority_AllLevels(t *testing.T) {
	useSet(t, SetEmoji)
	want := []string{"🔥", "⚡", "🔹", "☕", "💤"}
	for level, w := range want {
		if got := Priority(level); got != w {
			t.Errorf("Priority(%d) = %q, want %q", level, got, w)
		}
	}
	if got := Priority(99); got != "  " {
		t.Fatalf("Priority(99) = %q", got)
	}
}

func TestPriorityLabel(t *testing.T) {
	useSet(t, SetEmoji)
	if got := PriorityLabel(0); got != "🔥 Critical (P0)" {
		t.Fatalf("PriorityLabel(0) = %q", got)
	}
	if got := PriorityLabel(2); got != "🔹 Medium (P2)" {
		t.Fatalf("PriorityLabel(2) = %q", got)
	}
	if got := PriorityLabel(99); got != "P99" {
		t.Fatalf("PriorityLabel(99) = %q", got)
	}
}

func TestDependencyType_All(t *testing.T) {
	useSet(t, SetEmoji)
	tests := []struct {
		typ  string
		want string
	}{
		{"root", "📍"},
		{"blocks", "⛔"},
		{"related", "🔗"},
		{"parent-child", "📦"},
		{"discovered-from", "🔍"},
		{"unknown", "•"},
	}
	for _, tt := range tests {
		if got := DependencyType(tt.typ); got != tt.want {
			t.Errorf("DependencyType(%q) = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestIssueStatusGraph_All(t *testing.T) {
	useSet(t, SetEmoji)
	tests := []struct {
		status string
		want   string
	}{
		{"open", "🔵"},
		{"in_progress", "🟡"},
		{"blocked", "🔴"},
		{"closed", "✅"},
		{"tombstone", "✅"},
		{"deferred", "⏸️"},
		{"pinned", "📌"},
		{"hooked", "🪝"},
		{"review", "👁️"},
		{"unknown", "⚪"},
	}
	for _, tt := range tests {
		if got := IssueStatusGraph(tt.status); got != tt.want {
			t.Errorf("IssueStatusGraph(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestIssueStatus_ExtendedAll(t *testing.T) {
	useSet(t, SetEmoji)
	tests := []struct {
		status string
		want   string
	}{
		{"pinned", "📌"},
		{"hooked", "🪝"},
		{"review", "👁️"},
		{"draft", "⏸️"},
	}
	for _, tt := range tests {
		if got := IssueStatus(tt.status); got != tt.want {
			t.Errorf("IssueStatus(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestAllNamesHaveGlyphs(t *testing.T) {
	names := []Name{
		Bug, Feature, Task, Epic, Chore, Target, Unlock, Warning, Shuffle, Chart,
		Clock, Calendar, Lightning, CheckCircle, Construction, Blocked, Pause, User,
		Hourglass, Siren, Star, Fire, Alarm, Link, New, Check, Cross,
		StatusOpen, StatusInProgress, StatusBlocked, StatusClosed, StatusUnknown,
		StatusDeferred, StatusPinned, StatusHooked, StatusReview,
		StatusGraphOpen, StatusGraphWork, TriageScoreMid,
		PriorityMedium, PriorityLow, PriorityBacklog,
		DepRoot, DepParentChild, DepDiscovered,
		SwimRefresh, SwimProhibited, Question, FileDefault,
		HistoryGit, HistoryBead, SessionAttach, ArrowUp, ArrowDown,
		Folder, Globe, Comment, Label, Navigation, Lightbulb, Bell,
		Microscope, Books, Brain, Building, Satellite, CriticalPath, Health,
		HistoryScroll, CommitRefactor, CommitTest, CommitChore, CommitStyle, CommitRevert,
		Legend, TrendUp, Eye,
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
