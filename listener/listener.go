package listener

import (
	"context"

	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/listener/types"
	"github.com/digilabs/crossweaver/store"
	log "github.com/sirupsen/logrus"
)

type IVoyagerListener interface {
	Start(ctx context.Context)
}

type DigiPayListener struct {
	ChainSpec       config.ChainSpecs
	digiChainClient digichain.DigiChainClient
	logger          *log.Entry
	eventProcessor  types.EventProcessor
	dbHandler       *store.DbHandler
}

func NewDigiPayListener(chainSpec config.ChainSpecs, digiChainClient digichain.DigiChainClient, eventProcessor types.EventProcessor, dbHandler *store.DbHandler, logger *log.Entry) DigiPayListener {
	return DigiPayListener{
		ChainSpec:       chainSpec,
		logger:          logger,
		digiChainClient: digiChainClient,
		eventProcessor:  eventProcessor,
		dbHandler:       dbHandler,
	}
}

// Start function starts the chainListener.
func (c DigiPayListener) Start(ctx context.Context) {
	// Last Proccessed Block From DB
	lastProcessedBlockFromDB, err := c.dbHandler.GetLastProcessedBlock(c.ChainSpec.ChainType, c.ChainSpec.ChainId, c.ChainSpec.ContractAddress)
	if err != nil {
		c.logger.Error("Last Processed Block not found in DB. Inserting from chain")
		c.dbHandler.UpdateLastProcessedBlock(c.ChainSpec.ChainType, c.ChainSpec.ChainId, c.ChainSpec.ContractAddress, c.ChainSpec.StartBlock, c.ChainSpec.StartEventNonce)
		lastProcessedBlockFromDB, _ = c.dbHandler.GetLastProcessedBlock(c.ChainSpec.ChainType, c.ChainSpec.ChainId, c.ChainSpec.ContractAddress)
	} else {
		c.logger.WithFields(log.Fields{"lastProcessedBlockHeight": lastProcessedBlockFromDB.BlockHeight, "lastProcessedEventNonce": lastProcessedBlockFromDB.EventNonce}).Info("Fetch LastProcessedBlock from Db")
	}
	c.logger.WithFields(log.Fields{"lastProcessedBlockHeight": lastProcessedBlockFromDB.BlockHeight, "lastProcessedEventNonce": lastProcessedBlockFromDB.EventNonce}).Info("Start Processing Events")
	c.eventProcessor.ProcessInboundEvents(lastProcessedBlockFromDB.BlockHeight, lastProcessedBlockFromDB.EventNonce)
}
