package config

type (
	Window struct {
		Dimensions     Dimensions `toml:"dimensions,omitempty"`
		Padding        FontOffset `toml:"padding,omitempty"`
		DynamicPadding bool       `toml:"dynamic_padding,omitempty"`
		Decorations    string     `toml:"decorations,omitempty"`
		Opacity        float32    `toml:"opacity,omitempty"`
		Blur           bool       `toml:"blur,omitempty"`
		StartupMode    string     `toml:"startup_mode,omitempty"`
		Title          string     `toml:"title,omitempty"`
		DynamicTitle   bool       `toml:"dynamic_title,omitempty"`
	}
	Dimensions struct {
		Columns int `toml:"columns,omitempty"`
		Lines   int `toml:"lines,omitempty"`
	}
)

type Scrolling struct {
	History    int32 `toml:"history"`
	Multiplier int32 `toml:"multiplier"`
}

type (
	Colors struct {
		Primary                      Primary         `toml:"primary,omitempty"`
		Cursor                       ColorsCursor    `toml:"cursor,omitempty"`
		ViModeCursor                 ColorsCursor    `toml:"vi_mode_cursor,omitempty"`
		Search                       ColorsSearch    `toml:"search,omitempty"`
		Hints                        ColorsHints     `toml:"hints,omitempty"`
		LineIndicator                ColorsGround    `toml:"line_indicator,omitempty"`
		FooterBar                    ColorsGround    `toml:"footer_bar,omitempty"`
		Selection                    ColorsSelection `toml:"selection,omitempty"`
		Normal                       Normal          `toml:"normal,omitempty"`
		Bright                       Normal          `toml:"bright,omitempty"`
		Dim                          Normal          `toml:"dim,omitempty"`
		IndexedColors                []IndexedColors `toml:"indexed_colors,omitempty"`
		TransparentBackgroundColors  bool            `toml:"transparent_background_colors,omitempty"`
		DrawBoldTextWithBrightColors bool            `toml:"draw_bold_text_with_bright_colors,omitempty"`
	}

	ColorsSearch struct {
		Matches      ColorsGround `toml:"matches,omitempty"`
		FocusedMatch ColorsGround `toml:"focused_match,omitempty"`
	}

	ColorsGround struct {
		Foreground string `toml:"foreground,omitempty"`
		Background string `toml:"background,omitempty"`
	}
	ColorsHints struct {
		Start ColorsGround `toml:"start,omitempty"`
		End   ColorsGround `toml:"end,omitempty"`
	}
	ColorsSelection struct {
		Text       string `toml:"text,omitempty"`
		Background string `toml:"background,omitempty"`
	}
	IndexedColors struct {
		Index int    `toml:"index,omitempty"`
		Color string `toml:"color,omitempty"`
	}

	Primary struct {
		Foreground       string `toml:"foreground,omitempty"`
		Background       string `toml:"background,omitempty"`
		DimForeground    string `toml:"dim_foreground,omitempty"`
		BrightForeground string `toml:"bright_foreground,omitempty"`
	}

	Normal struct {
		Black   string `toml:"black,omitempty"`
		Red     string `toml:"red,omitempty"`
		Green   string `toml:"green,omitempty"`
		Yellow  string `toml:"yellow,omitempty"`
		Blue    string `toml:"blue,omitempty"`
		Magenta string `toml:"magenta,omitempty"`
		Cyan    string `toml:"cyan,omitempty"`
		White   string `toml:"white,omitempty"`
	}

	ColorsCursor struct {
		Text   string `toml:"text,omitempty"`
		Cursor string `toml:"cursor,omitempty"`
	}
)

type (
	Font struct {
		Normal            FontNormal `toml:"normal,omitempty"`
		Bold              FontNormal `toml:"bold,omitempty"`
		Italic            FontNormal `toml:"italic,omitempty"`
		BoldItalic        FontNormal `toml:"bold_italic,omitempty"`
		Size              float32    `toml:"size,omitempty"`
		Offset            FontOffset `toml:"offset,omitempty"`
		GlyphOffset       FontOffset `toml:"glyph_offset,omitempty"`
		BuiltinBoxDrawing bool       `toml:"builtin_box_drawing,omitempty"`
	}

	FontNormal struct {
		Family string `toml:"family,omitempty"`
		Style  string `toml:"style,omitempty"`
	}
	FontOffset struct {
		X int `toml:"x,omitempty"`
		Y int `toml:"y,omitempty"`
	}
)

type (
	Cursor struct {
		Style CursorStyle `toml:"style,omitempty"`
	}

	CursorStyle struct {
		Shape    string `toml:"shape,omitempty"`
		Blinking string `toml:"blinking,omitempty"`
	}
)

type (
	Bell struct {
		Animation string `toml:"animation,omitempty"`
		Duration  int32  `toml:"duration,omitempty"`
		Color     string `toml:"color,omitempty"`
	}
)

type Selection struct {
	SemanticEscapeChars string `toml:"semantic_escape_chars,omitempty"`
	SaveToClipboard     bool   `toml:"save_to_clipboard,omitempty"`
}
