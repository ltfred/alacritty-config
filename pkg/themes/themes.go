package themes

import (
	_ "embed"

	"github.com/BurntSushi/toml"
)

//go:embed catppuccin-frappe.toml
var catppuccinFrappe string

//go:embed catppuccin-latte.toml
var catppuccinLatte string

//go:embed catppuccin-macchiato.toml
var catppuccinMacchiato string

//go:embed catppuccin-mocha.toml
var catppuccinMocha string

//go:embed snazzy.toml
var snazzy string

//go:embed solarized-dark.toml
var solarizedDark string

//go:embed solarized-light.toml
var solarizedLight string

//go:embed one-dark.toml
var oneDark string

var themes = []Theme{
	{
		Name: "catppuccin-frappe",
		data: catppuccinFrappe,
	},
	{
		Name: "catppuccin-latte",
		data: catppuccinLatte,
	},
	{
		Name: "catppuccin-macchiato",
		data: catppuccinMacchiato,
	},
	{
		Name: "catppuccin-mocha",
		data: catppuccinMocha,
	},
	{
		Name: "snazzy",
		data: snazzy,
	},
	{
		Name: "solarized-dark",
		data: solarizedDark,
	},
	{
		Name: "solarized-light",
		data: solarizedLight,
	},
	{
		Name: "one-dark",
		data: oneDark,
	},
}

type Theme struct {
	Name string
	data string

	Colors Colors `toml:"colors"`
}

func (t Theme) GetColorsOriginData() string {
	return t.data
}

func GetThemes() []Theme {
	for i := range themes {
		_, _ = toml.Decode(themes[i].data, &themes[i])
	}
	return themes
}

type Colors struct {
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
}

type Font struct {
	Normal     string  `toml:"normal"`
	Bold       string  `toml:"bold"`
	Italic     string  `toml:"italic"`
	BoldItalic string  `toml:"bold_italic"`
	Size       float64 `toml:"size"`
}

type Cursor struct {
	Style struct {
		Shape    string `toml:"shape"`
		Blinking string `toml:"blinking"`
	} `toml:"style"`
	BlinkInterval int `toml:"blink_interval"`
}
