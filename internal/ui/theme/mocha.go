package theme

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// Mocha returns the Catppuccin Mocha theme
func Mocha() Theme {
	p := Palette{
		Text: "#cdd6f4", Subtext0: "#a6adc8", Overlay0: "#6c7086",
		Surface0: "#313244", Surface1: "#45475a",
		Base: "#1e1e2e", Mantle: "#181825", Crust: "#11111b",
		Blue: "#89b4fa", Lavender: "#b4befe", Sapphire: "#74c7ec", Sky: "#89dceb",
		Teal: "#94e2d5", Green: "#a6e3a1", Yellow: "#f9e2af", Peach: "#fab387",
		Maroon: "#eba0ac", Red: "#f38ba8", Mauve: "#cba6f7", Flamingo: "#f2cdcd",
		Pink: "#f5c2e7", Rosewater: "#f5e0dc",
	}

	m := chroma.MustNewStyle("catppuccin-mocha", chroma.StyleEntries{
		chroma.Text:        p.Text,
		chroma.Comment:     "italic " + p.Overlay0,
		chroma.Punctuation: p.Subtext0,
		chroma.NameTag:     p.Blue,  // JSON keys
		chroma.String:      p.Green, // values
		chroma.Number:      p.Peach,
		chroma.Literal:     p.Peach,
		chroma.Keyword:     p.Mauve,
		chroma.Operator:    p.Sapphire,
	})
	styles.Register(m)

	return Theme{
		Name:               "catppuccin-mocha",
		BorderActive:       lipgloss.Color(p.Peach),
		BorderInactive:     lipgloss.Color(p.Blue),
		Title:              lipgloss.Color(p.Blue),
		Header:             lipgloss.Color(p.Text),
		Status:             lipgloss.Color(p.Subtext0),
		ListSelectedText:   lipgloss.Color(p.Mauve),
		ListSelectedBorder: lipgloss.Color(p.Mauve),
		ChromaStyle:        "catppuccin-mocha",
	}
}
