package theme

import "github.com/charmbracelet/lipgloss"

// Palette defines the color palette for a theme
type Palette struct {
	Text, Subtext0, Overlay0, Surface0, Surface1, Base, Mantle, Crust                                        string
	Blue, Lavender, Sapphire, Sky, Teal, Green, Yellow, Peach, Maroon, Red, Mauve, Flamingo, Pink, Rosewater string
}

// Theme defines a complete UI theme
type Theme struct {
	Name               string
	BorderActive       lipgloss.Color
	BorderInactive     lipgloss.Color
	Title              lipgloss.Color
	Header             lipgloss.Color
	Status             lipgloss.Color
	ListSelectedText   lipgloss.Color
	ListSelectedBorder lipgloss.Color
	ChromaStyle        string // name registered with chroma
}

// Current is the active theme used by the UI
var Current = Mocha()
