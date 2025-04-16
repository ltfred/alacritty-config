package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

var Cfg *Config

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

func ReadConfig(path string) (*Config, error) {
	readFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = toml.Unmarshal(readFile, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func WriteConfig(c *Config) error {
	marshal, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, marshal, 0777)
}

func SetConfig() {
	path, _ := GetConfigPath()
	c, err := ReadConfig(path)
	if err != nil {
		panic(err)
	}
	Cfg = c
}

func GetConfigPath() (string, error) {
	var path string
	switch runtime.GOOS {
	case "windows":
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, ".config/alacritty/alacritty.toml")
	}

	return path, nil
}
