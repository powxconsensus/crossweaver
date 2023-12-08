package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gogo/protobuf/proto"
)

type ChainType int32

const (
	NONE_CHAIN ChainType = 0
	DIGI_CHAIN ChainType = 1
	EVM_CHAIN  ChainType = 2
)

var ChainType_name = map[int32]string{
	0: "NONE_CHAIN",
	1: "DIGI_CHAIN",
	2: "EVM_CHAIN",
}

func (x ChainType) String() string {
	return proto.EnumName(ChainType_name, int32(x))
}

// DigiPayLockerUnLocked represents a UnLocked event raised by the DigiPayLocker contract.
type CrossChainRequest struct {
	SrcChainId string
	SrcNonce   *big.Int
	Nonce      *big.Int
	Tokens     []common.Address
	Amounts    []*big.Int
	Sender     common.Address
	Recipient  common.Address
	Message    []byte
	Raw        types.Log // Blockchain specific contextual infos
}
