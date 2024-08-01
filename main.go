//go:build windows

/*
Copyright © 2024 Ryan Gravlin ryan.gravlin@gmail.com
*/
package main

import (
	"github.com/rgravlin/noitabackup/pkg/cmd"
	"github.com/spf13/cobra"
)

func main() {
	// disable the mousetrap
	cobra.MousetrapHelpText = ""
	cmd.Execute()
}
