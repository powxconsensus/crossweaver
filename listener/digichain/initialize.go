package initializer

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/store"
	"github.com/digilabs/crossweaver/types"
	"github.com/digilabs/crossweaver/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type DigiChainListener struct {
	digiChainClient digichain.DigiChainClient
	logger          *log.Entry
	ethPrivateKey   string
	from            common.Address
	dbHandler       *store.DbHandler
	errChn          chan<- error
}

func InitializeDigiPayListener(digiChainClient digichain.DigiChainClient, globalCfg config.GlobalCfg, dbHandler *store.DbHandler, logger *log.Entry, errChn chan<- error) (DigiChainListener, error) {
	//////////////////////////////////////////////
	///// INITIALIZE EVENT PROCESSOR ///////////////
	//////////////////////////////////////////////
	return DigiChainListener{
		digiChainClient: digiChainClient,
		logger:          logger,
		ethPrivateKey:   globalCfg.EthPrivateKey,
		dbHandler:       dbHandler,
		errChn:          errChn,
		from:            common.HexToAddress(globalCfg.From),
	}, nil
}

func (c DigiChainListener) Start() {
	c.logger.Info("Crosschain Request Fetching Started")
	for {
		crossChainRequests, err := c.digiChainClient.FetchCrossChainRequest(c.from)
		if err != nil {
			c.logger.WithFields(logrus.Fields{"Error": err}).Error("Error while fetching crosschain requests")
			time.Sleep(10000) // retry after 10 sec
			continue
		}
		for idx := 0; idx < len(crossChainRequests); idx++ {
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(crossChainRequests[idx])
			id, err := utils.CreateCrosschainDBId(crossChainRequests[idx])
			if err != nil {
				c.logger.WithFields(logrus.Fields{"Error": err}).Error("Error while creating ID for DB")
				continue
			}
			txq, err := c.dbHandler.GetTxqById(id)
			if txq.Id == id {
				continue
			}
			tx := &store.HanldeMsgTxq{
				Id:          id,
				Timestamp:   time.Now(),
				DestChainId: crossChainRequests[idx].DstChainId,
				Status:      types.TxReadyToExecute,
				MsgRequest:  reqBodyBytes.Bytes(),
			}
			c.dbHandler.AddToTxq(*tx)
			c.logger.WithFields(logrus.Fields{"DestinationChainID": crossChainRequests[idx].DstChainId, "SourceChainId": crossChainRequests[idx].SrcChainId, "SrcNonce": crossChainRequests[idx].SrcNonce}).Debug("Add Crosschain Request to queue")
		}
		time.Sleep(3000)
	}
}
