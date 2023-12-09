package eventprocessor

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/types"
	"github.com/ethereum/go-ethereum/common"

	log "github.com/sirupsen/logrus"
)

func (e *EvmEventProcessor) TransformLockedEvent(lockedEvents *types.DigiPayLockerLocked) (digichain.RawProposal, error) {
	destChainId := string(lockedEvents.DstChainId)
	destChainId = strings.TrimRight(destChainId, "\u0000") // removing null charater from right
	destChainId = strings.Trim(destChainId, "\x00")

	rmsg := digichain.NewCrossChainRequestMsg(
		lockedEvents.SrcChainId,
		destChainId,
		common.HexToAddress(e.chainSpec.ContractAddress),
		lockedEvents.Recipient,
		lockedEvents.Sender,
		lockedEvents.Tokens,
		lockedEvents.Amounts,
		lockedEvents.Nonce,
		lockedEvents.Raw.BlockNumber,
		lockedEvents.Raw.TxHash.String(),
	)
	fmt.Println("rmsg: ", rmsg)
	data, err := rmsg.GetMsgPacket()
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Fatalf("While getting msg packet")
	}

	txTypeMsg := digichain.NewCCTxTypeMsg(
		uint8(0),
		lockedEvents.SrcChainId,
		lockedEvents.Nonce,
		lockedEvents.DstChainId,
		big.NewInt(0),
	)
	txTypedata, err := txTypeMsg.GetMsgPacket()
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Fatalf("While getting TxType msg packet")
	}

	msg, err := e.digiChainClient.GetRawProposal(
		fmt.Sprintf("CrossChainRequest(%s)", digichain.HexString(hex.EncodeToString(txTypedata))),
		data,
	)
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Fatalf("error while generating raw proposal")
	}
	/////////////////////////////////////////////////////////////
	/// Transform LockedEvent to  NewLockedMsg ///////
	/////////////////////////////////////////////////////////////
	e.logger.WithFields(log.Fields{"msg": msg}).Info("Found LockedEvent. Transformed to CrossChainRequestMsg")
	return msg, nil
}

func (e *EvmEventProcessor) TransformUnLockedEvent(unLockedEvents *types.DigiPayLockerUnLocked) (digichain.RawProposal, error) {
	destChainId := string(unLockedEvents.DstChainId)
	destChainId = strings.TrimRight(destChainId, "\u0000") // removing null charater from right
	destChainId = strings.Trim(destChainId, "\x00")

	rmsg := digichain.NewCrossChainReplyMsg(
		unLockedEvents.SrcChainId,
		destChainId,
		common.HexToAddress(e.chainSpec.ContractAddress),
		unLockedEvents.Recipient,
		unLockedEvents.Sender,
		unLockedEvents.Tokens,
		unLockedEvents.Amounts,
		unLockedEvents.SrcNonce,
		unLockedEvents.Nonce,
		unLockedEvents.Raw.BlockNumber,
		unLockedEvents.Raw.TxHash.String(),
	)
	data, err := rmsg.GetMsgPacket()
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Fatalf("While getting msg packet")
	}

	txTypeMsg := digichain.NewCCTxTypeMsg(
		unLockedEvents.TxType,
		unLockedEvents.SrcChainId,
		unLockedEvents.SrcNonce,
		unLockedEvents.DstChainId,
		unLockedEvents.Nonce,
	)
	txTypedata, err := txTypeMsg.GetMsgPacket()
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Fatalf("While getting TxType msg packet")
	}
	msg, err := e.digiChainClient.GetRawProposal(
		fmt.Sprintf("CrossChainRequest(%s)", digichain.HexString(hex.EncodeToString(txTypedata))),
		data,
	)
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Fatalf("error while generating raw proposal")
	}
	/////////////////////////////////////////////////////////////
	/// Transform LockedEvent to  NewUnLockedMsg ///////
	/////////////////////////////////////////////////////////////
	e.logger.WithFields(log.Fields{"msg": msg}).Info("Found UnLockedEvent. Transformed to CrossChainReplyMsg")
	return msg, nil
}
