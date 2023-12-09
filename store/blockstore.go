package store

import (
	"errors"

	"github.com/digilabs/crossweaver/types"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProcessedBlock struct {
	gorm.Model
	Id          string `gorm:"primaryKey; not null"`
	ChainType   string `gorm:"primaryKey"`
	ChainId     string `gorm:"primaryKey"`
	Contract    string `gorm:"primaryKey"`
	BlockHeight uint64
	EventNonce  uint64
}

func (dataHandler DbHandler) GetLastProcessedBlock(chainType types.ChainType, chainID string, contract string) (ProcessedBlock, error) {
	var processedBlock ProcessedBlock
	result := dataHandler.db.First(&processedBlock, "chain_type = ? AND chain_id = ? AND contract = ?", chainType.String(), chainID, contract)
	return processedBlock, result.Error
}

func (dataHandler DbHandler) UpdateLastProcessedBlock(chainType types.ChainType, chainID string, contract string, newProcessedBlockHeight uint64, newProcessedEventnonce uint64) {
	// find ProcessedBlock
	processedBlock, err := dataHandler.GetLastProcessedBlock(chainType, chainID, contract)
	// Insert
	if errors.Is(err, gorm.ErrRecordNotFound) {
		processedBlock := ProcessedBlock{
			ChainType:   chainType.String(),
			ChainId:     chainID,
			Contract:    contract,
			BlockHeight: newProcessedBlockHeight,
			EventNonce:  newProcessedEventnonce,
		}
		dataHandler.logger.WithFields(log.Fields{"processedBlock": processedBlock}).Error("Last Proccessed Block Not Found in DB. Inserting new")
		dataHandler.db.Create(&processedBlock)
		return
	}
	// Update db
	dataHandler.db.Model(&processedBlock).Where(&ProcessedBlock{ChainType: chainType.String(), ChainId: chainID, Contract: contract}).Updates(ProcessedBlock{BlockHeight: newProcessedBlockHeight, EventNonce: newProcessedEventnonce}) // non-zero fields
}
