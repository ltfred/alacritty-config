package cmd

import (
	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/ltfred/alacritty-config/pkg/config"
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
	idx, err := fuzzyfinder.FindMulti(
		themes.Themes,
		func(i int) string {
			return themes.Themes[i]
		})
	if err != nil {
		color.Red(err.Error())
		return
	}
	c := config.Config{}
	path, _ := utils.GetConfigPath()
	readConfig, err := c.ReadConfig(path)
	if err != nil {
		color.Red(err.Error())
		return
	}
	th := themes.Themes[idx[0]]
	var colors config.Config
	err = toml.Unmarshal(themes.ThemesMap[th], &colors)
	if err != nil {
		color.Red(err.Error())
		return
	}
	readConfig.Colors = colors.Colors
	err = readConfig.WriteConfig()
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.Green("alacritty theme change to %s", th)

	return
}
