package types

type Status int64

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	TxUndefined Status = iota
	TxUnprocessed
	TxReadyToExecute
	TxPicked
	TxDispatched
	TxOnHold
	TxError
	TxCompleted
)

type CUSTOM_ERR int64

const (
	FundDeposited Status = iota
	FundDepositedWithMessage
	IUSCDDeposited
)

const (
	EVM Status = iota
	DbHandler
	NEAR
)

// type Status uint64
type IChannelMsg struct {
	SrcChainId   string
	SrcChainType Status
	RelayInfo    []byte
}
