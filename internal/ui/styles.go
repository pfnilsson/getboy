package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pfnilsson/getboy/internal/ui/theme"
)

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

func statusStyle() lipgloss.Style {
	return lipgloss.NewStyle().Faint(true).Foreground(theme.Current.Status)
}

func paneBaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(theme.Current.Title)
}

func activeBorder() lipgloss.Color   { return theme.Current.BorderActive }
func inactiveBorder() lipgloss.Color { return theme.Current.BorderInactive }
func chromaStyle() string            { return theme.Current.ChromaStyle }
func paneBorder() lipgloss.Border    { return lipgloss.RoundedBorder() }

func paneStyle(focused bool) lipgloss.Style {
	col := inactiveBorder()
	if focused {
		col = activeBorder()
	}
	return paneBaseStyle().BorderForeground(col)
}

func topLineStyle(focused bool) lipgloss.Style {
	col := inactiveBorder()
	if focused {
		col = activeBorder()
	}
	return lipgloss.NewStyle().Foreground(col)
}

func paneBadge(n int, focused bool) string {
	return fmt.Sprintf("[%d]", n)
}

func titledPane(content string, width int, focused bool, leftBadge string, leftTitle string) string {
	st := paneStyle(focused).Width(width)

	// Body without top border — we synthesize that ourselves.
	body := st.BorderTop(false).Render(content)
	fullW := lipgloss.Width(body)

	b := paneBorder()
	h, tl, tr := b.Top, b.TopLeft, b.TopRight
	if h == "" {
		h = "─"
	}

	border := topLineStyle(focused)

	// Style labels with the border color, so they "follow" focus.
	if leftBadge != "" {
		leftBadge = border.Bold(true).Render(leftBadge)
	}
	if leftTitle != "" {
		leftTitle = border.Bold(true).Render(leftTitle)
	}

	// Build left side, re-applying border style around every non-label piece.
	left := border.Render(tl + h)
	if leftBadge != "" {
		left += border.Render(" ") + leftBadge + border.Render(" "+h)
	}
	if leftTitle != "" {
		left += border.Render(" ") + leftTitle + border.Render(" ")
	}

	// Build right side with border style.
	right := border.Render(tr)

	// Fill with the border glyph, *also* styled with the border color.
	fill := fullW - lipgloss.Width(left) - lipgloss.Width(right)

	// Graceful degradation for tight widths (keep styling on every branch).
	if fill < 0 && leftTitle != "" {
		left = border.Render(tl+h+" ") + leftBadge + border.Render(" ")
		fill = fullW - lipgloss.Width(left) - lipgloss.Width(right)
	}
	if fill < 0 && leftBadge != "" {
		left = border.Render(tl+h) + leftBadge
		fill = fullW - lipgloss.Width(left) - lipgloss.Width(right)
	}
	if fill < 0 {
		left = border.Render(tl + h)
		fill = fullW - lipgloss.Width(left) - lipgloss.Width(right)
	}
	if fill < 0 {
		fill = 0
	}

	top := left + border.Render(strings.Repeat(h, fill)) + right
	return lipgloss.JoinVertical(lipgloss.Left, top, body)
}
