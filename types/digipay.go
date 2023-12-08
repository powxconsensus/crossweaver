package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DigiPayLockerLocked represents a Locked event raised by the DigiPayLocker contract.
type DigiPayLockerLocked struct {
	SrcChainId string
	DstChainId string
	Nonce      *big.Int
	Tokens     []common.Address
	Amounts    []*big.Int
	Sender     common.Address
	Recipient  common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// DigiPayLockerUnLocked represents a UnLocked event raised by the DigiPayLocker contract.
type DigiPayLockerUnLocked struct {
	TxType     uint8
	SrcChainId string
	DstChainId string
	SrcNonce   *big.Int
	Nonce      *big.Int
	Tokens     []common.Address
	Amounts    []*big.Int
	Sender     common.Address
	Recipient  common.Address
	Message    []byte
	Raw        types.Log // Blockchain specific contextual infos
}
