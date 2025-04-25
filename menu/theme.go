package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/themes"
	"image/color"
)

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
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
