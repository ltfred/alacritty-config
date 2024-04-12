package config

import (
	"os"

	"github.com/ltfred/alacritty-config/utils"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Import    []string  `toml:"import,omitempty"`
	Window    Window    `toml:"window,omitempty"`
	Colors    Colors    `toml:"colors,omitempty"`
	Font      Font      `toml:"font,omitempty"`
	Cursor    Cursor    `toml:"cursor,omitempty"`
	Scrolling Scrolling `toml:"scrolling,omitempty"`
	Bell      Bell      `toml:"bell,omitempty"`
	Selection Selection `toml:"selection,omitempty"`
}

func (c *Config) ReadConfig(path string) (Config, error) {
	readFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = toml.Unmarshal(readFile, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) WriteConfig() error {
	marshal, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	path, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, marshal, 0777)
}

func (c *Config) SetDefault() {
	c.Window = Window{
		Dimensions: Dimensions{
			Columns: 180,
			Lines:   50,
		},
		Padding:        FontOffset{},
		DynamicPadding: false,
		Decorations:    "Full",
		Opacity:        1.0,
		Blur:           false,
		StartupMode:    "Windowed",
		Title:          "Alacritty",
		DynamicTitle:   true,
	}

	c.Font = Font{
		Normal: FontNormal{
			Family: "",
			Style:  "Regular",
		},
		Bold:              FontNormal{},
		Italic:            FontNormal{},
		BoldItalic:        FontNormal{},
		Size:              15,
		Offset:            FontOffset{},
		GlyphOffset:       FontOffset{},
		BuiltinBoxDrawing: false,
	}

	c.Cursor = Cursor{
		Style: CursorStyle{
			Shape:    "Block",
			Blinking: "Off",
		},
		BlinkInterval: 750,
	}
}
