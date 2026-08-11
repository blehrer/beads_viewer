package ui

import (
	"fmt"
	"time"

	"github.com/Dicklesworthstone/beads_viewer/pkg/analysis"
	tea "github.com/charmbracelet/bubbletea"
)

// registerKeyHandlers attaches runtime KeyHandler funcs for registry dispatch (bv-p5kf.15).
func (m *Model) registerKeyHandlers() {
	if m == nil || m.keyRegistry == nil {
		return
	}
	r := m.keyRegistry

	register := func(f focus, key string, handler KeyHandler, allowFallthrough bool) {
		r.RegisterBinding(KeyBinding{
			Focus:            f,
			Key:              key,
			Handler:          handler,
			AllowFallthrough: allowFallthrough,
		})
	}

	// Board view keys (AllowFallthrough: view toggles fall through to focusList)
	boardKeys := []string{
		"h", "l", "j", "k", "left", "right", "up", "down",
		"home", "end", "G", "ctrl+d", "ctrl+u",
		"1", "2", "3", "4", "H", "L", "0", "$",
		"/", "n", "N", "y", "o", "c", "r", "s", "e", "d",
		"tab", "enter", "ctrl+j", "ctrl+k",
	}
	for _, k := range boardKeys {
		k := k
		register(focusBoard, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			model.pendingComboKey = ""
			return model.handleBoardKeys(msg), true
		}, true)
	}

	// Graph view
	graphKeys := []string{
		"h", "l", "j", "k", "left", "right", "up", "down",
		"H", "L", "ctrl+d", "ctrl+u", "pgup", "pgdown", "enter",
	}
	for _, k := range graphKeys {
		k := k
		register(focusGraph, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleGraphKeys(msg), true
		}, true)
	}

	// Tree view (g handled by gg-combo in tryGGCombo)
	treeKeys := []string{
		"h", "l", "j", "k", "left", "right", "up", "down",
		"G", "o", "O", "E", "esc", "enter", " ", "tab",
		"ctrl+d", "ctrl+u", "pgup", "pgdown",
	}
	for _, k := range treeKeys {
		k := k
		register(focusTree, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			model.pendingComboKey = ""
			return model.handleTreeKeys(msg), true
		}, true)
	}

	// Insights / attention panel
	insightsKeys := []string{
		"h", "l", "j", "k", "up", "down", "left", "right",
		"ctrl+j", "ctrl+k", "tab", "e", "x", "m", "enter", "esc",
	}
	for _, k := range insightsKeys {
		k := k
		register(focusInsights, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleInsightsKeys(msg), true
		}, true)
	}

	// Actionable view
	for _, k := range []string{"j", "k", "up", "down", "enter"} {
		k := k
		register(focusActionable, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleActionableKeys(msg), true
		}, true)
	}

	// History view
	historyKeys := []string{
		"h", "f", "g", "j", "k", "up", "down", "J", "K",
		"v", "tab", "enter", "y", "c", "F", "o", "/",
	}
	for _, k := range historyKeys {
		k := k
		register(focusHistory, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleHistoryKeys(msg), true
		}, true)
	}

	// Sprint view
	for _, k := range []string{"P", "esc", "j", "k", "up", "down"} {
		k := k
		register(focusSprint, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleSprintKeys(msg), true
		}, true)
	}

	// Flow matrix
	flowKeys := []string{
		"f", "g", "j", "k", "up", "down", "G", "end", "home", "tab", "enter", "esc", "q",
	}
	for _, k := range flowKeys {
		k := k
		register(focusFlowMatrix, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleFlowMatrixKeys(msg), true
		}, true)
	}

	// Label picker / recipe / repo overlays (no fallthrough)
	for _, k := range []string{"esc", "enter", "q"} {
		k := k
		register(focusLabelPicker, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleLabelPickerKeys(msg), true
		}, false)
	}
	for _, k := range []string{"esc", "q", "enter", "j", "k"} {
		k := k
		register(focusRecipePicker, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleRecipePickerKeys(msg), true
		}, false)
	}
	for _, k := range []string{"esc", "enter", " ", "a", "j", "k"} {
		k := k
		register(focusRepoPicker, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleRepoPickerKeys(msg), true
		}, false)
	}

	// Help overlays
	for _, k := range []string{"j", "k", "up", "down", "pgdown", "pgup", "ctrl+d", "ctrl+u", "g", "G", "q", " ", "x", "esc", "?"} {
		k := k
		register(focusHelp, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleHelpKeys(msg), true
		}, false)
	}
	for _, k := range []string{"esc", "q", "~"} {
		k := k
		register(focusContextHelp, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleContextHelpKeys(msg), true
		}, false)
	}

	// Detail pane: label picker keys only; other action keys fall through to focusList.
	register(focusDetail, "l", func(model Model, msg tea.KeyMsg) (Model, bool) {
		if model, ok := model.openLabelPicker(); ok {
			return model, true
		}
		return model, true
	}, false)
	register(focusDetail, "L", func(model Model, msg tea.KeyMsg) (Model, bool) {
		if model, ok := model.openLabelPicker(); ok {
			return model, true
		}
		return model, true
	}, false)

	// Global view toggles and list actions on focusList
	viewToggleKeys := []string{"b", "g", "a", "E", "i", "p", "h", "[", "f3", "]", "f4", "f", "!", "'", "w", "x", "l", "L"}
	for _, k := range viewToggleKeys {
		k := k
		register(focusList, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			return model.handleViewToggleKey(msg)
		}, false)
	}

	listKeys := []string{
		"enter", "home", "G", "end", "ctrl+d", "ctrl+u",
		"o", "c", "r", "t", "T", "C", "S", "s", "V", "U", "y",
	}
	for _, k := range listKeys {
		k := k
		register(focusList, k, func(model Model, msg tea.KeyMsg) (Model, bool) {
			if !model.listKeyHandled(k) {
				return model, false
			}
			return model.handleListKeys(msg), true
		}, false)
	}
}

// tryGGCombo handles the board/tree "gg" chord and single-g combo timer (bv-6fm0).
// Returns (updatedModel, comboCmd, handled).
func (m Model) tryGGCombo(msg tea.KeyMsg) (Model, tea.Cmd, bool) {
	if msg.String() != "g" {
		return m, nil, false
	}
	switch m.focused {
	case focusBoard:
		if m.pendingComboKey == "g" && m.pendingComboFocus == focusBoard && time.Since(m.pendingComboTime) < comboTimeout {
			m.board.MoveToTop()
			m.pendingComboKey = ""
			m.pendingComboTime = time.Time{}
			return m, nil, true
		}
		m.pendingComboKey = "g"
		m.pendingComboTime = time.Now()
		m.pendingComboFocus = focusBoard
		return m, comboTickCmd("g"), true
	case focusTree:
		if m.pendingComboKey == "g" && m.pendingComboFocus == focusTree && time.Since(m.pendingComboTime) < comboTimeout {
			m.tree.JumpToTop()
			m.pendingComboKey = ""
			m.pendingComboTime = time.Time{}
			return m, nil, true
		}
		m.pendingComboKey = "g"
		m.pendingComboTime = time.Now()
		m.pendingComboFocus = focusTree
		return m, comboTickCmd("g"), true
	}
	return m, nil, false
}

// dispatchRegistryKeys routes a key through the focus stack registry (bv-p5kf.16).
func (m Model) dispatchRegistryKeys(msg tea.KeyMsg) (Model, bool, tea.Cmd) {
	if m.keyRegistry == nil {
		return m, false, nil
	}

	keyStr := msg.String()

	// Label dashboard uses component Update and may return cmds — keep dedicated path.
	if m.focused == focusLabelDashboard {
		return m.dispatchLabelDashboardKeys(msg)
	}

	// History search/file-tree submodes consume all keys via history handler.
	if m.focused == focusHistory && (m.historyView.IsSearchActive() || m.historyView.FileTreeHasFocus()) {
		return m.handleHistoryKeys(msg), true, nil
	}

	stack := m.DispatchFocusStack()
	updated, handled, _ := m.keyRegistry.DispatchStack(stack, keyStr, m, msg)
	if handled {
		return updated, true, nil
	}

	// Detail viewport scroll for keys not handled by registry.
	if m.focused == focusDetail && !detailPassthroughKey(keyStr) && keyStr != "O" {
		return updated, false, nil
	}

	return updated, false, nil
}

func (m Model) dispatchLabelDashboardKeys(msg tea.KeyMsg) (Model, bool, tea.Cmd) {
	keyStr := msg.String()
	if selectedLabel, cmd := m.labelDashboard.Update(msg); selectedLabel != "" {
		m.currentFilter = "label:" + selectedLabel
		m.applyFilter()
		m.focused = focusList
		return m, true, cmd
	}
	if keyStr == "h" && len(m.labelDashboard.labels) > 0 {
		idx := m.labelDashboard.cursor
		if idx >= 0 && idx < len(m.labelDashboard.labels) {
			lh := m.labelDashboard.labels[idx]
			m.showLabelHealthDetail = true
			m.labelHealthDetail = &lh
			m.labelHealthDetailFlow = m.getCrossFlowsForLabel(lh.Label)
			return m, true, nil
		}
	}
	if keyStr == "d" && len(m.labelDashboard.labels) > 0 {
		idx := m.labelDashboard.cursor
		if idx >= 0 && idx < len(m.labelDashboard.labels) {
			lh := m.labelDashboard.labels[idx]
			m.labelDrilldownLabel = lh.Label
			m.labelDrilldownIssues = m.filterIssuesByLabel(lh.Label)
			m.showLabelDrilldown = true
			return m, true, nil
		}
	}
	return m, true, nil
}

func (m Model) listKeyHandled(key string) bool {
	switch key {
	case "enter", "home", "G", "end", "ctrl+d", "ctrl+u",
		"o", "c", "r", "t", "T", "C", "S", "s", "V", "U", "y":
		return true
	default:
		return false
	}
}

// handleViewToggleKey handles cross-view switch keys previously in Update()'s view-toggle block.
func (m Model) handleViewToggleKey(msg tea.KeyMsg) (Model, bool) {
	switch msg.String() {
	case "b":
		m.clearAttentionOverlay()
		m.isBoardView = !m.isBoardView
		m.isGraphView = false
		m.isActionableView = false
		m.isHistoryView = false
		if m.isBoardView {
			m.focused = focusBoard
			m.refreshBoardAndGraphForCurrentFilter()
		} else {
			m.focused = focusList
		}
		return m, true

	case "g":
		m.clearAttentionOverlay()
		m.isGraphView = !m.isGraphView
		m.isBoardView = false
		m.isActionableView = false
		m.isHistoryView = false
		if m.isGraphView {
			m.focused = focusGraph
			m.refreshBoardAndGraphForCurrentFilter()
		} else {
			m.focused = focusList
		}
		return m, true

	case "a":
		m.clearAttentionOverlay()
		m.isActionableView = !m.isActionableView
		m.isGraphView = false
		m.isBoardView = false
		m.isHistoryView = false
		if m.isActionableView {
			analyzer := analysis.NewAnalyzer(m.issues)
			plan := analyzer.GetExecutionPlan()
			m.actionableView = NewActionableModel(plan, m.theme)
			m.actionableView.SetSize(m.width, m.height-2)
			m.focused = focusActionable
		} else {
			m.focused = focusList
		}
		return m, true

	case "E":
		m.clearAttentionOverlay()
		if m.focused == focusTree {
			m.focused = focusList
		} else {
			m.isGraphView = false
			m.isBoardView = false
			m.isActionableView = false
			m.isHistoryView = false
			if m.snapshot != nil {
				m.tree.BuildFromSnapshot(m.snapshot)
			} else {
				m.tree.Build(m.issues)
			}
			m.tree.SetSize(m.width, m.height-2)
			m.focused = focusTree
		}
		return m, true

	case "i":
		m.clearAttentionOverlay()
		if m.focused == focusInsights {
			m.focused = focusList
		} else {
			m.isGraphView = false
			m.isBoardView = false
			m.isActionableView = false
			m.isHistoryView = false
			m.focused = focusInsights
			m.rebuildInsightsPanel()
		}
		return m, true

	case "p":
		m.showPriorityHints = !m.showPriorityHints
		m.updateListDelegate()
		if m.showPriorityHints {
			count := len(m.priorityHints)
			if count > 0 {
				m.statusMsg = fmt.Sprintf("Priority hints: ↑ increase ↓ decrease (%d suggestions)", count)
			} else {
				m.statusMsg = "Priority hints: No misalignments detected (analysis ongoing)"
			}
		} else {
			m.statusMsg = ""
		}
		return m, true

	case "h":
		m.clearAttentionOverlay()
		m.isHistoryView = !m.isHistoryView
		m.isGraphView = false
		m.isBoardView = false
		m.isActionableView = false
		if m.isHistoryView {
			bodyHeight := m.height - 1
			if bodyHeight < 5 {
				bodyHeight = 5
			}
			m.historyView.SetSize(m.width, bodyHeight)
			m.focused = focusHistory
		} else {
			m.focused = focusList
		}
		return m, true

	case "[", "f3":
		m.clearAttentionOverlay()
		m.isGraphView = false
		m.isBoardView = false
		m.isActionableView = false
		m.isHistoryView = false
		m.focused = focusLabelDashboard
		if !m.labelHealthCached {
			cfg := analysis.DefaultLabelHealthConfig()
			m.labelHealthCache = analysis.ComputeAllLabelHealth(m.issues, cfg, time.Now().UTC(), m.analysis)
			m.labelHealthCached = true
		}
		m.labelDashboard.SetData(m.labelHealthCache.Labels)
		m.labelDashboard.SetSize(m.width, m.height-1)
		m.statusMsg = fmt.Sprintf("Labels: %d total • critical %d • warning %d", m.labelHealthCache.TotalLabels, m.labelHealthCache.CriticalCount, m.labelHealthCache.WarningCount)
		m.statusIsError = false
		return m, true

	case "]", "f4":
		if !m.attentionCached {
			cfg := analysis.DefaultLabelHealthConfig()
			m.attentionCache = analysis.ComputeLabelAttentionScores(m.issues, cfg, time.Now().UTC())
			m.attentionCached = true
		}
		attText, _ := ComputeAttentionView(m.issues, max(40, m.width-4))
		m.isGraphView = false
		m.isBoardView = false
		m.isActionableView = false
		m.isHistoryView = false
		m.focused = focusInsights
		m.showAttentionView = true
		m.rebuildInsightsPanel()
		m.insightsPanel.labelAttention = m.attentionCache.Labels
		m.insightsPanel.extraText = attText
		panelHeight := m.height - 2
		if panelHeight < 3 {
			panelHeight = 3
		}
		m.insightsPanel.SetSize(m.width, panelHeight)
		return m, true

	case "f":
		m.clearAttentionOverlay()
		cfg := analysis.DefaultLabelHealthConfig()
		flow := analysis.ComputeCrossLabelFlow(m.issues, cfg)
		m.isGraphView = false
		m.isBoardView = false
		m.isActionableView = false
		m.isHistoryView = false
		m.focused = focusFlowMatrix
		m.flowMatrix = NewFlowMatrixModel(m.theme)
		m.flowMatrix.SetData(&flow, m.issues)
		panelHeight := m.height - 2
		if panelHeight < 3 {
			panelHeight = 3
		}
		m.flowMatrix.SetSize(m.width, panelHeight)
		return m, true

	case "!":
		activeCount := 0
		for _, a := range m.alerts {
			if !m.dismissedAlerts[alertKey(a)] {
				activeCount++
			}
		}
		if activeCount > 0 {
			m.showAlertsPanel = !m.showAlertsPanel
			m.alertsCursor = 0
		} else {
			m.statusMsg = "No active alerts"
			m.statusIsError = false
		}
		return m, true

	case "'":
		m.showRecipePicker = !m.showRecipePicker
		if m.showRecipePicker {
			m.recipePicker.SetSize(m.width, m.height-1)
			m.focused = focusRecipePicker
		} else {
			m.focused = focusList
		}
		return m, true

	case "w":
		if !m.workspaceMode || len(m.availableRepos) == 0 {
			m.statusMsg = "Repo filter available only in workspace mode"
			m.statusIsError = false
			return m, true
		}
		m.showRepoPicker = !m.showRepoPicker
		if m.showRepoPicker {
			m.repoPicker = NewRepoPickerModel(m.availableRepos, m.theme)
			m.repoPicker.SetActiveRepos(m.activeRepos)
			m.repoPicker.SetSize(m.width, m.height-1)
			m.focused = focusRepoPicker
		} else {
			m.focused = focusList
		}
		return m, true

	case "x":
		m.exportToMarkdown()
		return m, true

	case "l", "L":
		if m, ok := m.openLabelPicker(); ok {
			return m, true
		}
		return m, true

	default:
		return m, false
	}
}
