package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines a complete UI theme
type Theme struct {
	Name               string
	BorderActive       lipgloss.Color
	BorderInactive     lipgloss.Color
	Title              lipgloss.Color
	Header             lipgloss.Color
	Status             lipgloss.Color
	TabActive          lipgloss.Color
	ListSelectedText   lipgloss.Color
	ListSelectedBorder lipgloss.Color
	ChromaStyle        string // name registered with chroma
}

// Current is the active theme used by the UI
var Current = Mocha()
