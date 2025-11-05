package ui

import (
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

type CatppuccinPalette struct { // Mocha
	Text, Subtext0, Overlay0, Surface0, Surface1, Base, Mantle, Crust                                        string
	Blue, Lavender, Sapphire, Sky, Teal, Green, Yellow, Peach, Maroon, Red, Mauve, Flamingo, Pink, Rosewater string
}

type Theme struct {
	Name           string
	Palette        CatppuccinPalette
	BorderActive   lipgloss.Color
	BorderInactive lipgloss.Color
	Title          lipgloss.Color
	Header         lipgloss.Color
	Status         lipgloss.Color
	ChromaStyle    string // name registered with chroma
}

// Build the Catppuccin Mocha theme and register its chroma style.
func CatppuccinMocha() Theme {
	p := CatppuccinPalette{
		Text: "#cdd6f4", Subtext0: "#a6adc8", Overlay0: "#6c7086",
		Surface0: "#313244", Surface1: "#45475a",
		Base: "#1e1e2e", Mantle: "#181825", Crust: "#11111b",
		Blue: "#89b4fa", Lavender: "#b4befe", Sapphire: "#74c7ec", Sky: "#89dceb",
		Teal: "#94e2d5", Green: "#a6e3a1", Yellow: "#f9e2af", Peach: "#fab387",
		Maroon: "#eba0ac", Red: "#f38ba8", Mauve: "#cba6f7", Flamingo: "#f2cdcd",
		Pink: "#f5c2e7", Rosewater: "#f5e0dc",
	}

	// Register a chroma style with Mocha colors (JSON-focused; expand as needed).
	m := chroma.MustNewStyle("catppuccin-mocha", chroma.StyleEntries{
		chroma.Text:        p.Text,
		chroma.Comment:     "italic " + p.Overlay0,
		chroma.Punctuation: p.Subtext0,
		chroma.NameTag:     p.Blue,  // JSON keys
		chroma.String:      p.Green, // values
		chroma.Number:      p.Peach,
		chroma.Literal:     p.Peach,
		chroma.Keyword:     p.Mauve, // for other languages later
		chroma.Operator:    p.Sapphire,
	})
	styles.Register(m)

	return Theme{
		Name:           "catppuccin-mocha",
		Palette:        p,
		BorderActive:   lipgloss.Color(p.Mauve),
		BorderInactive: lipgloss.Color(p.Overlay0),
		Title:          lipgloss.Color(p.Blue),
		Header:         lipgloss.Color(p.Text),
		Status:         lipgloss.Color(p.Subtext0),
		ChromaStyle:    "catppuccin-mocha",
	}
}

var currentTheme = CatppuccinMocha()

// Pick best ANSI formatter based on terminal capability.
func bestFormatter() string {
	ct := strings.ToLower(os.Getenv("COLORTERM"))
	if strings.Contains(ct, "truecolor") || strings.Contains(ct, "24bit") {
		return "terminal16m"
	}
	if strings.Contains(strings.ToLower(os.Getenv("TERM")), "256color") {
		return "terminal256"
	}
	return "terminal"
}

// Lip Gloss style helpers (read from currentTheme)
func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(currentTheme.Header)
}

func statusStyle() lipgloss.Style {
	return lipgloss.NewStyle().Faint(true).Foreground(currentTheme.Status)
}

func paneBaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(currentTheme.Title)
}
func activeBorder() lipgloss.Color   { return currentTheme.BorderActive }
func inactiveBorder() lipgloss.Color { return currentTheme.BorderInactive }
func chromaStyle() string            { return currentTheme.ChromaStyle }
