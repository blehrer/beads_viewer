package icons

import "sort"

// GlossaryEntry documents one logical icon for in-TUI help (K / keywordprg-style lookup).
type GlossaryEntry struct {
	Name        Name
	Category    string
	Description string
}

// Glossary returns all documented icons, sorted by category then name.
func Glossary() []GlossaryEntry {
	out := append([]GlossaryEntry(nil), glossaryEntries...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// GlossaryByCategory groups Glossary entries for display.
func GlossaryByCategory() map[string][]GlossaryEntry {
	grouped := make(map[string][]GlossaryEntry)
	for _, e := range Glossary() {
		grouped[e.Category] = append(grouped[e.Category], e)
	}
	return grouped
}

// ContextualGlossary returns entries relevant to a selected issue (status, type, priority).
// issue may be nil — returns empty slice.
func ContextualGlossary(status, issueType string, priority int) []GlossaryEntry {
	if status == "" && issueType == "" && priority < 0 {
		return nil
	}
	seen := make(map[Name]struct{})
	var out []GlossaryEntry
	add := func(name Name) {
		if _, ok := seen[name]; ok {
			return
		}
		if e, ok := glossaryByName[name]; ok {
			seen[name] = struct{}{}
			out = append(out, e)
		}
	}

	if status != "" {
		switch status {
		case "open":
			add(StatusOpen)
		case "in_progress":
			add(StatusInProgress)
		case "blocked":
			add(StatusBlocked)
		case "closed", "tombstone":
			add(StatusClosed)
		case "deferred", "draft":
			add(Pause)
		case "pinned":
			add(StatusPinned)
		case "hooked":
			add(StatusHooked)
		case "review":
			add(StatusReview)
		default:
			add(StatusUnknown)
		}
	}

	if issueType != "" {
		switch issueType {
		case "bug":
			add(Bug)
		case "feature":
			add(Feature)
		case "task":
			add(Task)
		case "epic":
			add(Epic)
		case "chore":
			add(Chore)
		}
	}

	if priority >= 0 && priority <= 4 {
		switch priority {
		case 0:
			add(Fire)
		case 1:
			add(Lightning)
		case 2:
			add(PriorityMedium)
		case 3:
			add(PriorityLow)
		case 4:
			add(PriorityBacklog)
		}
	}

	return out
}

// Describe returns the human-readable description for a logical icon name.
func Describe(name Name) (string, bool) {
	e, ok := glossaryByName[name]
	if !ok {
		return "", false
	}
	return e.Description, true
}

var glossaryByName map[Name]GlossaryEntry

func init() {
	glossaryByName = make(map[Name]GlossaryEntry, len(glossaryEntries))
	for _, e := range glossaryEntries {
		glossaryByName[e.Name] = e
	}
}

// ponytail: descriptions mention TUI ● where list view uses lipgloss dots instead of emoji circles.
var glossaryEntries = []GlossaryEntry{
	// Status — list view uses theme-colored ●; export/robot JSON uses emoji when BV_ICON_SET=emoji.
	{Name: StatusOpen, Category: "Status", Description: "Open — ready to work. TUI: green ●. Export: 🟢."},
	{Name: StatusInProgress, Category: "Status", Description: "In progress — actively being worked. TUI: cyan ●. Export: 🔵."},
	{Name: StatusBlocked, Category: "Status", Description: "Blocked — waiting on a dependency. TUI: red ●. Export: 🔴."},
	{Name: StatusClosed, Category: "Status", Description: "Closed — done. TUI: muted ● (graph view: ✓). Export: ⚫/✅."},
	{Name: StatusUnknown, Category: "Status", Description: "Unknown or unrecognized status. TUI: muted ●. Export: ⚪."},
	{Name: Pause, Category: "Status", Description: "Deferred or draft — paused, not active work. Export: ⏸️."},
	{Name: StatusPinned, Category: "Status", Description: "Pinned — stays visible across filters. Export: 📌."},
	{Name: StatusHooked, Category: "Status", Description: "Hooked — attached to an agent/session. Export: 🪝."},
	{Name: StatusReview, Category: "Status", Description: "In review — awaiting approval. Export: 👁️."},

	// Priority (beads P0–P4)
	{Name: Fire, Category: "Priority", Description: "P0 Critical — drop everything."},
	{Name: Lightning, Category: "Priority", Description: "P1 High — important, do soon."},
	{Name: PriorityMedium, Category: "Priority", Description: "P2 Medium — normal queue."},
	{Name: PriorityLow, Category: "Priority", Description: "P3 Low — when time allows."},
	{Name: PriorityBacklog, Category: "Priority", Description: "P4 Backlog — someday/maybe."},

	// Issue type
	{Name: Bug, Category: "Type", Description: "Bug — something is broken."},
	{Name: Feature, Category: "Type", Description: "Feature — new capability."},
	{Name: Task, Category: "Type", Description: "Task — concrete unit of work."},
	{Name: Epic, Category: "Type", Description: "Epic — large goal spanning many tasks."},
	{Name: Chore, Category: "Type", Description: "Chore — maintenance or hygiene."},

	// Dependencies
	{Name: Blocked, Category: "Dependency", Description: "Blocks — must complete before dependent work proceeds."},
	{Name: Link, Category: "Dependency", Description: "Related — associated but not blocking."},
	{Name: DepParentChild, Category: "Dependency", Description: "Parent/child — hierarchical containment."},
	{Name: DepDiscovered, Category: "Dependency", Description: "Discovered-from — found while working another issue."},
	{Name: DepRoot, Category: "Dependency", Description: "Root — top of a dependency subtree."},

	// Triage & alerts
	{Name: Target, Category: "Triage", Description: "Recommended focus — top triage pick."},
	{Name: Unlock, Category: "Triage", Description: "Unblocks downstream work — high leverage."},
	{Name: Star, Category: "Triage", Description: "Quick win — low effort, high impact."},
	{Name: Warning, Category: "Triage", Description: "Warning — needs attention."},
	{Name: Siren, Category: "Triage", Description: "Alert — stale, cascade, or priority mismatch."},

	// General UI
	{Name: New, Category: "UI", Description: "New since comparison point (time-travel diff)."},
	{Name: CheckCircle, Category: "UI", Description: "Completed / closed (checkmark)."},
	{Name: Construction, Category: "UI", Description: "Work in progress / under construction."},
	{Name: Chart, Category: "UI", Description: "Metrics or analytics panel."},
	{Name: Clock, Category: "UI", Description: "Time-related — age, staleness, schedule."},
	{Name: User, Category: "UI", Description: "Assignee or agent."},
}
