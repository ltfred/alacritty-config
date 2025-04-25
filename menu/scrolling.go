package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
	"strconv"
)

func makeScrollingFormTab(_ fyne.Window) fyne.CanvasObject {
	// history
	history := widget.NewEntry()
	setHistory := func() {
		if config.Cfg.Scrolling.History != 0 {
			history.SetText(strconv.Itoa(int(config.Cfg.Scrolling.History)))
		} else {
			history.SetPlaceHolder("10000")
		}
	}
	setHistory()
	history.Validator = func(s string) error {
		if s != "" {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return fmt.Errorf("history must be a float")
			}
			if i < 0 || i > 10000 {
				return fmt.Errorf("history must be between 0 and 10000")
			}
		}
		return nil
	}
	// multiplier
	multiplier := widget.NewEntry()
	setMultiplier := func() {
		if config.Cfg.Scrolling.Multiplier != 0 {
			multiplier.SetText(strconv.Itoa(int(config.Cfg.Scrolling.Multiplier)))
		} else {
			multiplier.SetPlaceHolder("3")
		}
	}
	setMultiplier()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "History", Widget: history},
			{Text: "Multiplier", Widget: multiplier},
		},
		OnCancel: func() {
			setHistory()
			setMultiplier()
		},
		OnSubmit: func() {
			if history.Text == "" {
				history.SetText("10000")
			}
			if multiplier.Text == "" {
				multiplier.SetText("3")
			}
			config.Cfg.Scrolling.History, _ = strconv.ParseInt(history.Text, 10, 64)
			config.Cfg.Scrolling.Multiplier, _ = strconv.ParseInt(multiplier.Text, 10, 64)
		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"

	return form
}
