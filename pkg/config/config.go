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
		panic(err)
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

	general := make(map[string]any)
	if _, ok := config["general"]; ok {
		general = config["general"].(map[string]any)
	}
	anies := make([]any, 0)
	if _, ok := general["import"]; ok {
		anies = general["import"].([]any)
	}
	anies = append(anies, themePath)
	general["import"] = anies
	config["general"] = general

	return WriteConfig(config)
}
