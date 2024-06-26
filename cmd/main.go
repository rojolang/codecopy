package main

import (
	"os"

	"codecopy/ccopy"
	"codecopy/ui"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "--help" {
		ui.DisplayHelp()
		return
	}

	err := ccopy.Run(args)
	if err != nil {
		ui.DisplayError(err)
		os.Exit(1)
	}

	ui.DisplayHelpInfo()
}
