/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"github.com/rgravlin/noitabackup/pkg/lib"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the Noita save00 directory",
	Long: `Backs up the Noita save00 directory to %USERPROFILE%\NoitaBackup or
a specified destination directory through the environmental variable CONFIG_NOITA_DST_PATH.`,
	Run: func(cmd *cobra.Command, args []string) {
		lib.BackupNoita()
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
