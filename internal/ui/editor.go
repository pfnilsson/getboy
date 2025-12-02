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

// viewHeadersTab renders the headers tab with key-value pairs
func (m model) viewHeadersTab() string {
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Current.ListSelectedText)
	faintStyle := lipgloss.NewStyle().Faint(true)

	// Mode toggle indicator: [Structured / Raw]
	var modeIndicator string
	if m.headersRaw {
		modeIndicator = "  [" + faintStyle.Render("Structured") + " / " + selectedStyle.Render("Raw") + "]"
	} else {
		modeIndicator = "  [" + selectedStyle.Render("Structured") + " / " + faintStyle.Render("Raw") + "]"
	}

	var lines []string
	lines = append(lines, modeIndicator)

	// Style for input boxes - rounded border
	inputBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Current.BorderInactive).
		Padding(0, 1)

	activeInputBoxStyle := inputBoxStyle.
		BorderForeground(theme.Current.ListSelectedText)

	if m.headersRaw {
		// Raw mode - show the textarea for free-form editing
		lines = append(lines, "  "+m.headersRawText.View())
	} else {
		// Key-value mode with scrolling
		// Each header row is 3 lines tall (border + content + border)
		rowHeight := 3
		// Available height: editor height - pane borders (2) - mode indicator line (1) - scroll indicators (2) - padding (1)
		// Always reserve space for scroll indicators to prevent jumping
		availableHeight := m.editorHeight() - 6
		totalHeaders := len(m.headers)

		visibleRows := max(availableHeight/rowHeight, 1)

		// Calculate scroll window to keep selected row visible
		startIdx := 0
		if totalHeaders > visibleRows {
			// Center the selected row in the visible window
			startIdx = min(max(m.headerIdx-visibleRows/2, 0), totalHeaders-visibleRows)
		}
		endIdx := min(startIdx+visibleRows, totalHeaders)

		// Determine which scroll indicators to show
		hasAbove := startIdx > 0
		hasBelow := endIdx < totalHeaders

		// Always show scroll indicator line (empty or with arrow) to keep layout stable
		if hasAbove {
			lines = append(lines, faintStyle.Render("  ▲ more above"))
		} else {
			lines = append(lines, "") // Empty line to reserve space
		}

		for i := startIdx; i < endIdx; i++ {
			h := m.headers[i]
			isSelected := m.pane == paneEditor && m.activeTab == tabHeaders && m.headerIdx == i

			prefix := "  "
			if isSelected && !m.insertMode {
				prefix = "> "
			}

			// Determine box styles based on selection
			keyBoxStyle := inputBoxStyle
			valBoxStyle := inputBoxStyle

			if isSelected {
				if m.headerField == headerKey {
					keyBoxStyle = activeInputBoxStyle
				} else {
					valBoxStyle = activeInputBoxStyle
				}
			}

			// Render key and value in boxes
			keyView := h.key.View()
			valView := h.value.View()

			// Ensure minimum width for empty inputs
			keyWidth := max(lipgloss.Width(keyView), 15)
			valWidth := max(lipgloss.Width(valView), 20)

			keyBox := keyBoxStyle.Width(keyWidth).Render(keyView)
			valBox := valBoxStyle.Width(valWidth).Render(valView)

			// Separator styled to align vertically with boxes
			separator := lipgloss.NewStyle().
				Height(3).
				AlignVertical(lipgloss.Center).
				Render(" : ")

			// Prefix styled to align vertically
			prefixStyled := lipgloss.NewStyle().
				Height(3).
				AlignVertical(lipgloss.Center).
				Render(prefix)

			// Join horizontally to align the multi-line boxes
			row := lipgloss.JoinHorizontal(lipgloss.Center, prefixStyled, keyBox, separator, valBox)

			lines = append(lines, row)
		}

		// Always show scroll indicator line (empty or with arrow) to keep layout stable
		if hasBelow {
			lines = append(lines, faintStyle.Render("  ▼ more below"))
		} else {
			lines = append(lines, "") // Empty line to reserve space
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewBodyTab renders the body tab with the request body textarea
func (m model) viewBodyTab() string {
	selectedStyle := lipgloss.NewStyle().Foreground(theme.Current.ListSelectedText)
	normalStyle := lipgloss.NewStyle()

	bodyStyle := normalStyle
	bodyPrefix := "  "

	isEditing := m.pane == paneEditor && m.insertMode && m.activeTab == tabBody

	if m.pane == paneEditor && !m.insertMode && m.activeTab == tabBody {
		bodyPrefix = "> "
		bodyStyle = selectedStyle
	}

	bodyTitle := bodyStyle.Render(bodyPrefix + "Body")

	var bodyView string
	if isEditing {
		// Show plain textarea when editing
		bodyView = m.body.View()
	} else {
		// Show syntax-highlighted content when not editing
		content := m.body.Value()
		if content == "" {
			bodyView = m.body.View() // Show placeholder
		} else {
			bodyView = m.highlightBodyContent(content)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, bodyTitle, bodyView)
}
