package pkg

import (
	_ "embed"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/list"
)

//go:embed files/catppuccin-frappe.toml
var catppuccinFrappe string

//go:embed files/catppuccin-latte.toml
var catppuccinLatte string

//go:embed files/catppuccin-macchiato.toml
var catppuccinMacchiato string

//go:embed files/catppuccin-mocha.toml
var catppuccinMocha string

var themes = []theme{
	{
		name: "catppuccin-frappe",
		data: catppuccinFrappe,
	},
	{
		name: "catppuccin-latte",
		data: catppuccinLatte,
	},
	{
		name: "catppuccin-macchiato",
		data: catppuccinMacchiato,
	},
	{
		name: "catppuccin-mocha",
		data: catppuccinMocha,
	},
}

type theme struct {
	name string
	data string
}

func themeList() ([]list.Item, int) {
	l := make([]list.Item, 0, len(themes))
	maxWidth := 0
	for _, theme := range themes {
		l = append(l, item(theme.name))
		maxWidth = max(maxWidth, len(theme.name))
	}
	return l, maxWidth + 5
}

func decodeTheme(data string) Theme {
	var theme Theme
	_, _ = toml.Decode(data, &theme)
	return theme
}

type Theme struct {
	Colors struct {
		Primary struct {
			Background       string `toml:"background"`
			Foreground       string `toml:"foreground"`
			DimForeground    string `toml:"dim_foreground"`
			BrightForeground string `toml:"bright_foreground"`
		} `toml:"primary"`
		Cursor struct {
			Text   string `toml:"text"`
			Cursor string `toml:"cursor"`
		} `toml:"cursor"`
		ViModeCursor struct {
			Text   string `toml:"text"`
			Cursor string `toml:"cursor"`
		} `toml:"vi_mode_cursor"`
		Search struct {
			Matches struct {
				Foreground string `toml:"foreground"`
				Background string `toml:"background"`
			} `toml:"matches"`
			FocusedMatch struct {
				Foreground string `toml:"foreground"`
				Background string `toml:"background"`
			} `toml:"focused_match"`
		} `toml:"search"`
		FooterBar struct {
			Foreground string `toml:"foreground"`
			Background string `toml:"background"`
		} `toml:"footer_bar"`
		Hints struct {
			Start struct {
				Foreground string `toml:"foreground"`
				Background string `toml:"background"`
			} `toml:"start"`
			End struct {
				Foreground string `toml:"foreground"`
				Background string `toml:"background"`
			} `toml:"end"`
		} `toml:"hints"`
		Selection struct {
			Text       string `toml:"text"`
			Background string `toml:"background"`
		} `toml:"selection"`
		Normal struct {
			Black   string `toml:"black"`
			Red     string `toml:"red"`
			Green   string `toml:"green"`
			Yellow  string `toml:"yellow"`
			Blue    string `toml:"blue"`
			Magenta string `toml:"magenta"`
			Cyan    string `toml:"cyan"`
			White   string `toml:"white"`
		} `toml:"normal"`
		Bright struct {
			Black   string `toml:"black"`
			Red     string `toml:"red"`
			Green   string `toml:"green"`
			Yellow  string `toml:"yellow"`
			Blue    string `toml:"blue"`
			Magenta string `toml:"magenta"`
			Cyan    string `toml:"cyan"`
			White   string `toml:"white"`
		} `toml:"bright"`
		IndexedColors []struct {
			Index int    `toml:"index"`
			Color string `toml:"color"`
		} `toml:"indexed_colors"`
	} `toml:"colors"`
}
