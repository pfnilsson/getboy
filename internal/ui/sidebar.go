package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// viewSidebar renders the sidebar pane containing the list of saved requests.
func (m model) viewSidebar() string {
	var content string

	switch m.sidebarTab {
	case sidebarHistory:
		if len(m.history) == 0 {
			emptyStyle := lipgloss.NewStyle().
				Faint(true).
				Padding(1, 2)
			content = emptyStyle.Render("No history yet.\nSend a request to get started.")
		} else {
			content = m.sidebar.View()
		}
	case sidebarSaved:
		emptyStyle := lipgloss.NewStyle().
			Faint(true).
			Padding(1, 2)
		content = emptyStyle.Render("Saved requests\ncoming soon...")
	}

	tabs := []string{"History", "Saved"}

	sbBox := titledPaneWithTabs(
		content,
		m.sidebarWidth(),
		m.contentHeight(),
		m.pane == paneSidebar,
		paneBadge(1),
		"",
		tabs,
		int(m.sidebarTab),
	)
	return sbBox
}
