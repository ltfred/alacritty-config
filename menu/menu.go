package menu

import (
	"fyne.io/fyne/v2"
	"github.com/ltfred/alacritty-config/config"
)

const preferenceCurrentTutorial = "currentTutorial"

// OnChangeFuncs is a slice of functions that can be registered
// to run when the user switches tutorial.
var OnChangeFuncs []func()

type Menu struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	Menus = map[string]Menu{
		"Window": {
			Title: "Window",
			Intro: "This section documents the [window](https://alacritty.org/config-alacritty.html#s6) table of the configuration file.",
			View:  makeWindowFormTab,
		},
		"Scrolling": {
			Title: "Scrolling",
			Intro: "This section documents the [scrolling](https://alacritty.org/config-alacritty.html#s22) table of the configuration file.",
			View:  makeScrollingFormTab,
		},
		"Font": {
			Title: "Font",
			Intro: "This section documents the [font](https://alacritty.org/config-alacritty.html#s25) table of the configuration file.",
			View:  makeFontFormTab,
		},
		"Theme": {
			Title: "Theme",
			Intro: "This section documents the [colors](https://alacritty.org/config-alacritty.html#s34) table of the configuration file.",
			View:  themeScreen,
		},
		"Selection": {
			Title: "Selection",
			Intro: "This section documents the [selection](https://alacritty.org/config-alacritty.html#s54) table of the configuration file.",
			View:  makeSelectionFormTab,
		},
		"Cursor": {
			Title: "Cursor",
			Intro: "This section documents the [cursor](https://alacritty.org/config-alacritty.html#s57) table of the configuration file.",
			View:  makeCursorFormTab,
		},
		"Terminal": {
			Title: "Terminal",
			Intro: "This section documents the [terminal](https://alacritty.org/config-alacritty.html#s68) table of the configuration file.",
			View:  makeTerminalFormTab,
		},
	}
	MenuIndex = map[string][]string{
		"": {"Window", "Scrolling", "Font", "Theme", "Selection", "Cursor", "Terminal"},
	}
)

func writeConfig() {
	msg := "Cursor setting succeeded"
	if err := config.WriteConfig(config.Cfg); err != nil {
		msg = "Cursor setting failed"
	}
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Alacritty config",
		Content: msg,
	})
}
