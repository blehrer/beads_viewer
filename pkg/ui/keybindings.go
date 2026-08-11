package ui

import (
	"sort"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// KeyHandler processes a key event and returns the updated model and whether
// the key was handled. If handled is false, the key may fall through to other
// handlers (e.g., cross-view switches).
type KeyHandler func(m Model, msg tea.KeyMsg) (Model, bool)

// KeyBinding associates a key with a handler for a specific focus context.
type KeyBinding struct {
	Focus            focus      // Which view/focus context this binding applies to
	Key              string     // Key string (e.g., "j", "ctrl+d", "enter")
	Desc             string     // Human-readable description for help display
	Category         string     // Grouping category (e.g., "Navigation", "Actions")
	Handler          KeyHandler // The handler function to call
	AllowFallthrough bool       // If true, continue down the stack when Handler returns handled=false
}

// KeyRegistry manages key bindings organized by focus context. It provides
// a centralized dispatch mechanism for key events.
type KeyRegistry struct {
	mu       sync.RWMutex
	handlers map[focus]map[string]KeyHandler // focus -> key -> handler
	bindings map[focus][]KeyBinding          // focus -> ordered bindings (for help)
}

// NewKeyRegistry creates an empty key registry ready to accept bindings.
func NewKeyRegistry() *KeyRegistry {
	return &KeyRegistry{
		handlers: make(map[focus]map[string]KeyHandler),
		bindings: make(map[focus][]KeyBinding),
	}
}

// RegisterBinding adds a single key binding to the registry.
// If the same key is registered twice for the same focus, the later
// registration overwrites the earlier one.
func (r *KeyRegistry) RegisterBinding(b KeyBinding) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure the focus map exists
	if r.handlers[b.Focus] == nil {
		r.handlers[b.Focus] = make(map[string]KeyHandler)
	}

	// Register the handler (doc-only bindings omit Handler)
	if b.Handler != nil {
		r.handlers[b.Focus][b.Key] = b.Handler
	}

	// Track binding for help generation (replace if exists)
	existingBindings := r.bindings[b.Focus]
	found := false
	for i, existing := range existingBindings {
		if existing.Key == b.Key {
			existingBindings[i] = b
			found = true
			break
		}
	}
	if !found {
		r.bindings[b.Focus] = append(r.bindings[b.Focus], b)
	}
}

// RegisterView adds multiple bindings for a specific focus context.
// This is a convenience wrapper around RegisterBinding.
func (r *KeyRegistry) RegisterView(f focus, bindings []KeyBinding) {
	for _, b := range bindings {
		// Ensure the binding's focus matches the specified focus
		b.Focus = f
		r.RegisterBinding(b)
	}
}

// Dispatch looks up and executes the handler for the given focus and key.
// Returns:
//   - model: The potentially updated model
//   - handled: True if a handler was found and executed
//   - cmd: Any tea.Cmd returned by the handler (currently always nil; reserved for future use)
//
// If no handler is registered for the focus+key combination, handled is false
// and the model is returned unchanged.
func (r *KeyRegistry) Dispatch(f focus, key string, m Model, msg tea.KeyMsg) (Model, bool, tea.Cmd) {
	r.mu.RLock()
	focusHandlers := r.handlers[f]
	var handler KeyHandler
	if focusHandlers != nil {
		handler = focusHandlers[key]
	}
	r.mu.RUnlock()

	if handler == nil {
		return m, false, nil
	}

	updatedModel, handled := handler(m, msg)
	return updatedModel, handled, nil
}

// DispatchStack walks the focus stack until a handler reports handled=true.
// Bindings with AllowFallthrough continue to the next focus when handled=false.
func (r *KeyRegistry) DispatchStack(stack []focus, key string, m Model, msg tea.KeyMsg) (Model, bool, tea.Cmd) {
	for i, f := range stack {
		updated, handled, cmd := r.Dispatch(f, key, m, msg)
		if handled {
			return updated, true, cmd
		}
		m = updated
		if i < len(stack)-1 {
			binding := r.lookupBinding(f, key)
			if binding != nil && binding.Handler != nil && !binding.AllowFallthrough {
				break
			}
			continue
		}
	}
	return m, false, nil
}

func (r *KeyRegistry) lookupBinding(f focus, key string) *KeyBinding {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.bindings[f] {
		if b.Key == key {
			copy := b
			return &copy
		}
	}
	return nil
}

// AllBindingsForFocus returns all registered bindings for a specific focus,
// sorted by category then by key. Returns an empty slice if no bindings exist.
func (r *KeyRegistry) AllBindingsForFocus(f focus) []KeyBinding {
	r.mu.RLock()
	bindings := r.bindings[f]
	if len(bindings) == 0 {
		r.mu.RUnlock()
		return []KeyBinding{}
	}

	// Return a sorted copy to avoid exposing internal state
	result := make([]KeyBinding, len(bindings))
	copy(result, bindings)
	r.mu.RUnlock()

	sort.Slice(result, func(i, j int) bool {
		if result[i].Category != result[j].Category {
			return result[i].Category < result[j].Category
		}
		return result[i].Key < result[j].Key
	})

	return result
}

// AllBindings returns all registered bindings across all focus contexts,
// sorted by focus, then category, then key.
func (r *KeyRegistry) AllBindings() []KeyBinding {
	r.mu.RLock()
	if len(r.bindings) == 0 {
		r.mu.RUnlock()
		return []KeyBinding{}
	}

	var result []KeyBinding
	for _, bindings := range r.bindings {
		result = append(result, bindings...)
	}
	r.mu.RUnlock()

	sort.Slice(result, func(i, j int) bool {
		if result[i].Focus != result[j].Focus {
			return result[i].Focus < result[j].Focus
		}
		if result[i].Category != result[j].Category {
			return result[i].Category < result[j].Category
		}
		return result[i].Key < result[j].Key
	})

	return result
}

// HasHandler reports whether a runtime handler (not doc-only) exists for focus+key.
func (r *KeyRegistry) HasHandler(f focus, key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if focusHandlers := r.handlers[f]; focusHandlers != nil {
		_, exists := focusHandlers[key]
		return exists
	}
	return false
}

// HasBinding checks if a binding exists for the given focus and key.
func (r *KeyRegistry) HasBinding(f focus, key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if focusHandlers := r.handlers[f]; focusHandlers != nil {
		if _, exists := focusHandlers[key]; exists {
			return true
		}
	}
	for _, b := range r.bindings[f] {
		if b.Key == key {
			return true
		}
	}
	return false
}

// BindingsCount returns the total number of registered bindings.
func (r *KeyRegistry) BindingsCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, bindings := range r.bindings {
		count += len(bindings)
	}
	return count
}

// Clear removes all registered bindings. Primarily useful for testing.
func (r *KeyRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers = make(map[focus]map[string]KeyHandler)
	r.bindings = make(map[focus][]KeyBinding)
}

// registerKeyBindings populates the KeyRegistry with all view handlers.
// Called from NewModel to set up two-phase key dispatch (bv-3bsx).
// For now this registers the authoritative documentation bindings so the
// shortcuts/help surfaces stay in sync even before runtime dispatch migrates.
func (m *Model) registerKeyBindings() {
	if m == nil || m.keyRegistry == nil {
		return
	}

	for _, doc := range GetKeyBindingDocs() {
		for _, f := range focusesForBindingDoc(doc) {
			m.keyRegistry.RegisterBinding(KeyBinding{
				Focus:    f,
				Key:      doc.Key,
				Desc:     doc.Desc,
				Category: doc.Category,
			})
		}
	}
}

func focusesForBindingDoc(doc KeyBindingDoc) []focus {
	var focuses []focus
	seen := make(map[focus]struct{})

	addFocus := func(f focus) {
		if _, ok := seen[f]; ok {
			return
		}
		seen[f] = struct{}{}
		focuses = append(focuses, f)
	}

	for _, raw := range strings.Split(doc.Context, ",") {
		switch strings.TrimSpace(raw) {
		case "all":
			for _, f := range allDocumentedFocuses() {
				addFocus(f)
			}
		case "list":
			addFocus(focusList)
		case "detail":
			addFocus(focusDetail)
		case "board":
			addFocus(focusBoard)
		case "board-search":
			addFocus(focusBoard)
		case "graph":
			addFocus(focusGraph)
		case "insights":
			addFocus(focusInsights)
		case "attention":
			addFocus(focusInsights)
		case "history":
			addFocus(focusHistory)
		case "actionable":
			addFocus(focusActionable)
		case "label", "label-dashboard":
			addFocus(focusLabelDashboard)
		case "label-picker":
			addFocus(focusLabelPicker)
		case "recipe-picker":
			addFocus(focusRecipePicker)
		case "repo-picker":
			addFocus(focusRepoPicker)
		case "tree":
			addFocus(focusTree)
		case "flow", "flow-matrix":
			addFocus(focusFlowMatrix)
		case "sprint":
			addFocus(focusSprint)
		case "filter":
			addFocus(focusList)
		case "alerts":
			addFocus(focusList)
		case "help":
			addFocus(focusHelp)
		case "context-help":
			addFocus(focusContextHelp)
		}
	}

	return focuses
}

func allDocumentedFocuses() []focus {
	return []focus{
		focusList,
		focusDetail,
		focusBoard,
		focusGraph,
		focusInsights,
		focusHistory,
		focusActionable,
		focusLabelDashboard,
		focusLabelPicker,
		focusRecipePicker,
		focusRepoPicker,
		focusTree,
		focusFlowMatrix,
		focusSprint,
		focusHelp,
		focusContextHelp,
	}
}

// KeyBindingCasePolicy documents lowercase vs uppercase semantics for TUI keys.
// Lowercase keys are movement, view-local actions, or toggles; uppercase keys
// are alternate actions (scroll, column jump, hybrid search). Single-letter
// view switches (b/g/h/f/i/a/E) are lowercase only.
const KeyBindingCasePolicy = `TUI key case policy:
- lowercase: movement, view-local navigation, filters, and view toggles
- uppercase: alternate actions (board column jump H/L, hybrid toggle H, sort/triage S, scroll J/K in history)
- single-letter view switches are lowercase only (b board, g graph, h history, f flow, i insights, E tree)`

// KeyBindingDoc represents a key binding for documentation purposes (bv-xl6g).
type KeyBindingDoc struct {
	Key      string
	Desc     string
	Category string
	Context  string // Which view(s) this applies to
}

// GetKeyBindingDocs returns all key bindings for documentation/robot-help (bv-xl6g).
// This is separate from the registry to allow documentation even before handlers
// are registered (bv-3bsx migration).
func GetKeyBindingDocs() []KeyBindingDoc {
	return []KeyBindingDoc{
		// Global navigation
		{"j", "Move down", "Navigation", "all"},
		{"k", "Move up", "Navigation", "all"},
		{"G", "Go to end", "Navigation", "all"},
		{"home", "Go to start", "Navigation", "list,detail,board,graph,tree,actionable,history,flow-matrix,insights"},
		{"gg", "Go to start (combo)", "Navigation", "board,tree"},
		{"ctrl+d", "Page down", "Navigation", "all"},
		{"ctrl+u", "Page up", "Navigation", "all"},
		{"enter", "Open/select", "Navigation", "all"},
		{"esc", "Back/close", "Navigation", "all"},
		{"q", "Quit or close view", "Navigation", "all"},
		{"tab", "Toggle split focus", "Navigation", "list,detail"},
		{"<", "Shrink list pane", "Navigation", "list,detail"},
		{">", "Expand list pane", "Navigation", "list,detail"},
		{"wheel", "Scroll focused pane", "Mouse", "all"},

		// Help & reference overlays
		{"?", "Help overlay", "Help", "all"},
		{"f1", "Help overlay", "Help", "all"},
		{"~", "Context help", "Help", "all"},
		{"`", "Interactive tutorial", "Help", "all"},
		{";", "Shortcuts sidebar", "Help", "all"},
		{"f2", "Shortcuts sidebar", "Help", "all"},
		{"K", "Symbol reference (glyph glossary)", "Help", "list,detail,board,graph,insights,actionable,tree,flow-matrix,label-dashboard"},
		{"ctrl+j", "Scroll shortcuts sidebar down", "Help", "all"},
		{"ctrl+k", "Scroll shortcuts sidebar up", "Help", "all"},

		// View switching (lowercase only — see KeyBindingCasePolicy)
		{"a", "Actionable view", "Views", "list,detail"},
		{"b", "Board view", "Views", "list,detail"},
		{"g", "Graph view", "Views", "list,detail"},
		{"h", "History view", "Views", "list,detail"},
		{"i", "Insights panel", "Views", "list,detail"},
		{"E", "Tree view", "Views", "list,detail"},
		{"f", "Flow matrix view", "Views", "list,detail"},
		{"p", "Priority hints", "Views", "list,detail"},
		{"[", "Label dashboard", "Views", "list,detail"},
		{"f3", "Label dashboard", "Views", "list,detail"},
		{"]", "Attention view", "Views", "list,detail"},
		{"f4", "Attention view", "Views", "list,detail"},

		// List filters & search
		{"o", "Open issues only", "Filters", "list,board"},
		{"c", "Closed issues only", "Filters", "list,board"},
		{"r", "Ready (unblocked)", "Filters", "list,board"},
		{"l", "Label picker", "Filters", "list,detail"},
		{"L", "Label picker (Shift+L)", "Filters", "list,detail"},
		{"/", "Search/filter", "Filters", "list,history"},
		{"ctrl+s", "Toggle semantic search", "Filters", "list"},
		{"H", "Toggle hybrid search", "Filters", "list"},
		{"alt+h", "Cycle hybrid preset", "Filters", "list"},
		{"s", "Cycle sort mode", "Filters", "list"},
		{"S", "Apply triage recipe sort", "Filters", "list"},

		// Filter-mode input (list search active)
		{"esc", "Cancel filter", "Filter", "filter"},
		{"ctrl+s", "Toggle semantic while filtering", "Filter", "filter"},
		{"enter", "Apply filter", "Filter", "filter"},

		// Actions
		{"t", "Time travel prompt", "Actions", "list,detail"},
		{"T", "Time travel HEAD~5", "Actions", "list,detail"},
		{"x", "Export to markdown", "Actions", "list,detail"},
		{"y", "Copy issue ID", "Actions", "all"},
		{"C", "Copy full issue", "Actions", "list,detail"},
		{"O", "Open in $EDITOR", "Actions", "list,detail"},
		{"'", "Recipe picker", "Actions", "list"},
		{"U", "Self-update check", "Actions", "list"},
		{"V", "Cass sessions", "Actions", "list"},
		{"!", "Toggle alerts panel", "Actions", "list,detail"},
		{"w", "Repo picker (workspace)", "Actions", "list"},
		{"ctrl+r", "Force refresh", "Actions", "all"},
		{"f5", "Force refresh", "Actions", "all"},

		// Graph view
		{"h", "Move left", "Graph", "graph"},
		{"l", "Move right", "Graph", "graph"},
		{"j", "Move down", "Graph", "graph"},
		{"k", "Move up", "Graph", "graph"},
		{"pgup", "Page up", "Graph", "graph"},
		{"pgdown", "Page down", "Graph", "graph"},

		// Board view
		{"h", "Previous column", "Board", "board"},
		{"l", "Next column", "Board", "board"},
		{"H", "First column", "Board", "board"},
		{"L", "Last column", "Board", "board"},
		{"1", "Jump to Open column", "Board", "board"},
		{"2", "Jump to In Progress column", "Board", "board"},
		{"3", "Jump to Blocked column", "Board", "board"},
		{"4", "Jump to Closed column", "Board", "board"},
		{"0", "First card in column", "Board", "board"},
		{"$", "Last card in column", "Board", "board"},
		{"tab", "Toggle detail panel", "Board", "board"},
		{"ctrl+j", "Scroll detail down", "Board", "board"},
		{"ctrl+k", "Scroll detail up", "Board", "board"},
		{"/", "Board search", "Board", "board"},
		{"n", "Next search match", "Board", "board,board-search"},
		{"N", "Previous search match", "Board", "board,board-search"},
		{"s", "Cycle swimlane mode", "Board", "board"},
		{"e", "Toggle empty columns", "Board", "board"},
		{"d", "Toggle card expand", "Board", "board"},

		// Board search submode
		{"esc", "Cancel board search", "Board", "board-search"},
		{"enter", "Finish board search", "Board", "board-search"},
		{"backspace", "Delete search char", "Board", "board-search"},

		// Tree view
		{"h", "Collapse/parent", "Tree", "tree"},
		{"l", "Expand/child", "Tree", "tree"},
		{" ", "Toggle expand", "Tree", "tree"},
		{"o", "Expand all", "Tree", "tree"},
		{"O", "Collapse all", "Tree", "tree"},
		{"E", "Close tree view", "Tree", "tree"},
		{"tab", "Toggle detail (split)", "Tree", "tree"},

		// Insights view
		{"h", "Previous panel", "Insights", "insights,attention"},
		{"l", "Next panel", "Insights", "insights,attention"},
		{"e", "Toggle explanations", "Insights", "insights,attention"},
		{"x", "Calculation proof", "Insights", "insights,attention"},
		{"m", "Heatmap toggle", "Insights", "insights,attention"},
		{"ctrl+j", "Scroll detail down", "Insights", "insights,attention"},
		{"ctrl+k", "Scroll detail up", "Insights", "insights,attention"},

		// Flow matrix view
		{"f", "Close flow matrix", "Flow", "flow-matrix"},
		{"tab", "Toggle panel", "Flow", "flow-matrix"},
		{"g", "Go to start", "Flow", "flow-matrix"},
		{"G", "Go to end", "Flow", "flow-matrix"},

		// History view
		{"v", "Toggle git/bead mode", "History", "history"},
		{"tab", "Cycle focus panes", "History", "history"},
		{"J", "Detail scroll down", "History", "history"},
		{"K", "Detail scroll up", "History", "history"},
		{"f", "Toggle file tree", "History", "history"},
		{"F", "Toggle file tree", "History", "history"},
		{"g", "Jump to graph for bead", "History", "history"},
		{"y", "Copy commit SHA", "History", "history"},
		{"c", "Cycle confidence filter", "History", "history"},
		{"o", "Open commit in browser", "History", "history"},
		{"h", "Close history view", "History", "history"},

		// Label dashboard
		{"h", "Label health detail", "Labels", "label-dashboard"},
		{"d", "Label drilldown", "Labels", "label-dashboard"},
		{"enter", "Filter list by label", "Labels", "label-dashboard"},
		{"esc", "Close label dashboard", "Labels", "label-dashboard"},

		// Label picker overlay
		{"esc", "Cancel label picker", "Labels", "label-picker"},
		{"enter", "Apply label filter", "Labels", "label-picker"},

		// Recipe picker overlay
		{"esc", "Close recipe picker", "Recipes", "recipe-picker"},
		{"q", "Close recipe picker", "Recipes", "recipe-picker"},
		{"enter", "Apply recipe", "Recipes", "recipe-picker"},

		// Repo picker overlay (workspace)
		{" ", "Toggle repo", "Workspace", "repo-picker"},
		{"a", "Select all repos", "Workspace", "repo-picker"},
		{"enter", "Apply repo filter", "Workspace", "repo-picker"},
		{"esc", "Close repo picker", "Workspace", "repo-picker"},

		// Alerts panel
		{"j", "Next alert", "Alerts", "alerts"},
		{"k", "Previous alert", "Alerts", "alerts"},
		{"enter", "Jump to issue", "Alerts", "alerts"},
		{"d", "Dismiss alert", "Alerts", "alerts"},
		{"!", "Close alerts panel", "Alerts", "alerts"},

		// Sprint view (when active)
		{"P", "Close sprint view", "Sprints", "sprint"},
		{"j", "Next sprint", "Sprints", "sprint"},
		{"k", "Previous sprint", "Sprints", "sprint"},

		// Help overlay scroll/dismiss
		{"j", "Scroll help down", "Help", "help"},
		{"k", "Scroll help up", "Help", "help"},
		{" ", "Open tutorial from help", "Help", "help"},
		{"q", "Close help", "Help", "help"},

		// Context help overlay
		{"esc", "Close context help", "Help", "context-help"},
		{"q", "Close context help", "Help", "context-help"},
		{"~", "Close context help", "Help", "context-help"},

		// Actionable view
		{"j", "Move down", "Actionable", "actionable"},
		{"k", "Move up", "Actionable", "actionable"},
	}
}

// KeyBindingDocsForRobot exports authoritative bindings for robot-help JSON.
func KeyBindingDocsForRobot() []map[string]string {
	docs := GetKeyBindingDocs()
	out := make([]map[string]string, len(docs))
	for i, doc := range docs {
		out[i] = map[string]string{
			"key":      doc.Key,
			"desc":     doc.Desc,
			"category": doc.Category,
			"context":  doc.Context,
		}
	}
	return out
}
