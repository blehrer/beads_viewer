package icons

// ConventionalCommitIcon returns the glyph for a conventional-commit type prefix.
func ConventionalCommitIcon(commitType string) string {
	switch commitType {
	case "feat":
		return Get(Feature)
	case "fix":
		return Get(Bug)
	case "docs":
		return Get(FileDefault)
	case "refactor":
		return Get(CommitRefactor)
	case "perf":
		return Get(Lightning)
	case "test":
		return Get(CommitTest)
	case "chore":
		return Get(CommitChore)
	case "ci":
		return Get(SwimRefresh)
	case "build":
		return Get(DepParentChild)
	case "style":
		return Get(CommitStyle)
	default:
		return ""
	}
}

// CommitRevertIcon returns the glyph for revert commits.
func CommitRevertIcon() string {
	return Get(CommitRevert)
}
