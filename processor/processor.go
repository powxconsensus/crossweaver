package processor

import (
	"context"
	"time"

	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/store"
	"github.com/digilabs/crossweaver/transmitter"
	"github.com/digilabs/crossweaver/types"
	"github.com/sirupsen/logrus"
)

func NewRequestProcessor(dbHandler *store.DbHandler, digiChainClient digichain.DigiChainClient, log *logrus.Entry) *RequestProcessor {
	chainRegistry := make(map[string]transmitter.IChainTransmitter)
	requestProcessor := &RequestProcessor{
		registry:        chainRegistry,
		dbHandler:       dbHandler,
		digiChainClient: digiChainClient,
		log:             log,
	}
	return requestProcessor
}

type RequestProcessor struct {
	registry        map[string]transmitter.IChainTransmitter
	dbHandler       *store.DbHandler
	digiChainClient digichain.DigiChainClient
	log             *logrus.Entry
}

func (requestProcessor *RequestProcessor) AddChainRelayer(relayChain transmitter.IChainTransmitter) {
	destinationID := relayChain.DestinationChainId()
	requestProcessor.registry[destinationID] = relayChain
}

func (requestProcessor *RequestProcessor) Start(ctx context.Context) {
	requestProcessor.log.Debug("Start Processor Service")
	// Listen to request requests from sql and process
	go func() {
		for {
			msgTxqArray, err := requestProcessor.dbHandler.GetTxqByStatus(types.TxReadyToExecute)
			if err != nil {
				panic(err)
			}
			for _, txq := range msgTxqArray {
				err = requestProcessor.dbHandler.UpdateTxqStatus(txq.Id, types.TxPicked)
				if err != nil {
					panic(err)
				}
				if requestProcessor.registry[txq.DestChainId] != nil {
					// Receive msgs from Listener and send it to executor
					requestProcessor.registry[txq.DestChainId].AddHandleRequestToMsgChannel(txq.MsgRequest)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
