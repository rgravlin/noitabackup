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

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore the latest backed up Noita save",
	Long: `Restores the latest backed up Noita save to the save00 directory or a specified source directory through the
environmental variable CONFIG_NOITA_SRC_PATH.  Preserves your current save by deleting save00.bak and renaming save00
to save00.bak.  It then restores the latest save file to the save00 directory.`,
	PreRunE: validateCommandOptions,
	Run: func(cmd *cobra.Command, args []string) {
		restore := internal.NewRestore(
			"latest",
			internal.NewBackup(
				false,
				viper.GetBool("auto-launch"),
				viper.GetInt("num-backups"),
				viper.GetString("source-path"),
				viper.GetString("destination-path"),
			),
		)
		restore.RestoreNoita()
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
