package theme

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// palette defines the Catppuccin Mocha color palette
type palette struct {
	text, subtext0, overlay0            string
	blue, sapphire, green, peach, mauve string
}

// Mocha returns the Catppuccin Mocha theme
func Mocha() Theme {
	p := palette{
		text:     "#cdd6f4",
		subtext0: "#a6adc8",
		overlay0: "#6c7086",
		blue:     "#89b4fa",
		sapphire: "#74c7ec",
		green:    "#a6e3a1",
		peach:    "#fab387",
		mauve:    "#cba6f7",
	}

	m := chroma.MustNewStyle("catppuccin-mocha", chroma.StyleEntries{
		chroma.Text:        p.text,
		chroma.Comment:     "italic " + p.overlay0,
		chroma.Punctuation: p.subtext0,
		chroma.NameTag:     p.blue,  // JSON keys
		chroma.String:      p.green, // values
		chroma.Number:      p.peach,
		chroma.Literal:     p.peach,
		chroma.Keyword:     p.mauve,
		chroma.Operator:    p.sapphire,
	})
	styles.Register(m)

	return Theme{
		Name:               "catppuccin-mocha",
		BorderActive:       lipgloss.Color(p.peach),
		BorderInactive:     lipgloss.Color(p.blue),
		Title:              lipgloss.Color(p.blue),
		Header:             lipgloss.Color(p.text),
		Status:             lipgloss.Color(p.subtext0),
		ListSelectedText:   lipgloss.Color(p.mauve),
		ListSelectedBorder: lipgloss.Color(p.mauve),
		ChromaStyle:        "catppuccin-mocha",
	}
}
