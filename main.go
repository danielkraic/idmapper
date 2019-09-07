package main

import (
	"log"

	"github.com/danielkraic/idmapper/app"
	"github.com/spf13/pflag"
)

var (
	// Version will be set during build
	Version = ""
	// Commit will be set during build
	Commit = ""
	// Build will be set during build
	Build = ""
)

func main() {
	pflag.StringP("addr", "a", "0.0.0.0:80", "HTTP service address.")
	configFile := pflag.StringP("config", "c", "", "path to config file")
	pflag.Parse()

	app, err := app.NewApp(Version, Commit, Build, *configFile)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
