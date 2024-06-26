package cmd

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/rgravlin/noitabackup/pkg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func validateCommandOptions(cmd *cobra.Command, args []string) error {
	var uiErr string

	numBackups := viper.GetInt(internal.ViperNumBackups)
	if numBackups > ConfigMaxNumBackupsToKeep || numBackups <= 0 {
		uiErr = fmt.Sprintf("%s: %d", internal.ErrNumBackups, viper.GetInt(internal.ViperNumBackups))
	}

	if path, err := internal.GetSourcePath(viper.GetString(internal.ViperSourcePath)); err != nil {
		uiErr = fmt.Sprintf("%s: %v", internal.ErrGettingSourcePath, err)
	} else {
		viper.Set(internal.ViperSourcePath, path)
	}

	if path, err := internal.GetDestinationPath(viper.GetString(internal.ViperDestinationPath)); err != nil {
		uiErr = fmt.Sprintf("%s: %v", internal.ErrGettingDestinationPath, err)
	} else {
		viper.Set(internal.ViperDestinationPath, path)
	}

	if path, err := internal.GetSteamPath(viper.GetString(internal.ViperSteamPath)); err != nil {
		uiErr = fmt.Sprintf("%s: %v", internal.ErrSteamPathNotExist, err)
	} else {
		viper.Set(internal.ViperSteamPath, path)
	}

	if uiErr != "" {
		go func() {
			window := new(app.Window)
			window.Option(
				app.Title(appName),
				app.MaxSize(unit.Dp(internal.ErrorWidth), unit.Dp(internal.ErrorHeight)),
				app.MinSize(unit.Dp(internal.ErrorWidth), unit.Dp(internal.ErrorHeight)),
			)
			ui := internal.NewErrorUI(uiErr)
			err := ui.Run(window)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		app.Main()
	}

	return nil
}
