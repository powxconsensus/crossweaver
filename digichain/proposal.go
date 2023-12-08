package digichain

import (
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (dc *DigiChainClient) GetRawProposal(
	proposalType string,
	data []byte,
) (RawProposal, error) {
	msg := NewRawProposal(
		"0x12345",
		dc.GetChainId(),
		proposalType,
		common.HexToAddress(dc.FromAddress()),
		uint64(time.Now().UnixMilli()),
		HexString(hex.EncodeToString(data)),
	)
	return msg, nil
}
