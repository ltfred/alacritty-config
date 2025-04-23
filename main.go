package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
	"github.com/ltfred/alacritty-config/menu"
)

func main() {
	config.SetConfig()
	a := app.NewWithID("io.fyne.demo")
	a.SetIcon(data.FyneLogo)
	//makeTray(a)
	//logLifecycle(a)
	w := a.NewWindow("Alacritty config")
	//topWindow = w

	w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("")

	intro := widget.NewRichTextFromMarkdown("")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t menu.Menu) {
		title.SetText(t.Title)
		title.Show()
		intro.ParseMarkdown(t.Intro)
		intro.Show()
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	split := container.NewHSplit(menu.MakeNav(setTutorial, true), tutorial)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(640, 460))

	w.ShowAndRun()

}

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	openSettings := func() {
		w := a.NewWindow("Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	showAbout := func() {
		w := a.NewWindow("About")
		w.SetContent(widget.NewLabel("Alacritty config"))
		w.Show()
	}
	aboutItem := fyne.NewMenuItem("About", showAbout)
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	cutShortcut := &fyne.ShortcutCut{Clipboard: a.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(cutShortcut, a.Clipboard(), w.Canvas().Focused())
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: a.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(copyShortcut, a.Clipboard(), w.Canvas().Focused())
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: a.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(pasteShortcut, a.Clipboard(), w.Canvas().Focused())
	})
	pasteItem.Shortcut = pasteShortcut

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://github.com/ltfred/alacritty-config")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("")
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	file.Items = append(file.Items, aboutItem)
	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator()),
		helpMenu,
	)
	return main
}

func shortcutFocused(s fyne.Shortcut, cb fyne.Clipboard, f fyne.Focusable) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = cb
	case *fyne.ShortcutCut:
		sh.Clipboard = cb
	case *fyne.ShortcutPaste:
		sh.Clipboard = cb
	}
	if focused, ok := f.(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
