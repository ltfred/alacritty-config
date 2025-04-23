package menu

import (
	"fmt"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
	"github.com/ltfred/alacritty-config/themes"
)

const preferenceCurrentTutorial = "currentTutorial"

// OnChangeFuncs is a slice of functions that can be registered
// to run when the user switches tutorial.
var OnChangeFuncs []func()

type Menu struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	Menus = map[string]Menu{
		"Window": {
			Title: "Window",
			Intro: "This section documents the [window](https://alacritty.org/config-alacritty.html#s6) table of the configuration file.",
			View:  makeWindowFormTab,
		},
		"Scrolling": {
			Title: "Scrolling",
			Intro: "This section documents the [scrolling](https://alacritty.org/config-alacritty.html#s22) table of the configuration file.",
			View:  makeScrollingFormTab,
		},
		"Font": {
			Title: "Font",
			Intro: "This section documents the [font](https://alacritty.org/config-alacritty.html#s25) table of the configuration file.",
			View:  makeFontFormTab,
		},
		"Theme": {
			Title: "Theme",
			Intro: "This section documents the [colors](https://alacritty.org/config-alacritty.html#s34) table of the configuration file.",
			View:  themeScreen,
		},
		"Selection": {
			Title: "Selection",
			Intro: "This section documents the [selection](https://alacritty.org/config-alacritty.html#s54) table of the configuration file.",
			View:  makeSelectionFormTab,
		},
		"Cursor": {
			Title: "Cursor",
			Intro: "This section documents the [cursor](https://alacritty.org/config-alacritty.html#s57) table of the configuration file.",
			View:  makeCursorFormTab,
		},
		"Terminal": {
			Title: "Terminal",
			Intro: "This section documents the [terminal](https://alacritty.org/config-alacritty.html#s68) table of the configuration file.",
			View:  makeTerminalFormTab,
		},
	}
	MenuIndex = map[string][]string{
		"": {"Window", "Scrolling", "Font", "Theme", "Selection", "Cursor", "Terminal"},
	}
)

func makeTerminalFormTab(_ fyne.Window) fyne.CanvasObject {
	shell := widget.NewEntry()
	shell.SetPlaceHolder("/bin/zsh")
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Shell", Widget: shell},
		},
		OnCancel: func() {},
		OnSubmit: func() {},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}

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

func makeSelectionFormTab(_ fyne.Window) fyne.CanvasObject {
	// semanticEscapeChars
	semanticEscapeChars := widget.NewEntry()
	setSemanticEscapeChars := func() {
		if config.Cfg.Selection.SemanticEscapeChars != "" {
			semanticEscapeChars.SetText(config.Cfg.Selection.SemanticEscapeChars)
		} else {
			semanticEscapeChars.SetPlaceHolder(",│`|:\"' ()[]{}<>\t")
		}
	}
	setSemanticEscapeChars()

	// saveToClipboard
	saveToClipboard := widget.NewRadioGroup([]string{"true", "false"}, func(s string) {
	})
	saveToClipboard.Horizontal = true
	setSaveToClipboard := func() {
		if config.Cfg.Selection.SaveToClipboard {
			saveToClipboard.SetSelected("true")
		} else {
			saveToClipboard.SetSelected("false")
		}
	}
	setSaveToClipboard()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Semantic escape chars", Widget: semanticEscapeChars},
			{Text: "Save to clipboard", Widget: saveToClipboard},
		},
		OnCancel: func() {
			setSemanticEscapeChars()
			setSaveToClipboard()
		},
		OnSubmit: func() {
			if semanticEscapeChars.Text == "" {
				semanticEscapeChars.SetText(",│`|:\"' ()[]{}<>\t")
			}
			config.Cfg.Selection.SemanticEscapeChars = semanticEscapeChars.Text
			if saveToClipboard.Selected == "true" {
				config.Cfg.Selection.SaveToClipboard = true
			} else {
				config.Cfg.Selection.SaveToClipboard = false
			}
			writeConfig()
		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}

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

type browser struct {
	current  string
	themes   []string
	themeMap map[string]fyne.Resource
	image    *widget.Card
}

func (b *browser) setTheme(theme string) {
	if theme == "" {
		return
	}
	b.current = theme
	b.image.SetImage(canvas.NewImageFromResource(b.themeMap[theme]))
}

func themeScreen(_ fyne.Window) fyne.CanvasObject {
	b := &browser{}
	b.setThemeList(themes.LoadThemes())
	b.current = b.themes[0]
	b.image = widget.NewCard("", "", nil)

	themeSelect := widget.NewSelect(b.themes, func(name string) {})
	themeSelect.SetSelected(b.current)
	b.setTheme(b.current)

	label := widget.NewLabelWithStyle("Preview", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	themeSelect.OnChanged = func(s string) {
		b.setTheme(s)
	}
	sButton := widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
		fmt.Println(themeSelect.Selected)
	})
	sButton.Importance = widget.HighImportance
	cButton := widget.NewButtonWithIcon("Reset", theme.CancelIcon(), func() {
		themeSelect.SetSelected("Catppiuth")
		b.setTheme("Catppiuth")
	})
	box := container.NewBorder(nil, nil, nil, container.NewHBox(cButton, sButton))

	return container.NewVBox(themeSelect, box, label, b.image)
}

func (b *browser) setThemeList(themes []themes.ThemeInfo) {
	ret := make([]string, len(themes))
	themeMap := make(map[string]fyne.Resource)
	for i := range themes {
		ret[i] = themes[i].Name
		themeMap[themes[i].Name] = themes[i].Image
	}
	b.themeMap = themeMap
	b.themes = ret
}

func MakeNav(setTutorial func(menu Menu), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return MenuIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := MenuIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Menus[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := Menus[uid]; ok {
				for _, f := range OnChangeFuncs {
					f()
				}
				OnChangeFuncs = nil // Loading a page registers a new cleanup.

				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func writeConfig() {
	msg := "Cursor setting succeeded"
	if err := config.WriteConfig(config.Cfg); err != nil {
		msg = "Cursor setting failed"
	}
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Alacritty config",
		Content: msg,
	})
}
