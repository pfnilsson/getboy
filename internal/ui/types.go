package ui

type focusPane int

const (
	paneSidebar focusPane = iota
	paneEditor
	paneResponse
)

type editorFocus int

const (
	edMethod editorFocus = iota
	edURL
	edBody
)

type reqItem struct {
	title  string
	desc   string
	method string
	url    string
	body   string
}

func (i reqItem) Title() string {
	return i.title
}

func (i reqItem) Description() string {
	return i.desc
}

func (i reqItem) FilterValue() string {
	return i.title + " " + i.url
}
