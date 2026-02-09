package pkg

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/ltfred/alacritty/pkg/themes"
)

type OverlayModel struct {
	background   themeModel
	foreground   confirmModel
	theme        themes.Theme
	leftWidth    int
	isConfirming bool
}

func (m *OverlayModel) Init() tea.Cmd {
	return nil
}

func (m *OverlayModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.background.maxHeight = msg.Height
		m.background.maxWidth = msg.Width - m.leftWidth

	case tea.KeyMsg:
	}
	return m, nil
}

func (m *OverlayModel) View() string {
	if m.foreground.title == "" && m.background.View() == "" {
		return ""
	}
	if m.foreground.View() == "" && m.background.View() != "" {
		return m.background.View()
	}
	if m.foreground.View() != "" && m.background.View() == "" {
		return m.background.View()
	}
	return Composite(
		m.foreground.View(),
		m.background.View(),
		Center,
		Center,
		0,
		0,
	)
}

func Composite(fg, bg string, xPos, yPos Position, xOff, yOff int) string {
	if fg == "" {
		return bg
	}
	if bg == "" {
		return fg
	}
	if strings.Count(fg, "\n") == 0 && strings.Count(bg, "\n") == 0 {
		// Simple case: both are single line
		return fg
	}

	fgWidth, fgHeight := lipgloss.Size(fg)
	bgWidth, bgHeight := lipgloss.Size(bg)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return fg
	}

	x, y := offsets(fg, bg, xPos, yPos, xOff, yOff)
	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	fgLines := lines(fg)
	bgLines := lines(bg)
	var sb strings.Builder

	for i, bgLine := range bgLines {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			sb.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := ansi.Truncate(bgLine, x, "")
			pos = ansi.StringWidth(left)
			sb.WriteString(left)
			if pos < x {
				sb.WriteString(whitespace(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		sb.WriteString(fgLine)
		pos += ansi.StringWidth(fgLine)

		right := ansi.TruncateLeft(bgLine, pos, "")
		bgWidth := ansi.StringWidth(bgLine)
		rightWidth := ansi.StringWidth(right)
		if rightWidth <= bgWidth-pos {
			sb.WriteString(whitespace(bgWidth - rightWidth - pos))
		}
		sb.WriteString(right)
	}
	return sb.String()
}

func offsets(fg, bg string, xPos, yPos Position, xOff, yOff int) (int, int) {
	var x, y int

	// Handle X axis positioning
	switch xPos {
	case Left:
		x = 0
	case Center:
		halfBackgroundWidth := lipgloss.Width(bg) / 2
		halfForegroundWidth := lipgloss.Width(fg) / 2
		x = halfBackgroundWidth - halfForegroundWidth
	case Right:
		x = lipgloss.Width(bg) - lipgloss.Width(fg)
	}

	// Handle Y axis positioning
	switch yPos {
	case Top:
		y = 0
	case Center:
		halfBackgroundHeight := lipgloss.Height(bg) / 2
		halfForegroundHeight := lipgloss.Height(fg) / 2
		y = halfBackgroundHeight - halfForegroundHeight
	case Bottom:
		y = lipgloss.Height(bg) - lipgloss.Height(fg)
	}

	return x + xOff, y + yOff
}

// clamp clamps a value between lower and upper bounds.
func clamp(v, lower, upper int) int {
	// Handle case where bounds are inverted (e.g., min > max)
	if lower > upper {
		lower, upper = upper, lower
	}
	if v < lower {
		return lower
	}
	if v > upper {
		return upper
	}
	return v
}

// lines normalises any non standard new lines within a string, and then splits and returns a slice
// of strings split on the new lines.
func lines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}

// whitescpace returns a string of whitespace characters of the requested length.
func whitespace(length int) string {
	return strings.Repeat(" ", length)
}

type Position int

const (
	Top Position = iota + 1
	Right
	Bottom
	Left
	Center
)
