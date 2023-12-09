package eventprocessor

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/digilabs/crossweaver/config"
	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/store"
	"github.com/digilabs/crossweaver/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

// EventProcessor queries events from Evm chain
type EvmEventProcessor struct {
	chainSpec              config.ChainSpecs
	digiChainClient        digichain.DigiChainClient
	logger                 *log.Entry
	ethclient              *ethclient.Client
	ethPrivateKey          string
	dbHandler              *store.DbHandler
	middlewareAddress      string
	forwarderRouterAddress string
	errChn                 chan<- error
	from                   common.Address
}

func NewEvmEventProcessor(chainSpec config.ChainSpecs, digiChainClient digichain.DigiChainClient, ethClient *ethclient.Client, from string, ethPrivateKey string, dbHandler *store.DbHandler, logger *log.Entry, errChn chan<- error) EvmEventProcessor {
	eventProcessor := EvmEventProcessor{
		chainSpec:       chainSpec,
		logger:          logger,
		ethclient:       ethClient,
		ethPrivateKey:   ethPrivateKey,
		digiChainClient: digiChainClient,
		dbHandler:       dbHandler,
		errChn:          errChn,
		from:            common.HexToAddress(from),
	}
	return eventProcessor
}

func (e EvmEventProcessor) ProcessInboundEvents(lastQueriedBlock uint64, lastProcessedEventNonce uint64) error {
	e.logger.Info("Block Fetching started")
	for {
		fmt.Println(lastQueriedBlock)
		////////////////////////////////////////////////////////////
		/////// 1. FETCH EVENTS FROM SOURCE CHAIN 	////////////////
		/////// 2. TRANSFORM EVENTS TO ROUTERCHAIN SDK.MSG /////////
		////////////////////////////////////////////////////////////
		startBlock, endBlock, err := e.CalculateStartAndEndBlocks(lastQueriedBlock)
		if err != nil {
			// e.logger.WithFields(log.Fields{"lastQueriedBlock": lastQueriedBlock}).Error("No new blocks to process")
			// Sleep as there are no new confirmed blocks
			blockTime, err := time.ParseDuration(e.chainSpec.BlockTime)
			if err != nil {
				e.logger.WithFields(log.Fields{"BlockTime": e.chainSpec.BlockTime}).Error("Error parsing blocktime")
				panic(err)
			}
			time.Sleep(blockTime)
			continue
		}
		e.logger.WithFields(log.Fields{"StartBlock": startBlock, "EndBlock": endBlock}).Debug("Querying for Events")
		newProcessedEventnonce, err := e.FetchAndProcessMsg(lastProcessedEventNonce, startBlock, endBlock)
		if err != nil {
			e.logger.WithFields(log.Fields{"lastProcessedEventNonce": lastProcessedEventNonce, "startBlock": startBlock, "endBlock": endBlock}).Error("Error in FetchAndProcessMsg")
			return err
		}
		e.logger.WithFields(log.Fields{"StartBlock": startBlock, "EndBlock": endBlock, "newProcessedEventnonce": newProcessedEventnonce}).Debug("Querying for Events")

		///////////////////////////////////////////////////////
		///// 4. UPDATE PROCESSED BLOCK AND EVENT NONCE ///////
		///////////////////////////////////////////////////////
		lastQueriedBlock = endBlock
		lastProcessedEventNonce = newProcessedEventnonce
		// Add them to DB to avoid repeating same blocks when orchestrator restarts
		e.dbHandler.UpdateLastProcessedBlock(e.chainSpec.ChainType, e.chainSpec.ChainId, e.chainSpec.ContractAddress, lastQueriedBlock, lastProcessedEventNonce)
	}
}

func (e EvmEventProcessor) FetchAndProcessMsg(lastProcessedEventNonce uint64, startBlock uint64, endBlock uint64) (uint64, error) {
	retryCount := 3
	var lockedEvents []*types.DigiPayLockerLocked
	var unlockedEvents []*types.DigiPayLockerUnLocked
	var voyagerError error

	///////////////////////////////////////////////////
	///// 1. Retry incase of digipay query errors  ////
	///////////////////////////////////////////////////
	for i := 0; i < retryCount; i++ {
		lockedEvents, unlockedEvents, voyagerError = e.QueryDigiPayEvents(startBlock, endBlock)
		if voyagerError == nil {
			break
		}
		e.logger.WithFields(log.Fields{"startBlock": startBlock, "endBlock": endBlock, "ChainId": e.chainSpec.ChainId, "ChainName": e.chainSpec.ChainName, "voyagerError": voyagerError}).Debug("Retrying voyager events query due to voyager Error")
		time.Sleep(30 * time.Second)
	}
	if voyagerError != nil {
		return 0, voyagerError
	}
	e.logger.WithFields(log.Fields{"startBlock": startBlock, "endBlock": endBlock, "lockedEvents": lockedEvents, "unlockedEvents": unlockedEvents}).Debug("Query Events")
	fmt.Println("lastProcessedEventNonce by chain", lastProcessedEventNonce, "chainName", e.chainSpec.ChainName)
	lastProcessedEventNonce, err := e.SortAndTransformInboundEventsByEventNonce(lastProcessedEventNonce, lockedEvents, unlockedEvents)
	if err != nil {
		return 0, err
	}
	return lastProcessedEventNonce, nil
}

// We use merge sort to sort individual event arrays. Also we transform source events to Routerchain Messages.
// The output of this function contains array of cosmos messages which are sorted by event nonce.
func (e EvmEventProcessor) SortAndTransformInboundEventsByEventNonce(
	lastProcessedEventNonce uint64,
	lockedEvents []*types.DigiPayLockerLocked,
	unlockedEvents []*types.DigiPayLockerUnLocked,
) (uint64, error) {
	totalClaimEvents := len(lockedEvents) + len(unlockedEvents)
	var count, i, j int

	// Sort slices using the custom sorting function
	sort.Slice(lockedEvents, func(i, j int) bool {
		return compareEventNonce(lockedEvents[i].Nonce, lockedEvents[j].Nonce)
	})
	sort.Slice(unlockedEvents, func(i, j int) bool {
		return compareEventNonce(unlockedEvents[i].Nonce, unlockedEvents[j].Nonce)
	})

	// Sort events sequentially starting with eventNonce = lastClaimEvent + 1.
	for count < totalClaimEvents {
		var msg digichain.RawProposal
		var err error
		if i < len(lockedEvents) {
			if lockedEvents[i].Nonce.Uint64() <= lastProcessedEventNonce {
				// It's already processed. skip processing again.
				i++
			} else if lockedEvents[i].Nonce.Uint64() == lastProcessedEventNonce+1 {
				////////////////////////////////////////////////
				///// 1. Transform RequestToRouterEvent ////////
				////////////////////////////////////////////////
				msg, err = e.TransformLockedEvent(lockedEvents[i])
				if err != nil {
					return 0, err
				}
				lastProcessedEventNonce = lockedEvents[i].Nonce.Uint64()
				// Can Be Optimized
				ires, err := e.digiChainClient.IsCrosschainRequestBroadcasted(e.from, lockedEvents[i].SrcChainId, lockedEvents[i].Nonce)
				i++
				if err != nil {
					e.logger.WithFields(log.Fields{"error": err}).Error("Error While Fetching IsCrosschainRequestBroadcasted")
					continue
				}
				if ires {
					e.logger.WithFields(log.Fields{"error": err}).Error("Tx Already Broadcasted")
					continue
				}
				// directly signing and broadcasting
				res, error := e.digiChainClient.SignRawTxAndBroadCast(msg)
				if error != nil {
					e.logger.WithFields(log.Fields{"error": error}).Error("LockedEvent: Error While Broadcasting tx")
					continue
				}
				e.logger.WithFields(log.Fields{"tx_hash": res.Data.TxHash}).Debug("LockedEvent: Broadcasted Tx")
			}
		}

		if j < len(unlockedEvents) {
			if unlockedEvents[j].Nonce.Uint64() <= lastProcessedEventNonce {
				// It's already processed. skip processing again.
				j++
			} else if unlockedEvents[j].Nonce.Uint64() == lastProcessedEventNonce+1 {
				////////////////////////////////////////////////
				///// 2. Transform ValsetUpdatedEvent //////////
				////////////////////////////////////////////////
				msg, err = e.TransformUnLockedEvent(unlockedEvents[j])
				if err != nil {
					return 0, err
				}
				lastProcessedEventNonce = unlockedEvents[j].Nonce.Uint64()
				// Can Be Optimized
				ires, err := e.digiChainClient.IsCrosschainRequestBroadcasted(e.from, unlockedEvents[j].DstChainId, unlockedEvents[j].Nonce)
				if err != nil {
					e.logger.WithFields(log.Fields{"error": err}).Error("Error While Fetching IsCrosschainRequestBroadcasted")
					continue
				}
				if ires {
					e.logger.WithFields(log.Fields{"error": err}).Error("Tx Already Broadcasted")
					continue
				}
				j++
				res, error := e.digiChainClient.SignRawTxAndBroadCast(msg)
				if error != nil {
					e.logger.WithFields(log.Fields{"error": error}).Error("UnlockedEvent: Error While Broadcasting tx")
					continue
				}
				e.logger.WithFields(log.Fields{"tx_hash": res.Data.TxHash}).Debug("UnlockedEvent: Broadcasted Tx")
			}
		}
		count = count + 1
	}
	return lastProcessedEventNonce, nil
}

func (e EvmEventProcessor) CalculateStartAndEndBlocks(lastQueriedBlock uint64) (uint64, uint64, error) {
	////////////////////////////////////////////////////////////////////////////////
	/////// 1. StartBlock = lastQueriedBlock + 1 	////////////////////////////////
	/////// 2. latestConfirmedBlock = latestBlock - ConfirmationsRequired //////////
	/////// 3. No Blocks to query if latestBlock < (StartBlock + ConfirmationsRequired)   /////////
	/////// 4. EndBlock = StartBlock + BlocksToSearch if (StartBlock + BlocksToSearch) < latestConfirmedBlock  /////////
	/////// 5. EndBlock = StartBlock + latestConfirmedBlock if (StartBlock + BlocksToSearch) > latestConfirmedBlock  /////////
	////////////////////////////////////////////////////////////////////////////////
	latestHeader, err := e.ethclient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		e.logger.WithFields(log.Fields{"error": err}).Error("Failed to get latest block")
		e.errChn <- err
		return 0, 0, err
	}
	latestBlock := latestHeader.Number.Uint64()
	startBlock := lastQueriedBlock + 1
	confirmationsRequired := e.chainSpec.ConfirmationsRequired
	blocksToSearch := e.chainSpec.BlocksToSearch
	latestConfirmedBlock := latestBlock - confirmationsRequired
	endBlock := startBlock

	// There are no new confirmed blocks to process
	if latestBlock < (startBlock + confirmationsRequired) {
		return startBlock, endBlock, errors.New("no new confirmed blocks to process")
	}

	// calculate endblock
	if (startBlock + blocksToSearch) < latestConfirmedBlock {
		endBlock = startBlock + blocksToSearch
	} else {
		endBlock = latestConfirmedBlock
	}
	e.logger.WithFields(log.Fields{"startBlock": startBlock, "endBlock": endBlock, "confirmationsRequired": confirmationsRequired, "blocksToSearch": blocksToSearch, "latestBlock": latestBlock}).Debug("Calculate Start and endblock for voyager events query")
	return startBlock, endBlock, nil
}

func compareEventNonce(num1, num2 *big.Int) bool {
	return num1.Cmp(num2) == -1
}
