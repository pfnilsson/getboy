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
		m.editorHeight(),
		m.pane == paneEditor,
		paneBadge(2),
		"Request",
		tabs,
		int(m.activeTab),
	)

	return edBox
}

// viewOverviewTab renders the overview tab (method, URL)
func (m model) viewOverviewTab() string {
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Current.ListSelectedText)
	normalStyle := lipgloss.NewStyle()

	// Determine which field is selected (when in editor pane but not insert mode)
	methodStyle := normalStyle
	urlStyle := normalStyle
	methodPrefix := "  "
	urlPrefix := "  "
	methodSelected := false

	if m.pane == paneEditor && m.activeTab == tabOverview {
		switch m.editorPart {
		case edMethod:
			methodSelected = true
			if !m.insertMode {
				methodPrefix = "> "
				methodStyle = selectedStyle
			}
		case edURL:
			if !m.insertMode {
				urlPrefix = "> "
				urlStyle = selectedStyle
			}
		}
	}

	methodLabel := methodStyle.Render(methodPrefix + "Method: ")
	urlLabel := urlStyle.Render(urlPrefix + "URL:    ")

	// Render method as a dropdown with arrows
	methodValue := m.methodValue()
	var methodDisplay string
	if methodSelected && m.insertMode {
		// Show arrows when editing (vertical arrows for j/k navigation)
		methodDisplay = selectedStyle.Render("▲ " + methodValue + " ▼")
	} else {
		methodDisplay = methodValue
	}

	methodLine := methodLabel + methodDisplay
	urlLine := urlLabel + m.url.View()

	return lipgloss.JoinVertical(lipgloss.Left, methodLine, urlLine)
}

// viewHeadersTab renders the headers tab (placeholder for now)
func (m model) viewHeadersTab() string {
	return lipgloss.NewStyle().
		Faint(true).
		Render("Headers tab - coming soon...")
}

// viewBodyTab renders the body tab with the request body textarea
func (m model) viewBodyTab() string {
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Current.ListSelectedText)
	normalStyle := lipgloss.NewStyle()

	bodyStyle := normalStyle
	bodyPrefix := "  "

	if m.pane == paneEditor && !m.insertMode && m.activeTab == tabBody {
		bodyPrefix = "> "
		bodyStyle = selectedStyle
	}

	bodyTitle := bodyStyle.Render(bodyPrefix + "Body")
	bodyView := m.body.View()

	return lipgloss.JoinVertical(lipgloss.Left, bodyTitle, bodyView)
}
