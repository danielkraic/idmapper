package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	checkConfig := pflag.Bool("config-check", false, "check configuration")
	printConfig := pflag.BoolP("print-config", "p", false, "print configuration")
	pflag.Parse()

	app, err := app.NewApp(Version, Commit, Build, *configFile)
	if err != nil {
		log.Fatal(err)
	}

	if *printConfig {
		app.PrintConfiguration()
	}

	app.SetupRedis()
	err = app.SetupPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	err = app.SetupIDMappers()
	if err != nil {
		log.Fatal(err)
	}

	if *checkConfig {
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	app.Run(signalChan)
}
