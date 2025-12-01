package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.recomputeLayout()
		return m, nil

	case tea.KeyMsg:
		// Handle escape to exit insert mode first
		if m.insertMode && msg.String() == "esc" {
			m.insertMode = false
			m.applyFocus()
			return m, nil
		}

		// In insert mode, pass keys directly to the focused input
		if m.insertMode && m.pane == paneEditor {
			var cmd tea.Cmd
			switch m.editorPart {
			case edMethod:
				// Method is a dropdown - handle j/k for cycling
				switch msg.String() {
				case "k", "up":
					m.prevMethod()
				case "j", "down":
					m.nextMethod()
				}
				return m, nil
			case edURL:
				m.url, cmd = m.url.Update(msg)
			case edBody:
				m.body, cmd = m.body.Update(msg)
			}
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.nextPane()
			return m, nil
		case "shift+tab":
			m.prevPane()
			return m, nil
		case "1":
			m.pane = paneSidebar
			m.insertMode = false
			m.applyFocus()
			return m, nil
		case "2":
			m.pane = paneEditor
			m.insertMode = false
			m.applyFocus()
			return m, nil
		case "3":
			m.pane = paneResponse
			m.insertMode = false
			m.applyFocus()
			return m, nil
		case "enter":
			if m.pane == paneSidebar {
				if it, ok := m.sidebar.SelectedItem().(reqItem); ok {
					m.setMethod(it.method)
					m.url.SetValue(it.url)
					m.body.SetValue(it.body)
					m.status = fmt.Sprintf("Loaded '%s'", it.title)
				}
				return m, nil
			}
			method := m.methodValue()
			url := m.ensureURL(m.url.Value())
			if strings.TrimSpace(url) == "" {
				m.status = "Enter a URL first"
				return m, nil
			}
			m.err = nil
			m.loading = true
			m.status = fmt.Sprintf("%s %s…", method, url)
			return m, doHTTP(method, url, m.body.Value())
		}

		var cmd tea.Cmd
		switch m.pane {
		case paneSidebar:
			m.sidebar, cmd = m.sidebar.Update(msg)
			return m, cmd
		case paneEditor:
			switch msg.String() {
			case "i":
				m.insertMode = true
				m.applyFocus()
				return m, nil
			case "up", "k":
				m.prevEditorPart()
				return m, nil
			case "down", "j":
				m.nextEditorPart()
				return m, nil
			case "[", "left":
				m.prevTab()
				return m, nil
			case "]", "right":
				m.nextTab()
				return m, nil
			}
			// Don't pass unhandled keys to text inputs when not in insert mode
			return m, nil
		case paneResponse:
			m.view, cmd = m.view.Update(msg)
			return m, cmd
		}

	case httpDoneMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			m.view.SetContent(fmt.Sprintf("Error: %v", msg.Err))
			m.status = "Request failed"
			return m, nil
		}
		m.view.SetContent(renderResponse(msg.Body))
		m.status = msg.Status
		return m, nil
	}

	return m, nil
}
