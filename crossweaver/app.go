package main

import (
	"github.com/digilabs/crossweaver/config"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var (
	Version = "0.0.1"
)

// Core Mode Flags
var CoreFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
	config.ResetFlag,
	config.LatestBlockFlag,
	config.MetricsFlag,
	config.MetricsPort,
}

var CoremodeCommand = cli.Command{
	Name:  "start",
	Usage: "Runs Crossweaver ",
	Description: "The start command is used to run crossweaver in Core Mode.\n" +
		"\tThe crossweaver directly talks to the Router chain\n" +
		"\tThe crossweaver will sign all incomming transaction\n" +
		"\tThe crossweaver can listen to various chains",
	Action: run,
	Flags:  CoreFlags,
}

func init() {
	app.Copyright = "Copyright 2023 Digi Labs"
	app.Name = "crossweaver"
	app.Usage = "crossweaver"
	app.Authors = []*cli.Author{{Name: "Digi Labs 2023"}}
	app.Version = Version
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&CoremodeCommand,
	}
}
