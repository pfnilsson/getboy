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
			m.applyFocus()
			return m, nil
		case "2":
			m.pane = paneEditor
			m.applyFocus()
			return m, nil
		case "3":
			m.pane = paneResponse
			m.applyFocus()
			return m, nil
		case "enter":
			if m.pane == paneSidebar {
				if it, ok := m.sidebar.SelectedItem().(reqItem); ok {
					m.method.SetValue(it.method)
					m.url.SetValue(it.url)
					m.body.SetValue(it.body)
					m.status = fmt.Sprintf("Loaded ‘%s’", it.title)
				}
				return m, nil
			}
			method := strings.ToUpper(strings.TrimSpace(m.method.Value()))
			if method == "" {
				method = "GET"
			}
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
			switch m.editorPart {
			case edMethod:
				m.method, cmd = m.method.Update(msg)
			case edURL:
				m.url, cmd = m.url.Update(msg)
			case edBody:
				m.body, cmd = m.body.Update(msg)
			}
			return m, cmd
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
