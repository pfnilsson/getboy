package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// App header
	header := headerStyle().Render("getboy")

	// ===== Sidebar ============================================================
	sb := m.sidebar.View()
	sbBox := TitledPane(
		sb,
		m.sidebarWidth(),
		m.pane == paneSidebar,
		PaneBadge(1, m.pane == paneSidebar),
		"Requests",
	)

	// ===== Editor (Request) ===================================================
	methodView := lipgloss.NewStyle().Width(10).Render("Method: " + m.method.View())
	urlView := lipgloss.NewStyle().Width(m.rightPaneWidth() - 12).Render("URL: " + m.url.View())
	edTop := lipgloss.JoinHorizontal(lipgloss.Top, methodView, urlView)

	bodyTitle := titleStyle().Faint(true).Render("Body") // inner section; keep as-is
	ed := lipgloss.JoinVertical(lipgloss.Left, edTop, bodyTitle, m.body.View())

	edBox := TitledPane(
		ed,
		m.rightPaneWidth(),
		m.pane == paneEditor,
		PaneBadge(2, m.pane == paneEditor),
		"Request",
	)

	// ===== Response ===========================================================
	resp := lipgloss.JoinVertical(lipgloss.Left, m.view.View())
	respBox := TitledPane(
		resp,
		m.rightPaneWidth(),
		m.pane == paneResponse,
		PaneBadge(3, m.pane == paneResponse),
		"Response",
	)

	right := lipgloss.JoinVertical(lipgloss.Left, edBox, respBox)

	// ===== Footer / Status ====================================================
	status := m.status
	if m.loading {
		status += "  ·  loading…"
	}
	if m.err != nil {
		status += "  ·  error: " + m.err.Error()
	}
	footer := statusStyle().Render(status)

	// ===== Layout =============================================================
	content := lipgloss.JoinHorizontal(lipgloss.Top, sbBox, right)
	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}
