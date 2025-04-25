package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/ltfred/alacritty-config/config"
)

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
