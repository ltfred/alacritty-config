package themes

import "fyne.io/fyne/v2"

type ThemeInfo struct {
	Name  string
	Image fyne.Resource
}

func LoadThemes() []ThemeInfo {
	resource1, _ := fyne.LoadResourceFromPath("/Users/liutao/Desktop/alacritty-config/themes/img.png")
	return []ThemeInfo{
		{"Catppiuth", resource1},
	}

}
