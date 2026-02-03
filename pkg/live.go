package pkg

import (
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type liveMode struct {
	theme     theme
	maxWidth  int
	maxHeight int
}

func (m liveMode) Init() tea.Cmd {
	return nil
}

func (m liveMode) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m liveMode) View() string {
	theme := decodeTheme(m.theme.data)
	backgroundStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Colors.Primary.Background))

	themeNameStyle := backgroundStyle.Bold(true).Padding(1, 1).Italic(true).
		Width(m.maxWidth).Foreground(lipgloss.Color(theme.Colors.Primary.Foreground)).Align(lipgloss.Center)

	colorBlockStyle := backgroundStyle.Padding(1, 1).Width(m.maxWidth)

	// 创建彩色方块
	colorNames := []string{"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White"}
	colors := make([]colorBlock, 0, len(colorNames))
	for _, name := range colorNames {
		normalColor := reflect.ValueOf(theme.Colors.Normal).FieldByName(name).String()
		brightColor := reflect.ValueOf(theme.Colors.Bright).FieldByName(name).String()
		colors = append(colors, colorBlock{
			normalColor: normalColor,
			brightColor: brightColor,
			name:        name,
		})
	}
	colorBlock := renderColorBlock(theme.Colors.Primary.Background, theme.Colors.Primary.Foreground, colors)

	alacritty := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Cyan,
		Dark:  theme.Colors.Normal.Cyan,
	}).Render("alacritty")

	branch := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Blue,
		Dark:  theme.Colors.Normal.Blue,
	}).Render("\uF418 main")

	stash := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Red,
		Dark:  theme.Colors.Normal.Red,
	}).Render("[+]")

	program := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Yellow,
		Dark:  theme.Colors.Normal.Yellow,
	}).Render("\uE626 v1.26")

	on := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Yellow,
		Dark:  theme.Colors.Normal.Yellow,
	}).Render("on")
	via := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Yellow,
		Dark:  theme.Colors.Normal.Yellow,
	}).Render("via")

	blank := backgroundStyle.Render(" ")

	k8s := backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: theme.Colors.Bright.Blue,
		Dark:  theme.Colors.Normal.Blue,
	}).Render("\uE81D (alacritty-env)")

	prompt := backgroundStyle.Padding(1, 1).Width(m.maxWidth).Render(alacritty + blank + on + blank + branch + blank + stash + blank + program + blank + via + blank + k8s)

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		themeNameStyle.Render(m.theme.name),
		colorBlockStyle.Render(colorBlock),
		prompt,
	)

	return backgroundStyle.Height(m.maxHeight).Width(m.maxWidth).Render(view)
}

type colorBlock struct {
	normalColor string
	brightColor string
	name        string
}

func renderColorBlock(background, foreground string, colors []colorBlock) string {
	var allElements []string
	allElements = append(allElements,
		lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Foreground(lipgloss.Color(foreground)).Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(7).Height(3).Foreground(lipgloss.Color(foreground)).Render("Normal"),
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(7).Height(3).Foreground(lipgloss.Color(foreground)).Render("Bright"),
			),
		),
		lipgloss.NewStyle().Width(2).Height(9).Background(lipgloss.Color(background)).Render(""),
	)

	for i, color := range colors {
		s := lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Foreground(lipgloss.Color(foreground)).Width(7).Align(lipgloss.Center).Render(color.name),
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(color.normalColor)).Width(7).Height(3).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(color.brightColor)).Width(7).Height(3).Render(""),
			),
		)

		if i > 0 {
			allElements = append(allElements, lipgloss.NewStyle().Width(2).Height(9).Background(lipgloss.Color(background)).Render(""))
		}
		allElements = append(allElements, s)
	}
	return lipgloss.NewStyle().Background(lipgloss.Color(background)).Render(lipgloss.JoinHorizontal(lipgloss.Left, allElements...))
}
