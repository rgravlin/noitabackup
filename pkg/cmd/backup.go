/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"github.com/rgravlin/noitabackup/pkg/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var numBackupsToKeep int

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the Noita save00 directory",
	Long: `Backs up the Noita save00 directory to %USERPROFILE%\NoitaBackup or
a specified destination directory through the environmental variable CONFIG_NOITA_DST_PATH.`,
	Run: func(cmd *cobra.Command, args []string) {
		lib.BackupNoita(false, numBackupsToKeep)
	},
}

func init() {
	rootCmd.PersistentFlags().IntVar(&numBackupsToKeep, "num-backups", 16, "Define the maximum number of backups to keep")
	err := viper.BindPFlag("num-backups", rootCmd.PersistentFlags().Lookup("num-backups"))
	if err != nil {
		log.Printf("error binding viper flag: %v", err)
		return
	}
	rootCmd.AddCommand(backupCmd)
}
