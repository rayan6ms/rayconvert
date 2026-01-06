package main

import (
	"os"

	"github.com/rayan6ms/rayconvert/internal/app"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	os.Exit(app.Run(os.Args[1:], app.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}))
}
