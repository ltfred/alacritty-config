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

	t               Theme
	backgroundStyle lipgloss.Style
}

func (m liveMode) Init() tea.Cmd {
	return nil
}

func (m liveMode) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m liveMode) View() string {
	m.t = decodeTheme(m.theme.data)
	m.backgroundStyle = lipgloss.NewStyle().Background(lipgloss.Color(m.t.Colors.Primary.Background))

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		m.renderColorName(),
		m.renderColorBlock(),
		m.renderPrompt(),
		m.renderArrow(),
	)

	return m.backgroundStyle.Height(m.maxHeight).Width(m.maxWidth).Render(view)
}

func (m liveMode) renderArrow() string {
	arrow := m.backgroundStyle.Bold(true).Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Green,
		Dark:  m.t.Colors.Normal.Green,
	}).Render("\U000F17A9 alacritty")
	cmd := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Primary.BrightForeground,
		Dark:  m.t.Colors.Primary.Foreground,
	}).Render(" themes")

	return m.backgroundStyle.Width(m.maxWidth).Render(lipgloss.JoinHorizontal(lipgloss.Center, arrow, cmd))
}

func (m liveMode) renderColorName() string {
	return m.backgroundStyle.Bold(true).Italic(true).
		Width(m.maxWidth).Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Primary.BrightForeground,
		Dark:  m.t.Colors.Primary.Foreground,
	}).Align(lipgloss.Center).Render(m.theme.name)
}

func (m liveMode) renderPrompt() string {
	alacritty := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Cyan,
		Dark:  m.t.Colors.Normal.Cyan,
	}).Render("alacritty")

	branch := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Blue,
		Dark:  m.t.Colors.Normal.Blue,
	}).Render("\uF418 master")

	stash := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Red,
		Dark:  m.t.Colors.Normal.Red,
	}).Render("[+]")

	program := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Yellow,
		Dark:  m.t.Colors.Normal.Yellow,
	}).Render("\uE627 v1.26")

	on := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Yellow,
		Dark:  m.t.Colors.Normal.Yellow,
	}).Render("on")
	via := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Yellow,
		Dark:  m.t.Colors.Normal.Yellow,
	}).Render("via")

	blank := m.backgroundStyle.Render(" ")

	k8s := m.backgroundStyle.Foreground(lipgloss.AdaptiveColor{
		Light: m.t.Colors.Bright.Blue,
		Dark:  m.t.Colors.Normal.Blue,
	}).Render("\uE81D (alacritty-env)")

	return m.backgroundStyle.Width(m.maxWidth).Render(alacritty + blank + on + blank + branch + blank + stash + blank + program + blank + via + blank + k8s)
}

type colorBlock struct {
	normalColor string
	brightColor string
	name        string
}

func (m liveMode) renderColorBlock() string {
	// 创建彩色方块
	colorNames := []string{"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White"}
	colors := make([]colorBlock, 0, len(colorNames))
	for _, name := range colorNames {
		normalColor := reflect.ValueOf(m.t.Colors.Normal).FieldByName(name).String()
		brightColor := reflect.ValueOf(m.t.Colors.Bright).FieldByName(name).String()
		colors = append(colors, colorBlock{
			normalColor: normalColor,
			brightColor: brightColor,
			name:        name,
		})
	}

	var allElements []string
	allElements = append(allElements,
		lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.backgroundStyle.Foreground(lipgloss.Color(m.t.Colors.Primary.Foreground)).Width(7).Render(""),
				m.backgroundStyle.Width(7).Render(""),
				m.backgroundStyle.Width(7).Height(3).Foreground(lipgloss.Color(m.t.Colors.Primary.Foreground)).Render("Normal"),
				m.backgroundStyle.Width(7).Render(""),
				m.backgroundStyle.Width(7).Height(3).Foreground(lipgloss.Color(m.t.Colors.Primary.Foreground)).Render("Bright"),
			),
		),
		m.backgroundStyle.Width(2).Height(9).Render(""),
	)

	for i, color := range colors {
		s := lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.backgroundStyle.Foreground(lipgloss.Color(m.t.Colors.Primary.Foreground)).Width(7).Align(lipgloss.Center).Render(color.name),
				m.backgroundStyle.Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(color.normalColor)).Width(7).Height(3).Render(""),
				m.backgroundStyle.Width(7).Render(""),
				lipgloss.NewStyle().Background(lipgloss.Color(color.brightColor)).Width(7).Height(3).Render(""),
			),
		)

		if i > 0 {
			allElements = append(allElements, m.backgroundStyle.Width(2).Height(9).Render(""))
		}
		allElements = append(allElements, s)
	}
	return m.backgroundStyle.Width(m.maxWidth).Padding(2, 0).Render(lipgloss.JoinHorizontal(lipgloss.Left, allElements...))
}
