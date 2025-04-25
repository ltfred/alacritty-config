package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
	"strconv"
)

func makeCursorFormTab(_ fyne.Window) fyne.CanvasObject {
	// style
	style := widget.NewSelect([]string{"Block", "Underline", "Beam"}, func(s string) {})
	setStyle := func() {
		if config.Cfg.Cursor.Style.Shape != "" {
			style.SetSelected(config.Cfg.Cursor.Style.Shape)
		} else {
			style.SetSelected("Block")
		}
	}
	setStyle()
	// blinking
	blinking := widget.NewSelect([]string{"Never", "Off", "On", "Always"}, func(s string) {})
	setBlinking := func() {
		if config.Cfg.Cursor.Style.Shape != "" {
			blinking.SetSelected(config.Cfg.Cursor.Style.Blinking)
		} else {
			blinking.SetSelected("Off")
		}
	}
	setBlinking()
	// blinkInterval
	blinkInterval := widget.NewEntry()
	setBlinkInterval := func() {
		if config.Cfg.Cursor.BlinkInterval != 0 {
			blinkInterval.SetText(strconv.Itoa(config.Cfg.Cursor.BlinkInterval))
		} else {
			blinkInterval.SetPlaceHolder("750")
		}
	}
	setBlinkInterval()
	blinkInterval.Validator = func(s string) error {
		if s != "" {
			if _, err := strconv.ParseInt(s, 10, 64); err != nil {
				return fmt.Errorf("blink interval must be a number")
			}
		}
		return nil
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Style", Widget: style},
			{Text: "Blinking", Widget: blinking},
			{Text: "Blink interval", Widget: blinkInterval},
		},
		OnCancel: func() {
			setStyle()
			setBlinking()
			setBlinkInterval()
		},
		OnSubmit: func() {
			config.Cfg.Cursor.Style.Shape = style.Selected
			config.Cfg.Cursor.Style.Blinking = blinking.Selected
			if blinkInterval.Text != "" {
				i, _ := strconv.ParseInt(blinkInterval.Text, 10, 64)
				config.Cfg.Cursor.BlinkInterval = int(i)
			} else {
				blinkInterval.SetText("750")
				config.Cfg.Cursor.BlinkInterval = 750
			}
			writeConfig()
		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}
