package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	width  int
	height int

	sidebar list.Model

	method textinput.Model
	url    textinput.Model
	body   textarea.Model
	view   viewport.Model

	pane       focusPane
	editorPart editorFocus

	status  string
	loading bool
	err     error
}

func New() tea.Model {
	items := []list.Item{
		reqItem{title: "Example GET", desc: "GET https://httpbin.org/get", method: "GET", url: "https://httpbin.org/get"},
		reqItem{title: "Echo POST", desc: "POST https://httpbin.org/post", method: "POST", url: "https://httpbin.org/post", body: `{"hello":"world"}`},
		reqItem{title: "Example JSON", desc: "GET https://jsonplaceholder.typicode.com/todos/1", method: "GET", url: "https://jsonplaceholder.typicode.com/todos/1"},
	}

	delegate := list.NewDefaultDelegate()
	sb := list.New(items, delegate, 24, 20)
	sb.Title = "Requests"
	sb.SetShowHelp(false)
	sb.SetFilteringEnabled(true)
	sb.SetShowStatusBar(false)

	mth := textinput.New()
	mth.Placeholder = "GET"
	mth.CharLimit = 10
	mth.Width = 8
	mth.Prompt = ""

	u := textinput.New()
	u.Placeholder = "https://example.com"
	u.CharLimit = 2048
	u.Prompt = ""

	t := textarea.New()
	t.SetWidth(40)
	t.SetHeight(6)
	t.Placeholder = "Request body (optional)"
	t.ShowLineNumbers = false

	vp := viewport.New(0, 0)
	vp.SetContent("Response will appear here…")

	return model{
		sidebar: sb,
		method:  mth,
		url:     u,
		body:    t,
		view:    vp,
		pane:    paneSidebar,
		status:  "tab: switch panes  •  enter: run  •  j/k: move  •  q: quit",
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) nextPane() {
	m.pane = (m.pane + 1) % 3
	m.applyFocus()
}

func (m *model) prevPane() {
	m.pane = (m.pane + 2) % 3
	m.applyFocus()
}

func (m *model) nextEditorPart() {
	m.editorPart = (m.editorPart + 1) % 3
	m.applyFocus()
}

func (m *model) prevEditorPart() {
	m.editorPart = (m.editorPart + 2) % 3
	m.applyFocus()
}

func (m *model) applyFocus() {
	m.method.Blur()
	m.url.Blur()
	m.body.Blur()

	switch m.pane {
	case paneSidebar:
	case paneEditor:
		switch m.editorPart {
		case edMethod:
			m.method.Focus()
		case edURL:
			m.url.Focus()
		case edBody:
			m.body.Focus()
		}
	case paneResponse:
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
