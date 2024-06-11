/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/rgravlin/noitabackup/pkg/lib"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName                   = "Noita Backup and Restore"
	ConfigMaxNumBackupsToKeep = 100
)

var cfgFile, sourcePath, destinationPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "noitabackup",
	Short: "A Noita backup and restore manager",
	Long: `A configurable Noita backup and restore manager and launcher.  Automates the tedious
task of stopping, backing up, restoring, and restarting Noita.  Includes both a GUI and command
line interface.`,
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
		go func() {
			window := new(app.Window)
			window.Option(
				app.Title(appName),
				app.MaxSize(unit.Dp(640), unit.Dp(105)),
				app.MinSize(unit.Dp(640), unit.Dp(105)),
			)
			err := lib.Run(window)
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
	rootCmd.PersistentFlags().StringVar(&sourcePath, "source-path", lib.GetDefaultSourcePath(), "Define the source Noita save00 path")
	rootCmd.PersistentFlags().StringVar(&destinationPath, "destination-path", lib.GetDefaultDestinationPath(), "Define the destination backup path")

	err := viper.BindPFlag("source-path", rootCmd.PersistentFlags().Lookup("source-path"))
	if err != nil {
		log.Printf("error binding viper flag: %v", err)
	}
	err = viper.BindPFlag("destination-path", rootCmd.PersistentFlags().Lookup("destination-path"))
	if err != nil {
		log.Printf("error binding viper flag: %v", err)
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
