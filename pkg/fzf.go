package pkg

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func FuzzyThemes() {
	p := tea.NewProgram(NewThemeChooseModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
