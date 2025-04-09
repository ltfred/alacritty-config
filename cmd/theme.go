package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/ltfred/alacritty-config/config"
	"github.com/ltfred/alacritty-config/themes"
	"github.com/ltfred/alacritty-config/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(themeCmd)
}

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Change alacrittty theme",
	Run:   theme,
}

func theme(cmd *cobra.Command, args []string) {
	i, err := displayTheme()
	if err != nil {
		color.Error.Prompt(err.Error())
		return
	}
	c := config.Config{}
	path, _ := utils.GetConfigPath()
	readConfig, err := c.ReadConfig(path)
	if err != nil {
		color.Error.Prompt(err.Error())
		return
	}
	th := themes.Themes[i]
	var colors config.Config
	err = toml.Unmarshal(themes.ThemesMap[th], &colors)
	if err != nil {
		color.Error.Prompt(err.Error())
		return
	}
	readConfig.Colors = colors.Colors
	err = readConfig.WriteConfig()
	if err != nil {
		color.Error.Prompt(err.Error())
		return
	}
	color.Info.Prompt("Alacritty theme change to %s", th)

	return
}

func colorPrint(col string) string {
	if len(col) > 2 {
		col = col[1:]
	}

	return color.Sprintf("<bg=%s>llo</>", col)
}

func displayTheme() (int, error) {
	idx, err := fuzzyfinder.FindMulti(
		themes.Themes,
		func(i int) string {
			return themes.Themes[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			th := themes.Themes[i]
			var colors config.Config
			err := toml.Unmarshal(themes.ThemesMap[th], &colors)
			if err != nil {
				color.Error.Prompt(err.Error())
				return ""
			}

			return fmt.Sprintf("        Normal Bright\n  Black %s     %s\n    Red %s     %s\n  Green %s     %s\n"+
				" Yellow %s     %s\n   Blue %s     %s\nMagenta %s     %s\n   Cyan %s     %s\n  White %s     %s",
				colorPrint(colors.Colors.Normal.Black), colorPrint(colors.Colors.Bright.Black),
				colorPrint(colors.Colors.Normal.Red), colorPrint(colors.Colors.Bright.Red),
				colorPrint(colors.Colors.Normal.Green), colorPrint(colors.Colors.Bright.Green),
				colorPrint(colors.Colors.Normal.Yellow), colorPrint(colors.Colors.Bright.Yellow),
				colorPrint(colors.Colors.Normal.Blue), colorPrint(colors.Colors.Bright.Blue),
				colorPrint(colors.Colors.Normal.Magenta), colorPrint(colors.Colors.Bright.Magenta),
				colorPrint(colors.Colors.Normal.Cyan), colorPrint(colors.Colors.Bright.Cyan),
				colorPrint(colors.Colors.Normal.White), colorPrint(colors.Colors.Bright.White),
			)
		}),
	)
	if err != nil {
		return 0, err
	}

	return idx[0], nil
}
