package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

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
