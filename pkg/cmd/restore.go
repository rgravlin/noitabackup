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
	"log"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore the latest backed up Noita save",
	Long: `Restores the latest backed up Noita save to the save00 directory.  Preserves your current
save by deleting save00.bak and renaming save00 to save00.bak.  It then restores the latest save
file to the save00 directory.`,
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
		err := lib.RestoreNoita("", false)
		if err != nil {
			log.Printf("failed to restore noita backup: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
