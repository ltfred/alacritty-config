package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config-gui/themes"
	"runtime"
	"strconv"
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
			{Text: "Shell", Widget: shell, HintText: ""},
		},
		OnCancel: func() {
		},
		OnSubmit: func() {
		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}

func makeCursorFormTab(_ fyne.Window) fyne.CanvasObject {
	// style
	style := widget.NewSelect([]string{"Block", "Underline", "Beam"}, func(s string) {})
	style.SetSelected("Block")
	//style.Horizontal = true
	// blinking
	blinking := widget.NewSelect([]string{"Never", "Off", "On", "Always"}, func(s string) {})
	blinking.SetSelected("Off")
	//blinking.Horizontal = true
	blinkInterval := widget.NewEntry()
	blinkInterval.SetPlaceHolder("750")
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
			{Text: "Style", Widget: style, HintText: ""},
			{Text: "Blinking", Widget: blinking, HintText: ""},
			{Text: "Blink interval", Widget: blinkInterval, HintText: ""},
		},
		OnCancel: func() {
			style.SetSelected("Block")
			blinking.SetSelected("Off")
			blinkInterval.SetText("750")
		},
		OnSubmit: func() {

		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}

func makeSelectionFormTab(_ fyne.Window) fyne.CanvasObject {
	// semanticEscapeChars
	semanticEscapeChars := widget.NewEntry()
	semanticEscapeChars.SetPlaceHolder(",│`|:\\\"' ()[]{}<>\\t")
	// saveToClipboard
	saveToClipboard := widget.NewRadioGroup([]string{"true", "false"}, func(s string) {
	})
	saveToClipboard.SetSelected("false")
	saveToClipboard.Horizontal = true

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Semantic escape chars", Widget: semanticEscapeChars, HintText: ""},
			{Text: "Save to clipboard", Widget: saveToClipboard, HintText: ""},
		},
		OnCancel: func() {
			semanticEscapeChars.SetText(",│`|:\\\"' ()[]{}<>\\t")
			saveToClipboard.SetSelected("false")
		},
		OnSubmit: func() {

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
	family.SetPlaceHolder(defaultFont)
	// size
	size := widget.NewEntry()
	size.SetPlaceHolder("11.25")
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
			{Text: "Family", Widget: family, HintText: ""},
			{Text: "Size", Widget: size, HintText: ""},
		},
		OnCancel: func() {
			family.SetText(defaultFont)
			size.SetText("11.25")
		},
		OnSubmit: func() {

		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}

func makeWindowFormTab(_ fyne.Window) fyne.CanvasObject {
	// title
	title := widget.NewEntry()
	title.SetPlaceHolder("Alacritty")
	// columns
	column := widget.NewEntry()
	column.SetPlaceHolder("180")
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
	line.SetPlaceHolder("50")
	line.Validator = func(s string) error {
		if s != "" {
			if _, err := strconv.ParseInt(s, 10, 64); err != nil {
				return fmt.Errorf("line must be a number")
			}
		}
		return nil
	}
	// decorations
	decorations := widget.NewRadioGroup([]string{"Full", "None", "Transparent", "Buttonless"}, func(s string) {})
	decorations.SetSelected("Full")
	decorations.Horizontal = true
	// opacity
	f := 1.0
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

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: title, HintText: ""},
			{Text: "Columns", Widget: column, HintText: ""},
			{Text: "Lines", Widget: line, HintText: ""},
			{Text: "Decorations", Widget: decorations},
			{Text: "Opacity", Widget: entry},
			{Text: "", Widget: slide},
		},
		OnCancel: func() {
			title.SetText("Alacritty")
			column.SetText("180")
			line.SetText("50")
			decorations.SetSelected("Full")
			entry.SetText("1.0")
		},
		OnSubmit: func() {
		},
	}
	form.SubmitText, form.CancelText = "Save", "Reset"
	return form
}

func makeScrollingFormTab(_ fyne.Window) fyne.CanvasObject {
	// history
	history := widget.NewEntry()
	history.SetPlaceHolder("10000")
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
	multiplier.SetPlaceHolder("3")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "History", Widget: history},
			{Text: "Multiplier", Widget: multiplier},
		},
		OnCancel: func() {
			history.SetText("10000")
			multiplier.SetText("3")
		},
		OnSubmit: func() {},
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
