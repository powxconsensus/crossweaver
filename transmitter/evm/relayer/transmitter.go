package relayer

import (
	"context"
	"encoding/json"

	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/store"
	"github.com/digilabs/crossweaver/transmitter/evm/calls/gateway"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/sirupsen/logrus"
)

type ChainTransmitter struct {
	iHandleRequestChannel chan []byte
	globalCfg             config.GlobalCfg
	chainConfig           config.ChainSpecs
	log                   *logrus.Entry
	dbHandler             *store.DbHandler
	gateway               *gateway.DigiPayContract
	client                *ethclient.Client
}

func NewChainTransmitter(globalCfg config.GlobalCfg, chainConfig config.ChainSpecs, client *ethclient.Client, dbHandler *store.DbHandler, gateway *gateway.DigiPayContract, log *logrus.Entry) ChainTransmitter {
	iHandleRequestChannel := make(chan []byte)
	return ChainTransmitter{
		iHandleRequestChannel: iHandleRequestChannel,
		globalCfg:             globalCfg,
		chainConfig:           chainConfig,
		log:                   log,
		dbHandler:             dbHandler,
		gateway:               gateway,
		client:                client,
	}
}

func (c ChainTransmitter) DestinationChainId() string {
	return c.chainConfig.ChainId
}

func (c ChainTransmitter) AddHandleRequestToMsgChannel(msg []byte) {
	c.iHandleRequestChannel <- msg
}

// Start function starts the relayer. Relayer routine is starting all the chains
// and passing them with a channel that accepts unified cross chain message format
func (c ChainTransmitter) Start(ctx context.Context, sysErr chan error) {
	// c.log.Info("Start Chain Relayer")
	for {
		select {
		case msg := <-c.iHandleRequestChannel:
			var handleMsg *digichain.CrossChainRequest
			err := json.Unmarshal(msg, &handleMsg)
			if err != nil {
				c.log.WithFields(logrus.Fields{"err": err}).Error("Error: While Transforming Event")
				continue
			}
			err = c.HandleRequest(handleMsg)
			if err != nil {
				c.log.WithFields(logrus.Fields{"err": err}).Error("Error: While Relaying Funds ")
				continue
			}
			continue
		}
	}
}
