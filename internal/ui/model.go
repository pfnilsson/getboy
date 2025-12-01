package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pfnilsson/getboy/internal/ui/theme"
)

// HTTP methods available in the dropdown
var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

type model struct {
	width  int
	height int

	sidebar list.Model

	methodIdx int // index into httpMethods
	url       textinput.Model
	body      textarea.Model
	view      viewport.Model

	pane       focusPane
	editorPart editorFocus
	activeTab  requestTab
	insertMode bool

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

	// Ensure all inputs start blurred (not in insert mode)
	u.Blur()
	t.Blur()

	vp := viewport.New(0, 0)
	vp.SetContent("Response will appear here…")

	return model{
		sidebar:   sb,
		methodIdx: 0, // Default to GET
		url:       u,
		body:      t,
		view:      vp,
		pane:      paneSidebar,
		activeTab: tabOverview,
		status:    "1/2/3: panes  j/k: select  enter: load",
	}
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
		// Headers tab has no editable fields yet
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
		// Headers tab has no editable fields yet
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
		// No editable fields yet
	case tabBody:
		m.editorPart = edBody
	}
}

func (m *model) applyFocus() {
	m.url.Blur()
	m.body.Blur()

	// Only focus text inputs when in insert mode
	// Method is a dropdown, not a text input, so it doesn't need focus
	if m.pane == paneEditor && m.insertMode {
		switch m.editorPart {
		case edMethod:
			// Method is a dropdown - no focus needed
		case edURL:
			m.url.Focus()
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
