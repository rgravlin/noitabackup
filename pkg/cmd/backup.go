/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/rgravlin/noitabackup/pkg/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var numBackupsToKeep int

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the Noita save00 directory",
	Long: `Backs up the Noita save00 directory to %USERPROFILE%\NoitaBackup or
a specified destination directory through the environmental variable CONFIG_NOITA_DST_PATH.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if numBackupsToKeep > ConfigMaxNumBackupsToKeep || numBackupsToKeep <= 0 {
			return fmt.Errorf("number of backups to keep must be between 1 and 100")
		}

		if path, err := lib.GetSourcePath(viper.GetString("source-path")); err != nil {
			return fmt.Errorf("error getting source path: %v", err)
		} else {
			viper.Set("source-path", path)
		}
		if path, err := lib.GetDestinationPath(viper.GetString("destination-path")); err != nil {
			return fmt.Errorf("error getting destination path: %v", err)
		} else {
			viper.Set("destination-path", path)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		lib.BackupNoita(false, numBackupsToKeep)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
