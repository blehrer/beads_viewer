package ui

import (
	"sort"
	"strings"

	"github.com/Dicklesworthstone/beads_viewer/pkg/icons"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// HintSurface selects which UI surface consumes registry hints.
type HintSurface int

const (
	HintFooter HintSurface = iota
	HintHelpOverlay
	HintSidebar
)

const footerHintLimit = 8

// helpOverlayPanelSpec describes one help overlay column panel.
type helpOverlayPanelSpec struct {
	title    string
	icon     icons.Name
	colorIdx int
	category string // registry category filter; empty = static content
	static   []keyDescPair
}

type keyDescPair struct {
	key  string
	desc string
}

// footerSubjectContext returns the UI context that governs footer key hints.
// It reflects modal/submode state that overrides generic list/view bindings.
func (m Model) footerSubjectContext() Context {
	if m.isBoardView && m.board.IsSearchMode() {
		return ContextBoardSearch
	}
	if m.focused == focusHistory && m.historyView.IsSearchActive() {
		return ContextHistorySearch
	}
	return m.CurrentContext()
}

// hintContext is deprecated; use footerSubjectContext.
func (m Model) hintContext() Context { return m.footerSubjectContext() }

// helpSubjectContext returns the view context shown in the ? help overlay
// (the view the user came from, not the help overlay itself).
func (m Model) helpSubjectContext() Context {
	if m.showHelp {
		return m.contextFromFocus(m.focusBeforeHelp)
	}
	return m.CurrentContext()
}

func (m Model) contextFromFocus(f focus) Context {
	switch f {
	case focusBoard:
		return ContextBoard
	case focusGraph:
		return ContextGraph
	case focusTree:
		return ContextGraph
	case focusInsights:
		if m.showAttentionView {
			return ContextAttention
		}
		return ContextInsights
	case focusHistory:
		return ContextHistory
	case focusActionable:
		return ContextActionable
	case focusLabelDashboard:
		return ContextLabelDashboard
	case focusFlowMatrix:
		return ContextFlowMatrix
	case focusDetail:
		if m.timeTravelMode {
			return ContextTimeTravel
		}
		return ContextDetail
	case focusList:
		if m.timeTravelMode {
			return ContextTimeTravel
		}
		if m.list.FilterState() != list.Unfiltered {
			return ContextFilter
		}
		if m.isSplitView {
			return ContextSplit
		}
		if m.showDetails {
			return ContextDetail
		}
		return ContextList
	default:
		return ContextList
	}
}

func contextToDocContext(ctx Context) string {
	switch ctx {
	case ContextList:
		return "list"
	case ContextDetail:
		return "detail"
	case ContextSplit:
		return "list"
	case ContextFilter:
		return "filter"
	case ContextBoard:
		return "board"
	case ContextBoardSearch:
		return "board-search"
	case ContextGraph:
		return "graph"
	case ContextInsights:
		return "insights"
	case ContextAttention:
		return "attention"
	case ContextHistory:
		return "history"
	case ContextActionable:
		return "actionable"
	case ContextLabelDashboard:
		return "label-dashboard"
	case ContextLabelPicker:
		return "label-picker"
	case ContextRecipePicker:
		return "recipe-picker"
	case ContextRepoPicker:
		return "repo-picker"
	case ContextFlowMatrix:
		return "flow-matrix"
	case ContextHelp:
		return "help"
	case ContextContextHelp:
		return "context-help"
	case ContextTutorial:
		return "tutorial"
	case ContextGlyphHelp:
		return "glyph-help"
	case ContextQuitConfirm:
		return "quit-confirm"
	case ContextUpdateModal:
		return "update-modal"
	case ContextHistorySearch:
		return "history-search"
	case ContextAlerts:
		return "alerts"
	case ContextAgentPrompt:
		return "agent-prompt"
	case ContextCassSession:
		return "cass-session"
	case ContextLabelHealthDetail:
		return "label-health-detail"
	case ContextLabelDrilldown:
		return "label-drilldown"
	case ContextLabelGraphAnalysis:
		return "label-graph-analysis"
	case ContextTimeTravelInput:
		return "list"
	case ContextTimeTravel:
		return "list"
	default:
		return string(ctx)
	}
}

func docContextsForUI(ctx Context, surface HintSurface) []string {
	primary := contextToDocContext(ctx)
	if primary == "" {
		return nil
	}

	if surface == HintFooter && ctx == ContextFilter {
		return []string{"filter"}
	}

	seen := map[string]struct{}{primary: {}}
	contexts := []string{primary}

	add := func(raw ...string) {
		for _, c := range raw {
			if c == "" {
				continue
			}
			if _, ok := seen[c]; ok {
				continue
			}
			seen[c] = struct{}{}
			contexts = append(contexts, c)
		}
	}

	switch surface {
	case HintFooter, HintHelpOverlay:
		if ctx.AllowsGlobalFallthrough() {
			switch ctx {
			case ContextInsights, ContextAttention:
				add("insights", "attention", "list", "detail")
			case ContextGraph:
				add("graph", "list", "detail")
			case ContextBoard, ContextBoardSearch:
				add("board", "board-search", "list", "detail")
			case ContextHistory, ContextActionable, ContextFlowMatrix, ContextLabelDashboard:
				add(contextToDocContext(ctx), "list", "detail")
			case ContextSplit, ContextDetail, ContextTimeTravel:
				add("list", "detail")
			case ContextFilter:
				add("filter", "list")
			default:
				if primary != "list" {
					add("list", "detail")
				}
			}
		}
		if surface == HintHelpOverlay {
			add("all")
		}
		if surface == HintFooter && ctx.AllowsGlobalFallthrough() {
			add("all")
		}
	case HintSidebar:
		add("all")
	}

	return contexts
}

func docAppliesToContexts(doc KeyBindingDoc, docContexts []string) bool {
	for _, raw := range strings.Split(doc.Context, ",") {
		raw = strings.TrimSpace(raw)
		for _, want := range docContexts {
			if raw == want {
				return true
			}
		}
	}
	return false
}

func footerExcludedKey(ctx Context, key string) bool {
	switch ctx {
	case ContextHistory:
		switch key {
		case "h", "q", "esc", "tab", "enter":
			return false
		}
	case ContextInsights, ContextAttention:
		switch key {
		case "h", "l", "e", "enter", "?", "]", "f4", "f":
			return false
		case "ctrl+j", "ctrl+k", "x", "m":
			return true
		}
	case ContextGraph:
		switch key {
		case "h", "l", "j", "k", "enter", "g":
			return false
		}
	case ContextBoard, ContextBoardSearch:
		switch key {
		case "h", "l", "j", "k", "G", "enter", "b", "n", "N", "esc":
			return false
		}
		return true
	case ContextFlowMatrix:
		switch key {
		case "j", "k", "tab", "enter", "esc", "f":
			return false
		}
	case ContextActionable:
		switch key {
		case "j", "k", "enter", "a", "?":
			return false
		}
	case ContextFilter:
		switch key {
		case "esc", "ctrl+s", "enter":
			return false
		}
		return true
	case ContextContextHelp:
		switch key {
		case "esc", "q", "~":
			return false
		}
		return true
	case ContextTutorial:
		switch key {
		case "esc", "q", "t", "tab", "l", "h", "j", "k", " ", "left", "right", "ctrl+d", "ctrl+u":
			return false
		}
		return true
	case ContextGlyphHelp:
		switch key {
		case "j", "k", "esc", "K":
			return false
		}
		return true
	case ContextQuitConfirm:
		switch key {
		case "esc", "y", "Y":
			return false
		}
		return true
	case ContextUpdateModal:
		switch key {
		case "esc", "q", "enter", "n", "N":
			return false
		}
		return true
	case ContextHistorySearch:
		switch key {
		case "esc", "enter":
			return false
		}
		return true
	case ContextAlerts:
		switch key {
		case "j", "k", "enter", "d", "!", "esc", "q":
			return false
		}
		return true
	case ContextAgentPrompt:
		switch key {
		case "esc", "q", "enter", "y", "n":
			return false
		}
		return true
	case ContextCassSession:
		switch key {
		case "V", "esc", "enter", "q":
			return false
		}
		return true
	case ContextLabelHealthDetail:
		switch key {
		case "esc", "q", "enter", "h", "d":
			return false
		}
		return true
	case ContextLabelDrilldown, ContextLabelGraphAnalysis:
		switch key {
		case "esc", "q", "g", "enter", "d":
			return false
		}
		return true
	case ContextHelp:
		switch key {
		case "j", "k", "q", " ", "esc", "x", "g", "G", "ctrl+d", "ctrl+u", "pgup", "pgdown":
			return false
		}
		return true
	case ContextLabelPicker, ContextRecipePicker, ContextRepoPicker:
		switch key {
		case "j", "k", "enter", "esc", "space", "q":
			return false
		}
		return true
	case ContextTimeTravelInput:
		switch key {
		case "enter", "esc":
			return false
		}
		return true
	case ContextSplit, ContextDetail, ContextTimeTravel, ContextList:
		switch key {
		case "enter", "t", "S", "l", "L", "ctrl+r", "f5", "K", "?", "w", "tab", "C", "O", "x", "esc":
			return false
		}
	}

	switch key {
	case "j", "k", "up", "down", "left", "right",
		"ctrl+d", "ctrl+u", "home", "gg", "wheel":
		return true
	default:
		return false
	}
}

func footerKeyPriority(ctx Context, key string) int {
	switch ctx {
	case ContextHistory:
		switch key {
		case "h", "q", "esc":
			return 0
		case "tab", "enter":
			return 1
		}
	case ContextInsights, ContextAttention:
		switch key {
		case "h", "l", "e", "enter", "?", "]", "f4", "f":
			return 0
		}
	case ContextFilter:
		switch key {
		case "esc", "ctrl+s", "enter":
			return 0
		}
	case ContextTutorial:
		switch key {
		case "esc", "q", "l", "h", "j", "k", "t", " ":
			return 0
		}
	case ContextContextHelp:
		switch key {
		case "esc", "q", "~":
			return 0
		}
	case ContextQuitConfirm:
		switch key {
		case "esc", "y", "Y":
			return 0
		}
	case ContextGraph:
		switch key {
		case "h", "l", "j", "k", "enter", "g":
			return 0
		}
	case ContextBoard, ContextBoardSearch:
		switch key {
		case "h", "l", "j", "k", "G", "enter", "b":
			return 0
		}
	case ContextList, ContextSplit, ContextDetail, ContextTimeTravel:
		switch key {
		case "enter", "?", "l", "L", "t", "S", "K", "ctrl+r", "f5", "w":
			return 0
		}
	}
	return 50
}

func footerCategoryPriority(ctx Context, category string) int {
	switch ctx {
	case ContextInsights, ContextAttention:
		switch category {
		case "Insights":
			return 0
		case "Views":
			return 1
		case "Actions":
			return 2
		case "Help":
			return 3
		}
	case ContextHistory:
		if category == "History" {
			return 0
		}
	case ContextGraph:
		if category == "Graph" {
			return 0
		}
	case ContextBoard, ContextBoardSearch:
		if category == "Board" {
			return 0
		}
	case ContextFlowMatrix:
		if category == "Flow" {
			return 0
		}
	case ContextFilter:
		if category == "Filter" {
			return 0
		}
	case ContextTutorial:
		if category == "Tutorial" {
			return 0
		}
	}

	switch category {
	case "Actions":
		return 10
	case "Filters", "Filter":
		return 20
	case "Views":
		return 30
	case "Help":
		return 40
	case "Navigation":
		return 50
	default:
		return 60
	}
}

func docContextSpecificity(doc KeyBindingDoc, docContexts []string) int {
	best := 0
	for _, raw := range strings.Split(doc.Context, ",") {
		raw = strings.TrimSpace(raw)
		for i, want := range docContexts {
			if raw != want {
				continue
			}
			score := len(docContexts) - i
			if raw == "all" {
				score = 1
			}
			if score > best {
				best = score
			}
		}
	}
	return best
}

func footerSkipBinding(ctx Context, doc KeyBindingDoc) bool {
	if ctx == ContextInsights || ctx == ContextAttention {
		switch doc.Key {
		case "h", "b", "g", "a", "i", "E":
			if doc.Category == "Views" {
				return true
			}
		}
	}
	// O (open in $EDITOR) is handled only on focusList/focusDetail in Update().
	switch ctx {
	case ContextFlowMatrix, ContextActionable:
		if doc.Key == "O" && doc.Category == "Actions" {
			return true
		}
	}
	return false
}

// HintsFor returns bindings from the authoritative docs filtered for ctx and surface.
func (r *KeyRegistry) HintsFor(ctx Context, surface HintSurface, limit int) []KeyBinding {
	_ = r // registry reserved for future handler-aware hint filtering
	docContexts := docContextsForUI(ctx, surface)
	if len(docContexts) == 0 {
		return nil
	}

	type scored struct {
		b     KeyBinding
		score int
	}
	scoredBindings := make(map[string]scored)

	for _, doc := range GetKeyBindingDocs() {
		if !docAppliesToContexts(doc, docContexts) {
			continue
		}
		if surface == HintFooter && footerExcludedKey(ctx, doc.Key) {
			continue
		}
		if surface == HintFooter && footerSkipBinding(ctx, doc) {
			continue
		}
		specificity := docContextSpecificity(doc, docContexts)
		candidate := scored{
			b: KeyBinding{
				Key:      doc.Key,
				Desc:     doc.Desc,
				Category: doc.Category,
			},
			score: specificity,
		}
		if existing, ok := scoredBindings[doc.Key]; !ok || candidate.score > existing.score {
			scoredBindings[doc.Key] = candidate
		}
	}

	result := make([]KeyBinding, 0, len(scoredBindings))
	for _, s := range scoredBindings {
		result = append(result, s.b)
	}

	if surface == HintFooter && (ctx == ContextInsights || ctx == ContextAttention) {
		result = collapseAttentionFooterKeys(result)
	}

	sort.SliceStable(result, func(i, j int) bool {
		pi := footerKeyPriority(ctx, result[i].Key)
		pj := footerKeyPriority(ctx, result[j].Key)
		if pi != pj {
			return pi < pj
		}
		ci := footerCategoryPriority(ctx, result[i].Category)
		cj := footerCategoryPriority(ctx, result[j].Category)
		if ci != cj {
			return ci < cj
		}
		return result[i].Key < result[j].Key
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}

func collapseAttentionFooterKeys(bindings []KeyBinding) []KeyBinding {
	hasBracket := false
	hasF4 := false
	for _, b := range bindings {
		switch b.Key {
		case "]", "]/F4":
			hasBracket = true
		case "f4":
			hasF4 = true
		}
	}
	filtered := bindings[:0]
	for _, b := range bindings {
		if b.Key == "f4" {
			if hasBracket {
				continue
			}
			b.Key = "]/F4"
			hasBracket = true
		}
		if b.Key == "]" {
			b.Key = "]/F4"
		}
		filtered = append(filtered, b)
	}
	if !hasBracket && hasF4 {
		// ponytail: f4 alone still advertises the combined chord label
		for i := range filtered {
			if filtered[i].Key == "f4" {
				filtered[i].Key = "]/F4"
			}
		}
	}
	return filtered
}

func formatKeyForHint(key string) string {
	switch strings.ToLower(key) {
	case "enter":
		return "⏎"
	case "ctrl+r":
		return "Ctrl+R"
	case "ctrl+s":
		return "ctrl+s"
	case "f4":
		return "F4"
	default:
		return key
	}
}

func footerDescLabel(ctx Context, key, desc string) string {
	if label, ok := footerKeyLabels[key]; ok {
		if key == "enter" {
			switch ctx {
			case ContextFilter:
				return "select"
			case ContextHistory:
				return "jump"
			case ContextLabelPicker, ContextRecipePicker, ContextRepoPicker:
				return "apply"
			}
		}
		if key == "esc" {
			switch ctx {
			case ContextHistory:
				return "close"
			case ContextDetail, ContextSplit:
				return "back"
			case ContextTutorial:
				return "close"
			}
		}
		if key == "l" && ctx == ContextTutorial {
			return "next"
		}
		if key == "h" && ctx == ContextTutorial {
			return "prev"
		}
		if key == "h" && ctx == ContextHistory {
			return "close"
		}
		return label
	}

	desc = strings.TrimSpace(desc)
	if desc == "" {
		return ""
	}
	words := strings.Fields(desc)
	if len(words) == 0 {
		return desc
	}
	label := strings.ToLower(words[0])
	if len(words) > 1 && (words[0] == "Go" || words[0] == "Toggle" || words[0] == "Open") {
		label = strings.ToLower(words[1])
	}
	return label
}

var footerKeyLabels = map[string]string{
	"enter":  "details",
	"t":      "diff",
	"S":      "triage",
	"l":      "labels",
	"L":      "labels",
	"K":      "symbols",
	"?":      "help",
	"f":      "flow",
	"]":      "attention",
	"f4":     "attention",
	"e":      "explain",
	"ctrl+r": "refresh",
	"f5":     "refresh",
	"C":      "copy",
	"O":      "edit",
	"x":      "export",
	"w":      "repos",
	"tab":    "focus",
	"esc":    "close",
	"q":      "close",
	"`":      "tutorial",
	"~":      "close",
	"g":      "list",
	"b":      "list",
	"a":      "list",
	"space":  "toggle",
}

func formatFooterHint(ctx Context, key, desc string, keyStyle lipgloss.Style) string {
	label := footerDescLabel(ctx, key, desc)
	if label == "" {
		return keyStyle.Render(formatKeyForHint(key))
	}
	return keyStyle.Render(formatKeyForHint(key)) + " " + label
}

func expandFooterHintKeys(key string) []string {
	var keys []string
	if strings.Contains(key, "/") {
		for _, part := range strings.Split(key, "/") {
			part = strings.TrimSpace(part)
			if part != "" {
				keys = append(keys, part)
			}
		}
	} else if key == "⏎" {
		keys = []string{"enter"}
	} else {
		keys = []string{key}
	}
	for i, k := range keys {
		if strings.EqualFold(k, "f4") {
			keys[i] = "f4"
		}
	}
	return keys
}

func (m Model) footerHintBindings() []KeyBinding {
	switch {
	case m.showHelp:
		return nil
	case m.showGlyphHelp:
		return []KeyBinding{
			{Key: "j", Desc: "Scroll down", Category: "Help"},
			{Key: "k", Desc: "Scroll up", Category: "Help"},
			{Key: "K", Desc: "Close symbol reference", Category: "Help"},
			{Key: "esc", Desc: "Close symbol reference", Category: "Help"},
		}
	case m.isBoardView && m.board.IsSearchMode():
		return []KeyBinding{
			{Key: "n", Desc: "Next match", Category: "Board"},
			{Key: "N", Desc: "Previous match", Category: "Board"},
			{Key: "enter", Desc: "Finish board search", Category: "Board"},
			{Key: "esc", Desc: "Cancel board search", Category: "Board"},
		}
	default:
		if m.keyRegistry == nil {
			return nil
		}
		return m.keyRegistry.HintsFor(m.footerSubjectContext(), HintFooter, footerHintLimit)
	}
}

func (m Model) footerAdvertisedKeys() []string {
	var keys []string
	for _, b := range m.footerHintBindings() {
		keys = append(keys, expandFooterHintKeys(b.Key)...)
	}
	return keys
}

func (m *Model) footerHintsFromRegistry(keyStyle lipgloss.Style) []string {
	bindings := m.footerHintBindings()
	ctx := m.footerSubjectContext()
	hints := make([]string, 0, len(bindings))
	for _, b := range bindings {
		hints = append(hints, formatFooterHint(ctx, b.Key, b.Desc, keyStyle))
	}
	return hints
}

// footerShowsListLabelHint reports whether the left labelHint strip should show
// list-centric shortcuts (l:labels, enter:detail). Overlays and specialized
// views use their own hint surfaces instead.
func (m Model) footerShowsListLabelHint() bool {
	switch m.footerSubjectContext() {
	case ContextList, ContextSplit, ContextDetail, ContextTimeTravel, ContextBoard, ContextAttention:
		return true
	default:
		return false
	}
}

func (m *Model) helpOverlayPanelSpecs() []helpOverlayPanelSpec {
	ctx := m.helpSubjectContext()

	always := []helpOverlayPanelSpec{
		{title: "Navigation", icon: icons.Navigation, colorIdx: 0, category: "Navigation"},
		{title: "Global", icon: icons.Globe, colorIdx: 2, category: "Help"},
	}

	var contextual []helpOverlayPanelSpec

	switch ctx {
	case ContextGraph:
		contextual = append(contextual, helpOverlayPanelSpec{title: "Graph View", icon: icons.Chart, colorIdx: 4, category: "Graph"})
	case ContextBoard, ContextBoardSearch:
		contextual = append(contextual, helpOverlayPanelSpec{title: "Board View", icon: icons.Task, colorIdx: 4, category: "Board"})
	case ContextInsights, ContextAttention:
		contextual = append(contextual, helpOverlayPanelSpec{title: "Insights", icon: icons.Lightbulb, colorIdx: 5, category: "Insights"})
	case ContextHistory:
		contextual = append(contextual, helpOverlayPanelSpec{title: "History", icon: icons.HistoryScroll, colorIdx: 0, category: "History"})
	case ContextFlowMatrix:
		contextual = append(contextual, helpOverlayPanelSpec{title: "Flow Matrix", icon: icons.DepDiscovered, colorIdx: 3, category: "Flow"})
	case ContextLabelDashboard:
		contextual = append(contextual, helpOverlayPanelSpec{title: "Labels", icon: icons.Label, colorIdx: 1, category: "Labels"})
	case ContextActionable:
		contextual = append(contextual, helpOverlayPanelSpec{title: "Actionable", icon: icons.Epic, colorIdx: 1, category: "Actionable"})
	}

	if ctx.AllowsGlobalFallthrough() || ctx == ContextList || ctx == ContextSplit || ctx == ContextDetail || ctx == ContextFilter || ctx == ContextTimeTravel {
		contextual = append(contextual,
			helpOverlayPanelSpec{title: "Views", icon: icons.Eye, colorIdx: 1, category: "Views"},
			helpOverlayPanelSpec{title: "Filters & Sort", icon: icons.DepDiscovered, colorIdx: 3, category: "Filters"},
			helpOverlayPanelSpec{title: "Actions", icon: icons.Lightning, colorIdx: 1, category: "Actions"},
		)
	}

	if ctx == ContextFilter {
		contextual = append(contextual, helpOverlayPanelSpec{title: "Filter Input", icon: icons.DepDiscovered, colorIdx: 3, category: "Filter"})
	}

	statusPanel := helpOverlayPanelSpec{
		title:    "Status",
		icon:     icons.Health,
		colorIdx: 2,
		static: []keyDescPair{
			{"◌ metrics", "Phase 2 metrics computing"},
			{"⚠ age", "Snapshot getting stale"},
			{"⚠ STALE", "Snapshot is stale"},
			{"✗ bg", "Background worker errors"},
			{"↻ recov", "Worker self-healed"},
			{"⚠ dead", "Worker unresponsive"},
			{"polling", "Live reload uses polling"},
		},
	}

	panels := append(always, contextual...)
	panels = append(panels, statusPanel)
	return panels
}

func (m *Model) helpOverlayShortcuts(category string, subject Context) []keyDescPair {
	if m.keyRegistry == nil || category == "" {
		return nil
	}

	// Merge subject + global doc contexts for the requested category.
	docContexts := docContextsForUI(subject, HintHelpOverlay)
	var pairs []keyDescPair
	seen := make(map[string]struct{})

	for _, doc := range GetKeyBindingDocs() {
		if doc.Category != category {
			continue
		}
		if !docAppliesToContexts(doc, docContexts) {
			continue
		}
		if _, ok := seen[doc.Key]; ok {
			continue
		}
		seen[doc.Key] = struct{}{}
		pairs = append(pairs, keyDescPair{key: doc.Key, desc: doc.Desc})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].key < pairs[j].key
	})
	return pairs
}
