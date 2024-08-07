/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/rgravlin/noitabackup/pkg/internal"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName                   = "Noita Backup and Restore"
	ConfigDefaultSteamDir     = "C:\\Program Files (x86)\\Steam\\steam.exe"
	ConfigMaxNumBackupsToKeep = 64
	ConfigDefaultNumBackups   = 16
	ConfigMaxNumWorkers       = 32
	ConfigDefaultNumWorkers   = 4
)

var (
	cfgFile, sourcePath, destinationPath, steamPath string
	numBackupsToKeep, numCopyWorkers                int
	autoLaunch                                      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "noitabackup",
	Short: "A Noita backup and restore manager",
	Long: `A configurable Noita backup and restore manager and launcher.  Automates the tedious task of starting,
backing up, restoring, and restarting Noita.  Includes both a GUI and command line interface.`,
	PreRunE: validateCommandOptions,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			window := new(app.Window)
			window.Option(
				app.Title(appName),
				app.MaxSize(unit.Dp(internal.DefaultWidth), unit.Dp(internal.DefaultMinHeight)),
				app.MinSize(unit.Dp(internal.DefaultWidth), unit.Dp(internal.DefaultMinHeight)),
			)
			ui := internal.NewUI(viper.GetBool(internal.ViperAutoLaunch))
			err := ui.Run(window)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
		app.Main()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.noitabackup.yaml)")
	rootCmd.PersistentFlags().StringVar(&sourcePath, internal.ViperSourcePath, internal.GetDefaultSourcePath(), "source Noita save00 path")
	rootCmd.PersistentFlags().StringVar(&destinationPath, internal.ViperDestinationPath, internal.GetDefaultDestinationPath(), "destination backup path")
	rootCmd.PersistentFlags().StringVar(&steamPath, internal.ViperSteamPath, ConfigDefaultSteamDir, "path for your Steam executable")
	rootCmd.PersistentFlags().IntVar(&numBackupsToKeep, internal.ViperNumBackups, ConfigDefaultNumBackups, "maximum number of backups to keep")
	rootCmd.PersistentFlags().IntVar(&numCopyWorkers, internal.ViperNumWorkers, ConfigDefaultNumWorkers, "total number of go routine workers (advanced usage)")
	rootCmd.PersistentFlags().BoolVar(&autoLaunch, internal.ViperAutoLaunch, false, "auto-launch Noita after backup/restore operation")

	commands := []string{
		internal.ViperSourcePath,
		internal.ViperDestinationPath,
		internal.ViperNumBackups,
		internal.ViperNumWorkers,
		internal.ViperAutoLaunch,
		internal.ViperSteamPath,
	}

	for _, cmd := range commands {
		if err := viper.BindPFlag(cmd, rootCmd.PersistentFlags().Lookup(cmd)); err != nil {
			log.Printf("error binding viper flag: %v", err)
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".noitabackup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".noitabackup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			return
		}
	}
}
