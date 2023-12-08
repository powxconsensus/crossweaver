package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/digichain"
	digiChainInitializer "github.com/digilabs/crossweaver/listener/digichain"
	evmInitializer "github.com/digilabs/crossweaver/listener/evm/initializer"
	"github.com/digilabs/crossweaver/processor"
	evmTransmitterInitializer "github.com/digilabs/crossweaver/transmitter/evm/initializer"
	"github.com/digilabs/crossweaver/transmitter/evm/relayer"
	"github.com/digilabs/crossweaver/types"

	logging "github.com/digilabs/crossweaver/logger"
	"github.com/digilabs/crossweaver/store"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
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
	evmChainTransmitters := []relayer.ChainTransmitter{}

	//////////////////////////////////////////////////////////
	///////// INITIALIZE logger, datadog tracer ///////////////
	//////////////////////////////////////////////////////////
	logLevel, err := log.ParseLevel(cliContext.String(config.VerbosityFlag.Name))
	if err != nil {
		logLevel = log.InfoLevel
	}
	logger := logging.InitLogger(cliContext, log.Fields{}, logLevel)
	ctx, _ := context.WithCancel(context.Background())

	logger.Info("Initialse Crossweaver")
	errChn := make(chan error)

	///////////////////////////////
	///////// RESET ///////////////
	///////////////////////////////
	reset := cliContext.Bool(config.ResetFlag.Name)

	///////////////////////////////////////////
	///////// INITIALIZE CONFIG ///////////////
	///////////////////////////////////////////
	cfg, err := config.GetConfig(cliContext)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("Error Loading Config")
		panic(err)
	}
	digiChainClient := digichain.NewDigiChainClient(cfg.GlobalConfig.DigiChainRPC, cfg.GlobalConfig.EthPrivateKey)
	err = cfg.FetchSpecs(digiChainClient)
	if err != nil {
		panic(err)
	}
	logger.WithFields(log.Fields{
		"chainsCount": len(cfg.Chains),
		"chainSpecs":  len(cfg.ChainSpecs),
	}).Info("Initialize Config")

	/////////////////////////////////////////////////////
	///////// INITIALIZE DB /////////////////////////////
	/////////////////////////////////////////////////////
	logger.WithFields(log.Fields{"DbPath": cfg.GlobalConfig.DbPath}).Info("Initialize DB")
	dbHandler, err := store.InitialiseDB(cfg.GlobalConfig.DbPath, logger, reset)
	if err != nil {
		panic(err)
	}

	for _, chainConfig := range cfg.ChainSpecs {
		listenerLogger := logging.InitLogger(cliContext, log.Fields{"svc": "listener", "ChainTpe": chainConfig.ChainType, "ChainId": chainConfig.ChainId, "ChainName": chainConfig.ChainName}, log.DebugLevel)
		chainRelayerLogger := logging.InitLogger(cliContext, log.Fields{"svc": "chainTransmiiter", "ChainTpe": chainConfig.ChainType, "ChainId": chainConfig.ChainId, "ChainName": chainConfig.ChainName}, log.DebugLevel)

		/////////////////////////////////////////////////////
		///////// INITIALIZE CHAIN LISTENER /////////////////
		/////////////////////////////////////////////////////
		listenerLogger.Info("Initialize DigiPay Listener")
		switch chainConfig.ChainType {
		case types.EVM_CHAIN:
			////////////////////////////////////////////////
			///////// START EVM CHAIN LISTENER /////////////
			////////////////////////////////////////////////
			listenerLogger.Info("Starting EVM DigiPay Listener")
			chainListener, err := evmInitializer.InitializeDigiPayListener(chainConfig, digiChainClient, cfg.GlobalConfig, dbHandler, listenerLogger, errChn)
			if err != nil {
				logger.WithFields(log.Fields{"Err": err}).Error("Error while initializing digipay Listener")
				panic(err)
			}
			go func() {
				chainListener.Start(ctx)
			}()
			chainTransmiiter, _ := evmTransmitterInitializer.InitializeChainTransmitter(cfg.GlobalConfig, chainConfig, dbHandler, chainRelayerLogger, errChn)
			evmChainTransmitters = append(evmChainTransmitters, chainTransmiiter)
			go func() {
				chainTransmiiter.Start(ctx, errChn)
			}()
		default:
			panic(fmt.Errorf("type '%s' not recognized", chainConfig.ChainType))
		}
	}

	/////////////////////////////////////////////////////
	///////// INITIALIZE & START DISPATCHER /////////////
	/////////////////////////////////////////////////////
	nlogger := logging.InitLogger(cliContext, log.Fields{"svc": "dispatcher"}, log.DebugLevel)
	listener, err := digiChainInitializer.InitializeDigiPayListener(digiChainClient, cfg.GlobalConfig, dbHandler, nlogger, errChn)
	go func() {
		listener.Start()
	}()

	requestProcessor := processor.NewRequestProcessor(dbHandler, digiChainClient, logger)
	for _, chainTransmiter := range evmChainTransmitters {
		requestProcessor.AddChainRelayer(chainTransmiter)
	}
	go func() {
		requestProcessor.Start(ctx)
	}()

	return nil
}
