package abi

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	U256                          abi.Type
	U8                            abi.Type
	U64                           abi.Type
	Bytes32                       abi.Type
	Address                       abi.Type
	Bytes                         abi.Type
	String                        abi.Type
	AddressArray                  abi.Type
	U256Array                     abi.Type
	CROSS_CHAIN_REQUEST_INTERFACE abi.Arguments
	CROSS_CHAIN_REPLY_INTERFACE   abi.Arguments
	TX_TYPE_INTERFACE             abi.Arguments
)

func init() {
	U256, _ = abi.NewType("uint256", "", nil)
	U64, _ = abi.NewType("uint64", "", nil)
	U8, _ = abi.NewType("uint8", "", nil)
	Bytes32, _ = abi.NewType("bytes32", "", nil)
	Address, _ = abi.NewType("address", "", nil)
	Bytes, _ = abi.NewType("bytes", "", nil)
	String, _ = abi.NewType("string", "", nil)
	AddressArray, _ = abi.NewType("address[]", "", nil)
	U256Array, _ = abi.NewType("uint256[]", "", nil)

	CROSS_CHAIN_REQUEST_INTERFACE = abi.Arguments{
		abi.Argument{Type: String},       //src_chain_id
		abi.Argument{Type: String},       //dst_chain_id
		abi.Argument{Type: Address},      // src_contract
		abi.Argument{Type: Address},      // recipient
		abi.Argument{Type: Address},      // depositor
		abi.Argument{Type: AddressArray}, // tokens
		abi.Argument{Type: U256Array},    // amounts
		abi.Argument{Type: U256},         // src_nonce
		abi.Argument{Type: U64},          // src_block_number
		abi.Argument{Type: String},       // src_tx_hash
	}

	CROSS_CHAIN_REPLY_INTERFACE = abi.Arguments{
		abi.Argument{Type: String},       //src_chain_id
		abi.Argument{Type: String},       //dst_chain_id
		abi.Argument{Type: Address},      // src_contract
		abi.Argument{Type: Address},      // recipient
		abi.Argument{Type: Address},      // depositor
		abi.Argument{Type: AddressArray}, // tokens
		abi.Argument{Type: U256Array},    // amounts
		abi.Argument{Type: U256},         // src_nonce
		abi.Argument{Type: U256},         // dst_nonce
		abi.Argument{Type: U64},          // dst_block_number
		abi.Argument{Type: String},       // dst_tx_hash
	}

	TX_TYPE_INTERFACE = abi.Arguments{
		abi.Argument{Type: U8},      //request type
		abi.Argument{Type: String},  //src_chain_id
		abi.Argument{Type: U256},    // src_nonce
		abi.Argument{Type: String},  //dst_chain_id
		abi.Argument{Type: U256},    // dst_nonce
		abi.Argument{Type: Address}, // validator
	}
}
