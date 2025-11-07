package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// viewEditor renders the editor pane containing the method, URL, and body inputs.
func (m model) viewEditor() string {
	// Render content based on active tab
	var content string
	switch m.activeTab {
	case tabOverview:
		content = m.viewOverviewTab()
	case tabHeaders:
		content = m.viewHeadersTab()
	case tabBody:
		content = m.viewBodyTab()
	}

	// Define tabs
	tabs := []string{"Overview", "Headers", "Body"}

	edBox := titledPaneWithTabs(
		content,
		m.rightPaneWidth(),
		m.pane == paneEditor,
		paneBadge(2),
		"Request",
		tabs,
		int(m.activeTab),
	)

	return edBox
}

// viewOverviewTab renders the overview tab (method, URL, body)
func (m model) viewOverviewTab() string {
	methodView := lipgloss.NewStyle().Width(10).Render("Method: " + m.method.View())
	urlView := lipgloss.NewStyle().Width(m.rightPaneWidth() - 12).Render("URL: " + m.url.View())
	edTop := lipgloss.JoinHorizontal(lipgloss.Top, methodView, urlView)

	bodyTitle := titleStyle().Faint(true).Render("Body")
	bodyView := m.body.View()

	return lipgloss.JoinVertical(lipgloss.Left, edTop, bodyTitle, bodyView)
}

// viewHeadersTab renders the headers tab (placeholder for now)
func (m model) viewHeadersTab() string {
	return lipgloss.NewStyle().
		Faint(true).
		Render("Headers tab - coming soon...")
}

// viewBodyTab renders the body tab (placeholder for now)
func (m model) viewBodyTab() string {
	return lipgloss.NewStyle().
		Faint(true).
		Render("Body tab - coming soon...")
}
