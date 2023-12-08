package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func main() {
	/*
	* When SIGINT or SIGTERM is caught write to the quitChannel
	 */
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	// Reading env vars
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	if err := app.Run(os.Args); err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Error("Error During Startup")
		os.Exit(1)
	}
	//////////////////////////////////////
	///////// WAIT FOR QUIT MESSAGE //////
	//////////////////////////////////////
	<-quitChannel
}

func run(cliContext *cli.Context) error {
	return nil
}
