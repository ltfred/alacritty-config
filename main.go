package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config-gui/menu"
	"net/url"
)

var topWindow fyne.Window

func main() {
	a := app.NewWithID("io.fyne.demo")
	//a.SetIcon(data.FyneLogo)
	//makeTray(a)
	//logLifecycle(a)
	w := a.NewWindow("Fyne Demo")
	//topWindow = w

	w.SetMainMenu(makeMenu(a, w))
	w.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("Component name")

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
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := fyne.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := fyne.NewMenuItem("Other", nil)
	mailItem := fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") })
	mailItem.Icon = theme.MailComposeIcon()
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		mailItem,
	)
	fileItem := fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") })
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") })
	dirItem.Icon = theme.FolderIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		fileItem,
		dirItem,
		otherItem,
	)

	openSettings := func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	showAbout := func() {
		w := a.NewWindow("About")
		w.SetContent(widget.NewLabel("About Fyne Demo app..."))
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
	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	file.Items = append(file.Items, aboutItem)
	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
	checkedItem.Action = func() {
		checkedItem.Checked = !checkedItem.Checked
		main.Refresh()
	}
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
