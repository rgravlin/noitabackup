/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"github.com/rgravlin/noitabackup/pkg/lib"
	"github.com/spf13/cobra"
	"log"
)

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch the Noita game from Steam",
	Long:  `Launches the Noita Steam game`,
	Run: func(cmd *cobra.Command, args []string) {
		err := lib.LaunchNoita()
		if err != nil {
			log.Printf("failed to launch noita: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
}
