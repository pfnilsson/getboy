package ui

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	header := headerStyle().Render("getboy")

	sbTitle := titleStyle().Render("Requests")
	sb := lipgloss.JoinVertical(lipgloss.Left, sbTitle, m.sidebar.View())
	sbBox := m.paneStyle(m.pane == paneSidebar).Width(m.sidebarWidth()).Render(sb)

	methodView := lipgloss.NewStyle().Width(10).Render("Method: " + m.method.View())
	urlView := lipgloss.NewStyle().Width(m.rightPaneWidth() - 12).Render("URL: " + m.url.View())
	edTop := lipgloss.JoinHorizontal(lipgloss.Top, methodView, urlView)
	bodyTitle := titleStyle().Faint(true).Render("Body")
	ed := lipgloss.JoinVertical(lipgloss.Left, edTop, bodyTitle, m.body.View())
	edBox := m.paneStyle(m.pane == paneEditor).Width(m.rightPaneWidth()).Render(ed)

	respTitle := titleStyle().Render("Response")
	resp := lipgloss.JoinVertical(lipgloss.Left, respTitle, m.view.View())
	respBox := m.paneStyle(m.pane == paneResponse).Width(m.rightPaneWidth()).Render(resp)

	right := lipgloss.JoinVertical(lipgloss.Left, edBox, respBox)

	status := m.status
	if m.loading {
		status += "  ·  loading…"
	}
	if m.err != nil {
		status += "  ·  error: " + m.err.Error()
	}
	footer := statusStyle().Render(status)

	content := lipgloss.JoinHorizontal(lipgloss.Top, sbBox, right)
	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}
