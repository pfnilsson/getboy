package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pfnilsson/getboy/internal/ui/theme"
)

// HTTP methods available in the dropdown
var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// headerRow represents a single header key-value pair
type headerRow struct {
	key   textinput.Model
	value textinput.Model
}

// newHeaderRow creates a new empty header row
func newHeaderRow() headerRow {
	k := textinput.New()
	k.Placeholder = "Header-Name"
	k.CharLimit = 256
	k.Prompt = ""

	v := textinput.New()
	v.Placeholder = "value"
	v.CharLimit = 4096
	v.Prompt = ""

	return headerRow{key: k, value: v}
}

type model struct {
	width  int
	height int

	sidebar list.Model

	methodIdx      int // index into httpMethods
	url            textinput.Model
	headers        []headerRow
	headersRawText textarea.Model // textarea for raw headers mode
	body           textarea.Model
	view           viewport.Model

	pane        focusPane
	editorPart  editorFocus
	activeTab   requestTab
	insertMode  bool
	headerIdx   int         // which header row is selected
	headerField headerField // key or value within the row
	headersRaw  bool        // toggle for raw view mode

	status  string
	loading bool
	err     error
}

// methodValue returns the currently selected HTTP method
func (m model) methodValue() string {
	if m.methodIdx >= 0 && m.methodIdx < len(httpMethods) {
		return httpMethods[m.methodIdx]
	}
	return "GET"
}

// setMethod sets the method by name, defaulting to GET if not found
func (m *model) setMethod(method string) {
	method = strings.ToUpper(strings.TrimSpace(method))
	for i, meth := range httpMethods {
		if meth == method {
			m.methodIdx = i
			return
		}
	}
	m.methodIdx = 0 // Default to GET
}

// nextMethod cycles to the next HTTP method
func (m *model) nextMethod() {
	m.methodIdx = (m.methodIdx + 1) % len(httpMethods)
}

// prevMethod cycles to the previous HTTP method
func (m *model) prevMethod() {
	m.methodIdx = (m.methodIdx + len(httpMethods) - 1) % len(httpMethods)
}

func New() tea.Model {
	items := []list.Item{
		reqItem{title: "Example GET", desc: "GET https://httpbin.org/get", method: "GET", url: "https://httpbin.org/get"},
		reqItem{title: "Echo POST", desc: "POST https://httpbin.org/post", method: "POST", url: "https://httpbin.org/post", body: `{"hello":"world"}`},
		reqItem{title: "Example JSON", desc: "GET https://jsonplaceholder.typicode.com/todos/1", method: "GET", url: "https://jsonplaceholder.typicode.com/todos/1"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(theme.Current.ListSelectedText).
		BorderForeground(theme.Current.ListSelectedBorder)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(theme.Current.ListSelectedText).
		BorderForeground(theme.Current.ListSelectedBorder)

	sb := list.New(items, delegate, 24, 20)
	sb.SetShowTitle(false)
	sb.SetShowHelp(false)
	sb.SetFilteringEnabled(true)
	sb.SetShowStatusBar(false)

	u := textinput.New()
	u.Placeholder = "https://example.com"
	u.CharLimit = 2048
	u.Prompt = ""

	t := textarea.New()
	t.SetWidth(40)
	t.SetHeight(6)
	t.Placeholder = "Request body (optional)"
	t.ShowLineNumbers = false
	t.Prompt = ""
	t.FocusedStyle.CursorLine = lipgloss.NewStyle()
	t.BlurredStyle.CursorLine = lipgloss.NewStyle()

	// Raw headers textarea
	rawHeaders := textarea.New()
	rawHeaders.SetWidth(40)
	rawHeaders.SetHeight(6)
	rawHeaders.Placeholder = ""
	rawHeaders.ShowLineNumbers = false
	rawHeaders.Prompt = ""
	rawHeaders.FocusedStyle.CursorLine = lipgloss.NewStyle()
	rawHeaders.BlurredStyle.CursorLine = lipgloss.NewStyle()

	// Ensure all inputs start blurred (not in insert mode)
	u.Blur()
	t.Blur()
	rawHeaders.Blur()

	vp := viewport.New(0, 0)
	vp.SetContent("Response will appear here…")

	// Start with one empty header row
	headers := []headerRow{newHeaderRow()}

	return model{
		sidebar:        sb,
		methodIdx:      0, // Default to GET
		url:            u,
		headers:        headers,
		headersRawText: rawHeaders,
		body:           t,
		view:           vp,
		pane:           paneSidebar,
		activeTab:      tabOverview,
		status:         "1/2/3: panes  j/k: select  enter: load",
	}
}

// addHeaderRow adds a new header row after the current one
func (m *model) addHeaderRow() {
	newRow := newHeaderRow()
	// Insert after current row
	idx := m.headerIdx + 1
	m.headers = append(m.headers[:idx], append([]headerRow{newRow}, m.headers[idx:]...)...)
	m.headerIdx = idx
	m.headerField = headerKey
}

// deleteHeaderRow removes the current header row if there's more than one
func (m *model) deleteHeaderRow() {
	if len(m.headers) <= 1 {
		// Can't delete the last row, just clear it
		m.headers[0].key.SetValue("")
		m.headers[0].value.SetValue("")
		return
	}
	m.headers = append(m.headers[:m.headerIdx], m.headers[m.headerIdx+1:]...)
	if m.headerIdx >= len(m.headers) {
		m.headerIdx = len(m.headers) - 1
	}
}

// headersToRaw converts headers to raw text format
func (m model) headersToRaw() string {
	var lines []string
	for _, h := range m.headers {
		k, v := h.key.Value(), h.value.Value()
		if k != "" || v != "" {
			lines = append(lines, k+": "+v)
		}
	}
	return strings.Join(lines, "\n")
}

// headersFromRaw parses raw text into header rows
func (m *model) headersFromRaw(raw string) {
	lines := strings.Split(raw, "\n")
	m.headers = nil
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row := newHeaderRow()
		if idx := strings.Index(line, ":"); idx >= 0 {
			row.key.SetValue(strings.TrimSpace(line[:idx]))
			row.value.SetValue(strings.TrimSpace(line[idx+1:]))
		} else {
			row.key.SetValue(line)
		}
		m.headers = append(m.headers, row)
	}
	if len(m.headers) == 0 {
		m.headers = []headerRow{newHeaderRow()}
	}
	m.headerIdx = 0
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) nextPane() {
	m.pane = (m.pane + 1) % 3
	m.insertMode = false
	m.applyFocus()
}

func (m *model) prevPane() {
	m.pane = (m.pane + 2) % 3
	m.insertMode = false
	m.applyFocus()
}

func (m *model) nextEditorPart() {
	switch m.activeTab {
	case tabOverview:
		// Overview has Method and URL (cycle between 0 and 1)
		if m.editorPart == edMethod {
			m.editorPart = edURL
		} else {
			m.editorPart = edMethod
		}
	case tabHeaders:
		m.editorPart = edHeaders
		// Move to next header row
		if m.headerIdx < len(m.headers)-1 {
			m.headerIdx++
		}
	case tabBody:
		// Body tab has only one field, no navigation needed
		m.editorPart = edBody
	}
	m.applyFocus()
}

func (m *model) prevEditorPart() {
	switch m.activeTab {
	case tabOverview:
		// Overview has Method and URL (cycle between 0 and 1)
		if m.editorPart == edURL {
			m.editorPart = edMethod
		} else {
			m.editorPart = edURL
		}
	case tabHeaders:
		m.editorPart = edHeaders
		// Move to previous header row
		if m.headerIdx > 0 {
			m.headerIdx--
		}
	case tabBody:
		// Body tab has only one field, no navigation needed
		m.editorPart = edBody
	}
	m.applyFocus()
}

func (m *model) nextTab() {
	m.activeTab = (m.activeTab + 1) % 3
	m.resetEditorPartForTab()
}

func (m *model) prevTab() {
	m.activeTab = (m.activeTab + 2) % 3
	m.resetEditorPartForTab()
}

// resetEditorPartForTab sets the editorPart to the first field of the current tab
func (m *model) resetEditorPartForTab() {
	switch m.activeTab {
	case tabOverview:
		m.editorPart = edMethod
	case tabHeaders:
		m.editorPart = edHeaders
		m.headerIdx = 0
		m.headerField = headerKey
	case tabBody:
		m.editorPart = edBody
	}
}

func (m *model) applyFocus() {
	m.url.Blur()
	m.body.Blur()
	m.headersRawText.Blur()
	// Blur all header inputs
	for i := range m.headers {
		m.headers[i].key.Blur()
		m.headers[i].value.Blur()
	}

	// Only focus text inputs when in insert mode
	// Method is a dropdown, not a text input, so it doesn't need focus
	if m.pane == paneEditor && m.insertMode {
		switch m.editorPart {
		case edMethod:
			// Method is a dropdown - no focus needed
		case edURL:
			m.url.Focus()
		case edHeaders:
			if m.headersRaw {
				// In raw mode, focus the textarea
				m.headersRawText.Focus()
			} else if m.headerIdx >= 0 && m.headerIdx < len(m.headers) {
				if m.headerField == headerKey {
					m.headers[m.headerIdx].key.Focus()
				} else {
					m.headers[m.headerIdx].value.Focus()
				}
			}
		case edBody:
			m.body.Focus()
		}
	}
}

func (m model) ensureURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return u
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return "https://" + u
	}
	return u
}

// getHeaders returns headers as a map, skipping empty keys
func (m model) getHeaders() map[string]string {
	result := make(map[string]string)
	for _, h := range m.headers {
		k := strings.TrimSpace(h.key.Value())
		if k != "" {
			result[k] = h.value.Value()
		}
	}
	return result
}

// getContentType returns the Content-Type header value (case-insensitive lookup)
func (m model) getContentType() string {
	for _, h := range m.headers {
		if strings.EqualFold(strings.TrimSpace(h.key.Value()), "content-type") {
			return strings.ToLower(strings.TrimSpace(h.value.Value()))
		}
	}
	return ""
}

// highlightBodyContent applies syntax highlighting based on Content-Type header
func (m model) highlightBodyContent(content string) string {
	ct := m.getContentType()

	switch {
	case strings.Contains(ct, "json"):
		return highlight(content, "json")
	case strings.Contains(ct, "xml"):
		return highlight(content, "xml")
	default:
		return content
	}
}
