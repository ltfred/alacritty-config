package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
	"runtime"
	"strconv"
)

func makeFontFormTab(_ fyne.Window) fyne.CanvasObject {
	defaultFont := ""
	switch runtime.GOOS {
	case "darwin":
		defaultFont = "Menlo"
	case "linux":
		defaultFont = "monospace"
	case "windows":
		defaultFont = "Consolas"
	}
	// family
	family := widget.NewEntry()
	setFamily := func() {
		if config.Cfg.Font.Normal.Family != "" {
			family.SetText(config.Cfg.Font.Normal.Family)
		} else {
			family.SetText(defaultFont)
		}
	}
	setFamily()
	// size
	size := widget.NewEntry()
	setSize := func() {
		if config.Cfg.Font.Size != 0 {
			size.SetText(fmt.Sprintf("%v", config.Cfg.Font.Size))
		} else {
			size.SetText("11.25")
		}
	}
	setSize()
	size.Validator = func(s string) error {
		if s != "" {
			if _, err := strconv.ParseFloat(s, 64); err != nil {
				return fmt.Errorf("size must be a float")
			}
		}
		return nil
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Family", Widget: family},
			{Text: "Size", Widget: size},
		},
		OnCancel: func() {
			setFamily()
			setSize()
		},
		OnSubmit: func() {
			if family.Text == "" {
				family.SetText(defaultFont)
			}
			if size.Text == "" {
				size.SetText("11.25")
			}
			s, _ := strconv.ParseFloat(size.Text, 64)
			config.Cfg.Font.Size = float32(s)
			config.Cfg.Font.Normal.Family = family.Text
			writeConfig()
		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}
