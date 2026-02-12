package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

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

func ParseConfig() map[string]any {
	path := filepath.Join(getConfigDir(), "alacritty.toml")
	var config any
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil
	}

	return config.(map[string]any)
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

	config := ParseConfig()

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
