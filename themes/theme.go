package themes

import _ "embed"

//go:embed cobalt2.toml
var colbalt2 []byte

var Themes = []string{
	"cobalt2",
}

var ThemesMap = map[string][]byte{
	"cobalt2": colbalt2,
}
