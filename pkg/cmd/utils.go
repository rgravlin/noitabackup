package cmd

import (
	"fmt"
	"github.com/rgravlin/noitabackup/pkg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validateCommandOptions(cmd *cobra.Command, args []string) error {
	if numBackupsToKeep > ConfigMaxNumBackupsToKeep || numBackupsToKeep <= 0 {
		return fmt.Errorf(internal.ErrNumBackups)
	}

	if path, err := internal.GetSourcePath(viper.GetString(internal.ViperSourcePath)); err != nil {
		return fmt.Errorf("%s: %v", internal.ErrGettingSourcePath, err)
	} else {
		viper.Set(internal.ViperSourcePath, path)
	}

	if path, err := internal.GetDestinationPath(viper.GetString(internal.ViperDestinationPath)); err != nil {
		return fmt.Errorf("%s: %v", internal.ErrGettingDestinationPath, err)
	} else {
		viper.Set(internal.ViperDestinationPath, path)
	}

	return nil
}
