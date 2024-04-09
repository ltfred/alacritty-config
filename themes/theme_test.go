package themes

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTheme(t *testing.T) {
	dir, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	themes := make([]string, 0)
	themesMapS := make([]string, 0)
	var embedStr bytes.Buffer
	embedStr.WriteString(`package themes

import _ "embed"
`)

	for _, d := range dir {
		if strings.HasSuffix(d.Name(), ".toml") {
			themeN := fmt.Sprintf("\"%s\"", strings.Split(d.Name(), ".")[0])
			themes = append(themes, themeN)
			themesMapS = append(themesMapS, strings.Join([]string{themeN, strings.Split(d.Name(), ".")[0]}, ": "))
		}

		embedStr.WriteString(fmt.Sprintf(`
//go:embed %s
var %s []byte
`, d.Name(), strings.Split(d.Name(), ".")[0]))
	}

	embedStr.WriteString(fmt.Sprintf(`
var Themes = []string{
	%s,
}`, strings.Join(themes, ",\n")))

	embedStr.WriteString(fmt.Sprintf(`
var ThemesMap = map[string][]byte{
	%s,
}`, strings.Join(themesMapS, ",\n")))

	file, err := os.Create("theme.go")
	if err != nil {
		panic(err)
	}
	file.Write(embedStr.Bytes())
}
