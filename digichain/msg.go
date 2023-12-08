package digichain

import (
	"math/big"

	"github.com/digilabs/crossweaver/abi"
	"github.com/ethereum/go-ethereum/common"
)

type CrossChainRequestMsg struct {
	src_chain_id string
	dst_chain_id string
	src_contract common.Address
	recipient    common.Address
	depositor    common.Address
	tokens       []common.Address
	amounts      []*big.Int
	src_nonce    *big.Int
	block_number uint64
	src_tx_hash  string
}

func NewCrossChainRequestMsg(
	src_chain_id string,
	dst_chain_id string,
	src_contract common.Address,
	recipient common.Address,
	depositor common.Address,
	tokens []common.Address,
	amounts []*big.Int,
	src_nonce *big.Int,
	block_number uint64,
	src_tx_hash string,
) CrossChainRequestMsg {
	crossChainRequestMsg := CrossChainRequestMsg{
		src_chain_id: src_chain_id,
		dst_chain_id: dst_chain_id,
		src_contract: src_contract,
		recipient:    recipient,
		depositor:    depositor,
		tokens:       tokens,
		amounts:      amounts,
		src_nonce:    src_nonce,
		block_number: block_number,
		src_tx_hash:  src_tx_hash,
	}
	return crossChainRequestMsg
}

func (c *CrossChainRequestMsg) GetMsgPacket() ([]byte, error) {
	data, err := abi.CROSS_CHAIN_REQUEST_INTERFACE.Pack(
		c.src_chain_id,
		c.dst_chain_id,
		c.src_contract,
		c.recipient,
		c.depositor,
		c.tokens,
		c.amounts,
		c.src_nonce,
		c.block_number,
		c.src_tx_hash,
	)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// Reply Of Crosschain Request From Other Chain
type CrossChainReplyMsg struct {
	src_chain_id string
	dst_chain_id string
	src_contract common.Address
	recipient    common.Address
	depositor    common.Address
	tokens       []common.Address
	amounts      []*big.Int
	src_nonce    *big.Int
	dst_nonce    *big.Int
	block_number uint64
	dst_tx_hash  string
}

func NewCrossChainReplyMsg(
	src_chain_id string,
	dst_chain_id string,
	src_contract common.Address,
	recipient common.Address,
	depositor common.Address,
	tokens []common.Address,
	amounts []*big.Int,
	src_nonce *big.Int,
	dst_nonce *big.Int,
	block_number uint64,
	dst_tx_hash string,
) CrossChainReplyMsg {
	crossChainReplyMsg := CrossChainReplyMsg{
		src_chain_id: src_chain_id,
		dst_chain_id: dst_chain_id,
		src_contract: src_contract,
		recipient:    recipient,
		depositor:    depositor,
		tokens:       tokens,
		amounts:      amounts,
		src_nonce:    src_nonce,
		dst_nonce:    dst_nonce,
		block_number: block_number,
		dst_tx_hash:  dst_tx_hash,
	}
	return crossChainReplyMsg
}

func (c *CrossChainReplyMsg) GetMsgPacket() ([]byte, error) {
	data, err := abi.CROSS_CHAIN_REPLY_INTERFACE.Pack(
		c.src_chain_id,
		c.dst_chain_id,
		c.src_contract,
		c.recipient,
		c.depositor,
		c.tokens,
		c.amounts,
		c.src_nonce,
		c.dst_nonce,
		c.block_number,
		c.dst_tx_hash,
	)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// get txt_type msg
type CrossChainRequestTypeData struct {
	request_type uint8
	src_chain_id string
	src_nonce    *big.Int
	dst_chain_id string
	dst_nonce    *big.Int
	validator    common.Address
}

func NewCCTxTypeMsg(
	request_type uint8,
	src_chain_id string,
	src_nonce *big.Int,
	dst_chain_id string,
	dst_nonce *big.Int,
) CrossChainRequestTypeData {
	crossChainRequestTypeData := CrossChainRequestTypeData{
		request_type: request_type,
		src_chain_id: src_chain_id,
		dst_chain_id: dst_chain_id,
		validator:    common.Address{},
		src_nonce:    src_nonce,
		dst_nonce:    dst_nonce,
	}
	return crossChainRequestTypeData
}

func (c *CrossChainRequestTypeData) GetMsgPacket() ([]byte, error) {
	data, err := abi.TX_TYPE_INTERFACE.Pack(
		c.request_type,
		c.src_chain_id,
		c.src_nonce,
		c.dst_chain_id,
		c.dst_nonce,
		c.validator,
	)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
