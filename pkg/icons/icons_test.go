package icons

import "testing"

func TestIssueType_EmojiDefault(t *testing.T) {
	Use(SetEmoji)
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
	Use(SetNerd)
	if got := IssueType("bug"); got != nerdIcons[Bug] {
		t.Fatalf("IssueType(bug) = %q, want nerd bug %q", got, nerdIcons[Bug])
	}
	if got := IssueType("unknown"); got != "•" {
		t.Fatalf("IssueType(unknown) = %q, want bullet", got)
	}
}

func TestIssueStatus_Switch(t *testing.T) {
	Use(SetEmoji)
	if got := IssueStatus("open"); got != "🟢" {
		t.Fatalf("emoji open = %q", got)
	}
	Use(SetNerd)
	if got := IssueStatus("open"); got != "\U000f0765" {
		t.Fatalf("nerd open = %q", got)
	}
}

func TestSetFromEnv(t *testing.T) {
	t.Setenv("BV_ICON_SET", "nerd-font")
	SetFromEnv()
	if ActiveSet() != SetNerd {
		t.Fatalf("ActiveSet() = %v, want SetNerd", ActiveSet())
	}
	t.Setenv("BV_ICON_SET", "emoji")
	SetFromEnv()
	if ActiveSet() != SetEmoji {
		t.Fatalf("ActiveSet() = %v, want SetEmoji", ActiveSet())
	}
}

func TestGet_FallsBackToEmoji(t *testing.T) {
	Use(SetNerd)
	Use(SetEmoji)
	if Get(Target) != "🎯" {
		t.Fatalf("Get(Target) = %q", Get(Target))
	}
}
