package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
	"runtime"
	"strconv"
)

func makeWindowFormTab(_ fyne.Window) fyne.CanvasObject {
	// title
	title := widget.NewEntry()
	setTitle := func() {
		if config.Cfg.Window.Title != "" {
			title.SetText(config.Cfg.Window.Title)
		} else {
			title.SetPlaceHolder("Alacritty")
		}
	}
	setTitle()
	// columns
	column := widget.NewEntry()
	setColumns := func() {
		if config.Cfg.Window.Dimensions.Columns != 0 {
			column.SetText(strconv.Itoa(int(config.Cfg.Window.Dimensions.Columns)))
		} else {
			column.SetPlaceHolder("180")
		}
	}
	setColumns()
	column.Validator = func(s string) error {
		if s != "" {
			if _, err := strconv.ParseInt(s, 10, 64); err != nil {
				return fmt.Errorf("column must be a number")
			}
		}
		return nil
	}
	// lines
	line := widget.NewEntry()
	setLines := func() {
		if config.Cfg.Window.Dimensions.Lines != 0 {
			line.SetText(strconv.Itoa(int(config.Cfg.Window.Dimensions.Lines)))
		} else {
			line.SetPlaceHolder("50")
		}
	}
	setLines()
	line.Validator = func(s string) error {
		if s != "" {
			if _, err := strconv.ParseInt(s, 10, 64); err != nil {
				return fmt.Errorf("line must be a number")
			}
		}
		return nil
	}

	// decorations
	decorationsChoices := []string{"Full", "None"}
	if runtime.GOOS == "darwin" {
		decorationsChoices = append(decorationsChoices, "Transparent", "Buttonless")
	}
	decorations := widget.NewSelect(decorationsChoices, func(s string) {})
	setDecorations := func() {
		if config.Cfg.Window.Decorations != "" {
			decorations.SetSelected(config.Cfg.Window.Decorations)
		} else {
			decorations.SetSelected("Full")
		}
	}
	setDecorations()

	// opacity
	var f float64
	if config.Cfg.Window.Opacity != 0 {
		f = float64(config.Cfg.Window.Opacity)
	} else {
		f = 1
	}
	data := binding.BindFloat(&f)
	slide := widget.NewSliderWithData(0, 1, data)
	slide.Step = 0.1
	entry := widget.NewEntryWithData(binding.FloatToStringWithFormat(data, "%.1f"))
	entry.Validator = func(s string) error {
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			return fmt.Errorf("opacity must be a float")
		}
		if f < 0 || f > 1 {
			return fmt.Errorf("opacity must be between 0 and 1.0")
		}
		return nil
	}
	// startup mode
	startupModeChoices := []string{"Windowed", "Maximized", "Fullscreen"}
	if runtime.GOOS == "darwin" {
		startupModeChoices = append(startupModeChoices, "SimpleFullscreen")
	}
	startupMode := widget.NewSelect(startupModeChoices, func(s string) {})
	setStartupMode := func() {
		if config.Cfg.Window.StartupMode != "" {
			startupMode.SetSelected(config.Cfg.Window.StartupMode)
		} else {
			startupMode.SetSelected("Windowed")
		}
	}
	setStartupMode()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: title},
			{Text: "Columns", Widget: column},
			{Text: "Lines", Widget: line},
			{Text: "Decorations", Widget: decorations},
			{Text: "Opacity", Widget: entry},
			{Text: "", Widget: slide},
			{Text: "Startup Mode", Widget: startupMode},
		},
		OnCancel: func() {
			setTitle()
			setColumns()
			setLines()
			setDecorations()
			setStartupMode()
		},
		OnSubmit: func() {
			if title.Text == "" {
				title.SetText("Alacritty")
			}
			if column.Text == "" {
				column.SetText("180")
			}
			if line.Text == "" {
				line.SetText("50")
			}
			config.Cfg.Window.Title = title.Text
			config.Cfg.Window.Dimensions.Columns, _ = strconv.ParseInt(column.Text, 10, 64)
			config.Cfg.Window.Dimensions.Lines, _ = strconv.ParseInt(line.Text, 10, 64)
			config.Cfg.Window.Decorations = decorations.Selected
			config.Cfg.Window.Opacity = float32(f)
			config.Cfg.Window.StartupMode = startupMode.Selected
			writeConfig()
		},
	}

	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}
