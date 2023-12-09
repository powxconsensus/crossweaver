package initializer

import (
	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/store"
	"github.com/digilabs/crossweaver/transmitter/evm/calls/gateway"
	"github.com/digilabs/crossweaver/transmitter/evm/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

func InitializeChainTransmitter(globalCfg config.GlobalCfg, chainConfig config.ChainSpecs, db *store.DbHandler, log *logrus.Entry, errChn chan<- error) (relayer.ChainTransmitter, error) {
	///////////////////////////////////////////
	///// INITIALIZE ETH CLIENT ///////////////
	///////////////////////////////////////////
	var rpcClient *rpc.Client
	var err error
	rpcClient, err = rpc.DialHTTP(chainConfig.ChainRpc)
	if err != nil {
		return relayer.ChainTransmitter{}, err
	}
	ethClient := ethclient.NewClient(rpcClient)

	//////////////////////////////////////////////
	///// INITIALIZE Voyager WRAPPER ///////////////
	//////////////////////////////////////////////
	if err != nil {
		return relayer.ChainTransmitter{}, err
	}
	digiPayAddressAddress := common.HexToAddress(chainConfig.ContractAddress)
	privateKey, err := crypto.HexToECDSA(globalCfg.EthPrivateKey)
	if err != nil {
		return relayer.ChainTransmitter{}, err
	}

	log.Info("Initializing the Chain Transmitter Instance", digiPayAddressAddress)
	gateway := gateway.NewDigiPayContract(ethClient, digiPayAddressAddress, privateKey, log)
	chainTransmitter := relayer.NewChainTransmitter(globalCfg, chainConfig, ethClient, db, gateway, log)
	return chainTransmitter, nil
}
