package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// Render each pane using dedicated view functions
	sbBox := m.viewSidebar()
	edBox := m.viewEditor()
	respBox := m.viewResponse()

	right := lipgloss.JoinVertical(lipgloss.Left, edBox, respBox)

	// ===== Footer / Status ====================================================
	status := m.status
	if m.insertMode {
		status = "-- INSERT --  esc: exit insert mode"
	}
	if m.loading {
		status += "  ·  loading…"
	}
	if m.err != nil {
		status += "  ·  error: " + m.err.Error()
	}
	footer := statusStyle().Render(status)

	// ===== Layout =============================================================
	content := lipgloss.JoinHorizontal(lipgloss.Top, sbBox, right)
	return lipgloss.JoinVertical(lipgloss.Left, content, footer)
}
