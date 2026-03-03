package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Dimensions struct {
	Columns int `toml:"columns"`
	Lines   int `toml:"lines"`
}

type Window struct {
	Decorations string     `toml:"decorations"`
	Dimensions  Dimensions `toml:"dimensions"`
	StartMode   string     `toml:"startup_mode"`
	Opacity     float64    `toml:"opacity"`
	Title       string     `toml:"title"`
}

type Cursor struct {
	BlinkInterval int         `toml:"blink_interval"`
	Style         CursorStyle `toml:"style"`
}

type CursorStyle struct {
	Shape    string `toml:"shape"`
	BlinkIng string `toml:"blinking"`
}

type Font struct {
	Size       float64   `toml:"size"`
	Normal     FontStyle `toml:"normal"`
	Bold       FontStyle `toml:"bold"`
	Italic     FontStyle `toml:"italic"`
	BoldItalic FontStyle `toml:"bold_italic"`
}

type FontStyle struct {
	Family string `toml:"family"`
	Style  string `toml:"style"`
}

type Config struct {
	Window Window `toml:"window"`
	Cursor Cursor `toml:"cursor"`
	Font   Font   `toml:"font"`
}

func getConfigDir() string {
	switch runtime.GOOS {
	case "darwin", "linux":
		return filepath.Join(os.Getenv("HOME"), ".config", "alacritty")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "alacritty")
	default:
		return ""
	}
}

func GetConfigMap() map[string]any {
	path := filepath.Join(getConfigDir(), "alacritty.toml")
	var config any
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil
	}

	return config.(map[string]any)
}

func GetConfigStruct() Config {
	path := filepath.Join(getConfigDir(), "alacritty.toml")
	var config Config
	toml.DecodeFile(path, &config)
	return config
}

func WriteConfig(config map[string]any) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(getConfigDir(), "alacritty.toml"), buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func SetTheme(name, data string) error {
	themePath := filepath.Join(getConfigDir(), fmt.Sprintf("%s.toml", name))

	var colors any
	if _, err := toml.Decode(data, &colors); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(colors); err != nil {
		return err
	}
	if err := os.WriteFile(themePath, buf.Bytes(), 0644); err != nil {
		return err
	}

	config := GetConfigMap()

	general, ok := config["general"].(map[string]any)
	if !ok {
		general = make(map[string]any)
	}
	importList, ok := general["import"].([]any)
	if !ok {
		importList = make([]any, 0)
	}

	importList = append(importList, themePath)
	general["import"] = importList
	config["general"] = general

	return WriteConfig(config)
}
