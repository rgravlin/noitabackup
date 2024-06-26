/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"github.com/rgravlin/noitabackup/pkg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the Noita save00 directory",
	Long: `Backs up the Noita save00 directory to %USERPROFILE%\NoitaBackup or a specified destination directory
through the environmental variable CONFIG_NOITA_DST_PATH.`,
	PreRunE: validateCommandOptions,
	Run: func(cmd *cobra.Command, args []string) {
		backup := internal.NewBackup(
			false,
			viper.GetBool(internal.ViperAutoLaunch),
			viper.GetInt(internal.ViperNumBackups),
			viper.GetString(internal.ViperSourcePath),
			viper.GetString(internal.ViperDestinationPath),
		)
		backup.BackupNoita()
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
