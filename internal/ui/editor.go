package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// viewEditor renders the editor pane containing the method, URL, and body inputs.
func (m model) viewEditor() string {
	methodView := lipgloss.NewStyle().Width(10).Render("Method: " + m.method.View())
	urlView := lipgloss.NewStyle().Width(m.rightPaneWidth() - 12).Render("URL: " + m.url.View())
	edTop := lipgloss.JoinHorizontal(lipgloss.Top, methodView, urlView)

	bodyTitle := titleStyle().Faint(true).Render("Body") // inner section; keep as-is
	ed := lipgloss.JoinVertical(lipgloss.Left, edTop, bodyTitle, m.body.View())

	edBox := titledPane(
		ed,
		m.rightPaneWidth(),
		m.pane == paneEditor,
		paneBadge(2),
		"Request",
	)

	return edBox
}
