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
	var status string
	if m.insertMode {
		status = "-- INSERT --  esc: exit"
	} else {
		switch m.pane {
		case paneSidebar:
			status = "1/2/3: panes  j/k: select  enter: load"
		case paneEditor:
			status = "1/2/3: panes  i: insert  j/k: fields"
			if m.activeTab == tabParams {
				status += "  a: add  d: delete"
			}
			if m.activeTab == tabHeaders {
				status += "  a: add  d: delete  r: toggle view"
			}
		case paneResponse:
			status = "1/2/3: panes  j/k: scroll"
		}
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
