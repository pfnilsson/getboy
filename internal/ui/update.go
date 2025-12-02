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
				// Method is a dropdown - handle j/k for cycling, enter to confirm
				switch msg.String() {
				case "k", "up":
					m.prevMethod()
				case "j", "down":
					m.nextMethod()
				case "enter":
					m.insertMode = false
					m.applyFocus()
				}
				return m, nil
			case edURL:
				// Enter confirms/exits insert mode for URL
				if msg.String() == "enter" {
					m.insertMode = false
					m.applyFocus()
					return m, nil
				}
				m.url, cmd = m.url.Update(msg)
			case edHeaders:
				if m.headersRaw {
					// Raw mode - pass all keys to the textarea
					m.headersRawText, cmd = m.headersRawText.Update(msg)
					return m, cmd
				}
				// Structured mode - handle tab to switch between key and value
				switch msg.String() {
				case "enter":
					// Enter confirms/exits insert mode in structured headers
					m.insertMode = false
					m.applyFocus()
					return m, nil
				case "tab":
					// Move from key to value, or value to next row's key
					if m.headerField == headerKey {
						m.headerField = headerValue
					} else if m.headerIdx < len(m.headers)-1 {
						m.headerField = headerKey
						m.headerIdx++
					}
					m.applyFocus()
					return m, nil
				case "shift+tab":
					// Move from value to key, or key to previous row's value
					if m.headerField == headerValue {
						m.headerField = headerKey
					} else if m.headerIdx > 0 {
						m.headerIdx--
						m.headerField = headerValue
					}
					m.applyFocus()
					return m, nil
				case "up":
					if m.headerIdx > 0 {
						m.headerIdx--
						m.applyFocus()
					}
					return m, nil
				case "down":
					if m.headerIdx < len(m.headers)-1 {
						m.headerIdx++
						m.applyFocus()
					}
					return m, nil
				default:
					// Pass other keys to the focused header input
					if m.headerIdx < len(m.headers) {
						if m.headerField == headerKey {
							m.headers[m.headerIdx].key, cmd = m.headers[m.headerIdx].key.Update(msg)
						} else {
							m.headers[m.headerIdx].value, cmd = m.headers[m.headerIdx].value.Update(msg)
						}
					}
				}
				return m, cmd
			case edBody:
				m.body, cmd = m.body.Update(msg)
			}
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
			return m, doHTTP(method, url, m.body.Value(), m.getHeaders())
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
			case "shift+tab", "left":
				m.prevTab()
				return m, nil
			case "tab", "right":
				m.nextTab()
				return m, nil
			case "h":
				// In headers tab, switch to key field
				if m.activeTab == tabHeaders {
					m.headerField = headerKey
				}
				return m, nil
			case "l":
				// In headers tab, switch to value field
				if m.activeTab == tabHeaders {
					m.headerField = headerValue
				}
				return m, nil
			case "r":
				// Toggle raw mode in headers tab
				if m.activeTab == tabHeaders {
					if m.headersRaw {
						// Switching from raw to structured - parse the raw text
						m.headersFromRaw(m.headersRawText.Value())
					} else {
						// Switching from structured to raw - convert to text
						m.headersRawText.SetValue(m.headersToRaw())
					}
					m.headersRaw = !m.headersRaw
					m.applyFocus()
				}
				return m, nil
			case "a":
				// Add new header row (if in headers tab)
				if m.activeTab == tabHeaders {
					m.addHeaderRow()
				}
				return m, nil
			case "d":
				// Delete current header row (if in headers tab and more than one row)
				if m.activeTab == tabHeaders {
					m.deleteHeaderRow()
				}
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
