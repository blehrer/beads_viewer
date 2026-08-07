package icons

// HistoryBeadStatus returns a compact status glyph for history list rows and related-bead panels.
// Emoji/default mode keeps single-width ○ ● ✓ for column alignment; nerd mode uses registry glyphs.
func HistoryBeadStatus(status string) string {
	if ActiveSet() == SetNerd {
		switch status {
		case "closed", "tombstone":
			return Get(Check)
		case "in_progress":
			return Get(StatusInProgress)
		case "blocked":
			return Get(StatusBlocked)
		default:
			return Get(StatusGraphOpen)
		}
	}
	switch status {
	case "closed", "tombstone":
		return Get(Check)
	case "in_progress":
		return "●"
	default:
		return "○"
	}
}

// LifecycleEvent returns the icon for a bead lifecycle event type (created, claimed, closed, …).
func LifecycleEvent(eventType string) string {
	switch eventType {
	case "created":
		return Get(New)
	case "claimed":
		return Get(User)
	case "closed":
		return Get(Check)
	case "reopened":
		return Get(SwimRefresh)
	case "modified":
		return Get(FileDefault)
	default:
		return "•"
	}
}

// TimelineMilestone returns a compact marker for lifecycle milestones in the one-line timeline.
func TimelineMilestone(kind string) string {
	if ActiveSet() == SetNerd {
		switch kind {
		case "created":
			return Get(StatusGraphOpen)
		case "claimed":
			return Get(StatusInProgress)
		case "closed":
			return Get(Check)
		}
		return "•"
	}
	switch kind {
	case "created":
		return "○"
	case "claimed":
		return "●"
	case "closed":
		return Get(Check)
	default:
		return "•"
	}
}

// HistoryViewModeIcon returns the history header icon for git-centric vs bead-centric mode.
func HistoryViewModeIcon(gitMode bool) string {
	if gitMode {
		return Get(HistoryGit)
	}
	return Get(HistoryBead)
}
