package ui

import "github.com/charmbracelet/lipgloss"

func (m model) paneStyle(active bool) lipgloss.Style {
	if active {
		return paneBaseStyle().BorderForeground(activeBorder())
	}
	return paneBaseStyle().BorderForeground(inactiveBorder())
}
