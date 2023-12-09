package relayer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	labi "github.com/digilabs/crossweaver/abi"
	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/transmitter/evm/calls/methods"
	"github.com/digilabs/crossweaver/types"
	"github.com/digilabs/crossweaver/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

func (ct *ChainTransmitter) HandleRequest(m *digichain.CrossChainRequest) error {
	id, err := utils.CreateCrosschainDBId(*m)
	if err != nil {
		ct.log.WithFields(logrus.Fields{"Error": err}).Error("Error while creating ID for DB")
		return err
	}
	// check if this msg is already executed or not?
	var tnonce big.Int
	src_nonce, _ := tnonce.SetString(m.SrcNonce, 10)
	excutedArgs, err := methods.Executed.GetAbiPackBytes(labi.DIGIPAY_ABI, src_nonce) // src_nonce is nonce on digichain
	res, err := ct.gateway.CallContract(excutedArgs)
	if err != nil {
		ct.log.WithFields(logrus.Fields{"Error": err}).Error("Error fetching if tx already executed or not")
		return err
	}
	executed := new(big.Int).SetBytes(res).Sign() != 0
	if executed {
		ct.log.WithFields(logrus.Fields{"Error": err, "Id": id}).Error("Already Executed")
		return err
	}
	payload, err := hex.DecodeString(utils.Remove0xPrefix(m.Payload))
	if err != nil {
		ct.log.WithFields(logrus.Fields{"Error": err}).Error("Error while decoding payload")
		return err
	}
	sigs := make([]string, len(m.Sigs))
	for idx := 0; idx < len(m.Sigs); idx++ {
		sigs[idx], _ = m.Sigs[idx].ToHex()
	}
	argsData, err := methods.HandleRequest.GetAbiPackBytes(labi.DIGIPAY_ABI, m.SrcChainId, m.DstChainId, src_nonce, payload, sigs)
	if err != nil {
		ct.log.Fatalf("Failed to pack input data: %v", err)
	}

	// Get the nonce for the sender's address
	nonce, err := ct.client.PendingNonceAt(context.Background(), ct.gateway.From)
	if err != nil {
		log.Fatalf("Failed to retrieve sender's nonce: %v", err)
		return nil
	}
	gasLimit, err := ct.gateway.SimulateTransaction(string(methods.HandleRequest), ct.gateway.DigiPayAddress, argsData)
	if err != nil {
		ct.log.WithFields(logrus.Fields{"RequestId": id, "Error": err}).Error("HandleRequest: Simulate Failing")
		err = ct.dbHandler.UpdateTxqStatus(id, types.TxError)
		if err != nil {
			ct.log.WithFields(logrus.Fields{"RequestId": id, "Error": err}).Error("HandleRequest: Error While updating status to Completed")
			return err
		}
		return fmt.Errorf("HandleRequest: Simulate Failing")
	}

	// retrieving suggested gas fees and gas price.
	tipCap, err := ct.client.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	feeCap, err := ct.client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	chainID, err := ct.client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	// Create the transaction
	tx := etypes.NewTx(&etypes.DynamicFeeTx{
		ChainID:   chainID,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Nonce:     nonce,
		To:        &ct.gateway.DigiPayAddress,
		Value:     big.NewInt(0),
		Gas:       gasLimit,
		Data:      []byte(argsData),
	})
	// Sign the transaction
	signedTx, err := etypes.SignTx(tx, etypes.LatestSignerForChainID(chainID), ct.gateway.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
		return err
	}
	err = ct.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		ct.log.WithFields(logrus.Fields{"RequestId": id, "Error": err}).Error("HandleRequest: Error While sending the Tx")
		dbErr := ct.dbHandler.AddTxqError(id, types.ContractError().ErrorID)
		if dbErr != nil {
			ct.log.WithFields(logrus.Fields{"RequestId": id, "Error": err}).Error("Error: While adding tx error to DB")
			return dbErr
		}
		return err
	}
	ct.log.WithFields(logrus.Fields{"hash": signedTx.Hash()}).Debug("HandleRequest: Transaction successfully submitted")
	err = ct.dbHandler.AddTxqHash(id, signedTx.Hash().Hex())
	if err != nil {
		ct.log.WithFields(logrus.Fields{"RequestId": id, "Error": err}).Error("Error: While adding tx hash to DB")
		return err
	}
	receipt, err := bind.WaitMined(context.Background(), ct.client, signedTx)
	if err != nil {
		ct.log.WithFields(logrus.Fields{"RequestId": id, "tsHash": signedTx.Hash().Hex(), "Error": err}).Error("HandleRequest: Error While waiting for Tx Receipt")
		return err
	}
	ct.log.WithFields(logrus.Fields{"hash": receipt.TxHash}).Info("HandleRequest: Transaction Completed")
	ct.log.WithFields(logrus.Fields{"CumulativeGasUsed": receipt.CumulativeGasUsed, "blockNumber": receipt.BlockNumber}).Info("HandleRequest: Receipt Details")
	err = ct.dbHandler.UpdateTxqStatus(id, types.TxCompleted)
	if err != nil {
		ct.log.WithFields(logrus.Fields{"RequestId": id, "Error": err}).Error("HandleRequest: Error While updating status to Completed")
		return err
	}
	return nil
}
