package pkg

import "testing"

func TestDecodeTheme(t *testing.T) {
	theme := decodeTheme(catppuccinFrappe)
	t.Log(theme)
}
