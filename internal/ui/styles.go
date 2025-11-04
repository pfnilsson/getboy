package ui

import "github.com/charmbracelet/lipgloss"

var (
	accent        = lipgloss.Color("12")
	muted         = lipgloss.Color("240")
	headerStyle   = lipgloss.NewStyle().Bold(true)
	statusStyle   = lipgloss.NewStyle().Faint(true)
	paneBaseStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(accent)
)

func (m model) paneStyle(active bool) lipgloss.Style {
	if active {
		return paneBaseStyle.BorderForeground(accent)
	}
	return paneBaseStyle.BorderForeground(muted)
}
