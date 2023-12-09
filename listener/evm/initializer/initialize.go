package initializer

import (
	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/listener"
	"github.com/digilabs/crossweaver/listener/evm/eventprocessor"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/digilabs/crossweaver/store"
	log "github.com/sirupsen/logrus"
)

func InitializeDigiPayListener(chainSpec config.ChainSpecs, digiChainClient digichain.DigiChainClient, globalCfg config.GlobalCfg, dbHandler *store.DbHandler, logger *log.Entry, errChn chan<- error) (listener.DigiPayListener, error) {
	///////////////////////////////////////////
	///// INITIALIZE ETH CLIENT ///////////////
	///////////////////////////////////////////
	var rpcClient *rpc.Client
	var err error
	rpcClient, err = rpc.DialHTTP(chainSpec.ChainRpc)
	if err != nil {
		errChn <- err
		panic(err)
	}
	ethClient := ethclient.NewClient(rpcClient)

	//////////////////////////////////////////////
	///// INITIALIZE EVENT PROCESSOR ///////////////
	//////////////////////////////////////////////
	eventProcessor := eventprocessor.NewEvmEventProcessor(chainSpec, digiChainClient, ethClient, globalCfg.From, globalCfg.EthPrivateKey, dbHandler, logger, errChn)

	///////////////////////////////////////////////
	///// INITIALIZE CHAIN LISTENER ///////////////
	///////////////////////////////////////////////
	chainListener := listener.NewDigiPayListener(chainSpec, digiChainClient, eventProcessor, dbHandler, logger)
	return chainListener, nil
}
