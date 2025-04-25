package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
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
