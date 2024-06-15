package cmd

import (
	"fmt"
	"github.com/rgravlin/noitabackup/pkg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validateCommandOptions(cmd *cobra.Command, args []string) error {
	if numBackupsToKeep > ConfigMaxNumBackupsToKeep || numBackupsToKeep <= 0 {
		return fmt.Errorf("number of backups to keep must be between 1 and 100")
	}

	if path, err := internal.GetSourcePath(viper.GetString("source-path")); err != nil {
		return fmt.Errorf("error getting source path: %v", err)
	} else {
		viper.Set("source-path", path)
	}

	if path, err := internal.GetDestinationPath(viper.GetString("destination-path")); err != nil {
		return fmt.Errorf("error getting destination path: %v", err)
	} else {
		viper.Set("destination-path", path)
	}

	return nil
}
