package ui

import (
	lg2 "charm.land/lipgloss/v2"
)

// canvasOverlayBox returns tier-2 overlay content rendered as a box (no Place fill)
// for alerts and help — the views migrated to lipgloss v2 Canvas compositing.
func (m Model) canvasOverlayBox(_ int, _ int) (string, bool) {
	switch {
	case m.showAlertsPanel:
		return m.renderAlertsPanelBox(), true
	case m.showHelp:
		mm := m
		return mm.renderHelpOverlayBox(), true
	default:
		return "", false
	}
}

// compositeCenteredOverlay stacks overlay on base using lipgloss v2 Canvas compositing
// (bv-sl44.6). The base view stays visible beneath the centered overlay.
func compositeCenteredOverlay(base, overlay string, width, height int) string {
	if overlay == "" {
		return base
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	ow := lg2.Width(overlay)
	oh := lg2.Height(overlay)
	x := (width - ow) / 2
	y := (height - oh) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	comp := lg2.NewCompositor(
		lg2.NewLayer(base).Z(0),
		lg2.NewLayer(overlay).X(x).Y(y).Z(1),
	)
	canvas := lg2.NewCanvas(width, height)
	return canvas.Compose(comp).Render()
}

// renderViewBody assembles tier-0 base, tier-1 sidebar, and tier-2 overlay/replacement.
func (m Model) renderViewBody() string {
	cw := m.mainContentWidth()
	bodyH := m.height - 1

	base := m.joinShortcutsSidebar(m.renderBaseView(cw, bodyH))
	if box, ok := m.canvasOverlayBox(cw, bodyH); ok {
		return compositeCenteredOverlay(base, box, m.width, bodyH)
	}
	if overlay, ok := m.renderOverlay(cw, bodyH); ok {
		return m.joinShortcutsSidebar(overlay)
	}
	return base
}
