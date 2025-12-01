package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pfnilsson/getboy/internal/ui/theme"
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
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Current.ListSelectedText)
	normalStyle := lipgloss.NewStyle()

	// Determine which field is selected (when in editor pane but not insert mode)
	methodStyle := normalStyle
	urlStyle := normalStyle
	bodyStyle := normalStyle
	methodPrefix := "  "
	urlPrefix := "  "
	bodyPrefix := "  "

	if m.pane == paneEditor && !m.insertMode {
		switch m.editorPart {
		case edMethod:
			methodPrefix = "> "
			methodStyle = selectedStyle
		case edURL:
			urlPrefix = "> "
			urlStyle = selectedStyle
		case edBody:
			bodyPrefix = "> "
			bodyStyle = selectedStyle
		}
	}

	methodLabel := methodStyle.Render(methodPrefix + "Method: ")
	urlLabel := urlStyle.Render(urlPrefix + "URL: ")

	methodView := lipgloss.NewStyle().Width(12).Render(methodLabel + m.method.View())
	urlView := lipgloss.NewStyle().Width(m.rightPaneWidth() - 14).Render(urlLabel + m.url.View())
	edTop := lipgloss.JoinHorizontal(lipgloss.Top, methodView, urlView)

	bodyTitle := bodyStyle.Render(bodyPrefix + "Body")
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
