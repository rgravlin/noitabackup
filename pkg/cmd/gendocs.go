/*
Package cmd
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log"
)

// gendocsCmd represents the gendocs command
var gendocsCmd = &cobra.Command{
	Use:   "gendocs",
	Short: "Generate command line documentation",
	Long:  `Generates markdown documentation and places in the ./docs folder`,
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gendocsCmd)
}
