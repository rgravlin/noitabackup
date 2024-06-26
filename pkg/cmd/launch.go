/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"github.com/rgravlin/noitabackup/pkg/internal"
	"github.com/spf13/cobra"
	"log"
)

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch the Noita game from Steam",
	Long:  `Launches the Noita Steam game`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.LaunchNoita(false)
		if err != nil {
			log.Printf("%s: %v", internal.ErrLaunchingNoita, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
}
