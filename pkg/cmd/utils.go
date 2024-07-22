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
		uiErr = fmt.Sprintf("%s: %d", internal.ErrNumBackups, numBackups)
	}

	numWorkers := viper.GetInt(internal.ViperNumWorkers)
	if numWorkers > ConfigMaxNumWorkers || numWorkers <= 0 {
		uiErr = fmt.Sprintf("%s: %d", internal.ErrNumWorkers, numWorkers)
	}

	if path, err := internal.GetSourcePath(viper.GetString(internal.ViperSourcePath)); err != nil {
		uiErr = fmt.Sprintf("%v", err)
	} else {
		viper.Set(internal.ViperSourcePath, path)
	}

	if path, err := internal.GetDestinationPath(viper.GetString(internal.ViperDestinationPath)); err != nil {
		uiErr = fmt.Sprintf("%v", err)
	} else {
		viper.Set(internal.ViperDestinationPath, path)
	}

	if path, err := internal.GetSteamPath(viper.GetString(internal.ViperSteamPath)); err != nil {
		uiErr = fmt.Sprintf("%v", err)
	} else {
		viper.Set(internal.ViperSteamPath, path)
	}

	if uiErr != "" {
		RunErrorUI(uiErr)
	}

	return nil
}

func RunErrorUI(error string) {
	go func() {
		window := new(app.Window)
		window.Option(
			app.Title(appName),
			app.MaxSize(unit.Dp(internal.ErrorWidth), unit.Dp(internal.ErrorHeight)),
			app.MinSize(unit.Dp(internal.ErrorWidth), unit.Dp(internal.ErrorHeight)),
		)
		ui := internal.NewErrorUI(error)
		err := ui.Run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
