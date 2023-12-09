package eventprocessor

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/digilabs/crossweaver/abi"
	"github.com/digilabs/crossweaver/types"
	"github.com/ethereum/go-ethereum"
	eabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
)

func (e *EvmEventProcessor) QueryDigiPayEvents(startBlock uint64, endBlock uint64) (
	lockedEvents []*types.DigiPayLockerLocked,
	unLockedEvents []*types.DigiPayLockerUnLocked,
	err error) {
	eventNames := []string{
		"Locked",
		"UnLocked",
	}
	var topics []common.Hash
	contract, err := eabi.JSON(strings.NewReader(abi.DIGIPAY_ABI))
	if err != nil {
		e.logger.Errorf("Failed to parse contract ABI: %v", err)
	}

	for _, eventName := range eventNames {
		event, exists := contract.Events[eventName]
		if !exists {
			e.logger.Errorf("event %s not found in ABI", eventName)
		}
		topics = append(topics, common.HexToHash(event.ID.Hex()))
	}
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(startBlock),
		ToBlock:   new(big.Int).SetUint64(endBlock), // latest
		Addresses: []common.Address{
			common.HexToAddress(e.chainSpec.ContractAddress),
		},
		Topics: [][]common.Hash{
			topics,
		},
	}

	eventLogs, err := e.ethclient.FilterLogs(context.Background(), query)
	if err != nil {
		e.logger.Errorf("Failed to retrieve logs: %v", err)
	}
	for _, eventLog := range eventLogs {
		event, err := contract.EventByID(common.HexToHash(eventLog.Topics[0].Hex()))
		if err != nil {
			e.logger.Errorf("Unknown event topic: %s", eventLog.Topics[0].Hex())
			continue
		}
		switch event.Name {
		case "Locked":
			Event := new(types.DigiPayLockerLocked)
			err := UnpackEventLog(Event, "Locked", eventLog)
			if err != nil {
				e.logger.Errorf("Failed to unpack locked event")
				continue
			}
			Event.Raw = eventLog
			lockedEvents = append(lockedEvents, Event)
		case "UnLocked":
			Event := new(types.DigiPayLockerUnLocked)
			fmt.Println(eventLog)
			err := UnpackEventLog(Event, "UnLocked", eventLog)
			if err != nil {
				e.logger.Errorf("Failed to unpack UnLocked event")
				continue
			}
			Event.Raw = eventLog
			unLockedEvents = append(unLockedEvents, Event)
		}
	}
	return
}

func UnpackEventLog(Event interface{}, eventName string, eventLog etypes.Log) error {
	contract, err := eabi.JSON(strings.NewReader(abi.DIGIPAY_ABI))
	if err != nil {
		return err
	}

	if err := contract.UnpackIntoInterface(Event, eventName, eventLog.Data); err != nil {
		return err
	}
	var indexed eabi.Arguments
	for _, arg := range contract.Events[eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return eabi.ParseTopics(Event, indexed, eventLog.Topics[1:])
}
