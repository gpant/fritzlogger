package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {

	app := cli.App("fritzlogger", "Parse TR064 System logs and merge them")
	app.Spec = "[-d] [-e] DIR [LOGFILE]"
	app.Version("v version", fmt.Sprintf("fritzlogger Version %v\nCopyright (c) 2022 George Pantazis\nGNU GPL 2.0", "1.1"))
	app.LongDesc = "fritzlogger will merge & sort fritz system log dumps including TR064 messages"

	var (
		// version = app.BoolOpt("v version", false, "Application Version")
		remove           = app.BoolOpt("d delete", false, "Delete Parsed Logs")
		dir              = app.StringArg("DIR", "", "Directory Containing Logs")
		mergelogfilename = app.StringArg("LOGFILE", "fritz.logs", "The Merged Log File")
		emulateonly      = app.BoolOpt("e emulate", false, "Only output to console, dont copy/move/delete logs")
		_                = emulateonly
	)

	app.Action = func() {
		var localpath string

		if *dir == "" {
			localpath, _ = os.Getwd()
		} else {
			localpath = *dir
		}

		if _, err := os.Stat(localpath); !os.IsNotExist(err) {
			parselogs(*remove, localpath, *mergelogfilename)
		} else {
			fmt.Printf("%s directory does not exist\n", localpath)
		}
	}

	app.Run(os.Args)
}
