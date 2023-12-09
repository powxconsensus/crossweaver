package store

import (
	"time"

	"github.com/digilabs/crossweaver/types"

	"gorm.io/gorm"
)

type HanldeMsgTxq struct {
	gorm.Model
	Id          string `gorm:"primaryKey; not null"`
	DestChainId string
	Timestamp   time.Time
	Status      types.Status
	MsgRequest  []byte
	ErrorCode   int64
	TxHash      string
}

func (dataHandler DbHandler) AddToTxq(newTxq HanldeMsgTxq) error {
	//if Id not present add it to queue
	txq, _ := dataHandler.GetTxqById(newTxq.Id)
	if txq.Id == newTxq.Id {
		dataHandler.logger.Debug("No DB action: Tx already exists with ID: ", newTxq.Id)
		return nil
	}
	// // Insert
	result := dataHandler.db.Create(&newTxq)
	dataHandler.logger.Debug("Added Tx with Id: ", newTxq.Id)
	return result.Error
}

func (dataHandler DbHandler) UpdateTxq(newTxq HanldeMsgTxq) error {
	var fundRelay HanldeMsgTxq
	//if Id not present add it to queue
	result := dataHandler.db.Model(&fundRelay).Where(&HanldeMsgTxq{Id: newTxq.Id}).
		Updates(&newTxq)
	dataHandler.logger.Debug("Updated Id: ", newTxq.Id)
	return result.Error
}

func (dataHandler DbHandler) GetTxqById(id string) (HanldeMsgTxq, error) {
	var fundRelay HanldeMsgTxq
	result := dataHandler.db.Where(&HanldeMsgTxq{Id: id}).Find(&fundRelay)
	return fundRelay, result.Error
}

func (dataHandler DbHandler) GetTxqByStatus(status types.Status) ([]HanldeMsgTxq, error) {
	fundRelayArray := []HanldeMsgTxq{}
	result := dataHandler.db.Where(&HanldeMsgTxq{Status: status}).Find(&fundRelayArray)
	return fundRelayArray, result.Error
}

func (dataHandler DbHandler) GetTxqByErrorCode(errorCode int64) ([]HanldeMsgTxq, error) {
	fundRelayArray := []HanldeMsgTxq{}
	result := dataHandler.db.Where(&HanldeMsgTxq{ErrorCode: errorCode}).Find(&fundRelayArray)
	return fundRelayArray, result.Error
}

func (dataHandler DbHandler) RemoveTxq(id string) error {
	var fundRelay HanldeMsgTxq
	result := dataHandler.db.Where(&HanldeMsgTxq{Id: id}).Find(&fundRelay)
	dataHandler.db.Unscoped().Delete(&fundRelay)
	dataHandler.logger.Debug("Removed Id: ", id)
	return result.Error

}

func (dataHandler DbHandler) DequeueTxq(limit int) error {
	var fundRelay HanldeMsgTxq
	result := dataHandler.db.Limit(limit).First(&fundRelay)
	dataHandler.db.Unscoped().Delete(&fundRelay)
	dataHandler.logger.Debug("Removed Id: ", fundRelay)
	return result.Error

}
func (dataHandler DbHandler) UpdateTxqStatus(id string, status types.Status) error {

	var fundRelay HanldeMsgTxq

	result := dataHandler.db.Model(&fundRelay).Where(&HanldeMsgTxq{Id: id}).
		Updates(&HanldeMsgTxq{Status: status, Timestamp: time.Now()})

	dataHandler.logger.Debug("Updated Id: ", id, " Status to: ", status)
	return result.Error
}

func (dataHandler DbHandler) AddTxqError(id string, errorCode int64) error {

	var fundRelay HanldeMsgTxq

	result := dataHandler.db.Model(&fundRelay).Where(&HanldeMsgTxq{Id: id}).
		Updates(&HanldeMsgTxq{Status: types.TxError, ErrorCode: errorCode, Timestamp: time.Now()})
	dataHandler.logger.Debug("Updated Id: ", id, "and add error: ", errorCode)
	return result.Error

}

func (dataHandler DbHandler) AddTxqHash(id string, hash string) error {
	var fundRelay HanldeMsgTxq
	result := dataHandler.db.Model(&fundRelay).Where(&HanldeMsgTxq{Id: id}).
		Updates(&HanldeMsgTxq{Status: types.TxDispatched, TxHash: hash, Timestamp: time.Now()})
	dataHandler.logger.Debug("Updated Id: ", id, "and added hash: ", hash)
	return result.Error
}
