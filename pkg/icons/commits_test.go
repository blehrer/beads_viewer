package icons

import "testing"

func TestConventionalCommitIcon(t *testing.T) {
	useSet(t, SetEmoji)
	tests := []struct {
		typ  string
		want string
	}{
		{"feat", Get(Feature)},
		{"fix", Get(Bug)},
		{"docs", Get(FileDefault)},
		{"refactor", Get(CommitRefactor)},
		{"perf", Get(Lightning)},
		{"test", Get(CommitTest)},
		{"chore", Get(CommitChore)},
		{"ci", Get(SwimRefresh)},
		{"build", Get(DepParentChild)},
		{"style", Get(CommitStyle)},
		{"unknown", ""},
	}
	for _, tt := range tests {
		if got := ConventionalCommitIcon(tt.typ); got != tt.want {
			t.Errorf("ConventionalCommitIcon(%q) = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestConventionalCommitIcon_Nerd(t *testing.T) {
	useSet(t, SetNerd)
	if got := ConventionalCommitIcon("feat"); got != Get(Feature) {
		t.Fatalf("ConventionalCommitIcon(feat) = %q, want nerd feature glyph", got)
	}
}

func TestCommitRevertIcon(t *testing.T) {
	useSet(t, SetEmoji)
	if got := CommitRevertIcon(); got != Get(CommitRevert) {
		t.Fatalf("CommitRevertIcon() = %q, want %q", got, Get(CommitRevert))
	}
}
