package digichain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type HexString string

type RawProposal struct {
	hash          string
	chain_id      string
	proposal_type string
	proposed_by   common.Address
	proposed_at   uint64
	data          HexString
	nonce         uint64
	signature     Signature
}

func NewRawProposal(
	hash, chain_id, proposal_type string, proposed_by common.Address, proposed_at uint64, data HexString,
) RawProposal {
	// at the time broadcasting proposal, add account nonce
	// after heading nonce, sign it
	// crosschainRequest := crosschainTypes.NewCrosschainRequestFromMsg(msg)
	// err, ethSigner, signature := e.CreateCrosschainLockedConfirmationSignature(msg)
	// if err != nil {
	// 	e.logger.WithFields(log.Fields{"error": err}).Fatalf("Error in Adding Signature")
	// 	return nil, err
	// }
	// msg.EthSigner = ethSigner
	// msg.Signature = signature
	return RawProposal{
		hash:          hash,
		chain_id:      chain_id,
		proposal_type: proposal_type,
		proposed_by:   proposed_by,
		proposed_at:   proposed_at,
		data:          data,
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Params interface{}
type RequestBody struct {
	JsonRpcId string `json:"id"`
	Params    Params `json:"params"`
	Method    string `json:"method"`
}

func (dc *DigiChainClient) NewRequestBody(method string, params Params) RequestBody {
	return RequestBody{
		JsonRpcId: "1",
		Params:    params,
		Method:    method,
	}
}

type GetContractParams struct {
	ChainIDs []string `json:"chain_ids"`
}

type GetCrossChainRequest struct {
	Validator string `json:"validator"`
}

type IsCrossChainRequestBroadcasted struct {
	Validator     common.Address `json:"validator"`
	SrcChainId    string         `json:"src_chain_id"`
	SrcChainNonce string         `json:"src_nonce"`
}

type GetAccountParams struct {
	Address common.Address `json:"address"`
}

type GetProposals struct {
	ChainIDs []string `json:"chain_ids"`
}

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
	V uint8  `json:"v"`
}

func (s *Signature) ToHex() (string, error) {
	// Convert R and S components to 64-character hexadecimal strings
	rHex := fmt.Sprintf("%064s", s.R)
	sHex := fmt.Sprintf("%064s", s.S)
	// Convert V component to a two-digit hexadecimal string
	vHex := fmt.Sprintf("%02x", s.V)
	// Combine the components into the final hexadecimal representation
	hexSig := fmt.Sprintf("0x%s%s%s", rHex, sHex, vHex)
	return hexSig, nil
}

type Transaction struct {
	Data      string         `json:"data"`
	Hash      string         `json:"hash"`
	ChainId   string         `json:"chain_id"`
	CreatedAt uint64         `json:"created_at"`
	Type      string         `json:"tx_type"`
	From      common.Address `json:"from"`
	Nonce     string         `json:"nonce"`
	Signature Signature      `json:"signature"`
}

type BroadcastTxParams struct {
	Transaction Transaction `json:"transaction"`
}

type TxResponse struct {
	Data struct {
		TxHash string `json:"tx_hash"`
	} `json:"data"`
	ID string `json:"id"`
}
