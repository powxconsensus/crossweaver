package digichain

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func (d *DigiChainClient) FetchTransactionByHash() {

}

func (d *DigiChainClient) BroadcastTx(msg RawProposal) (*TxResponse, error) {
	account, err := d.FetchAccount(msg.proposed_by)
	if err != nil {
		return nil, err
	}
	cp := BroadcastTxParams{
		Transaction: Transaction{
			Data:      string(msg.data),
			Hash:      msg.hash,
			ChainId:   msg.chain_id,
			CreatedAt: msg.proposed_at,
			Type:      msg.proposal_type,
			From:      msg.proposed_by,
			Nonce:     account.Account.TxNonce.Add(account.Account.TxNonce, big.NewInt(1)).Uint64(),
			Signature: msg.signature,
		},
	}
	rq := d.NewRequestBody("broadcast_transaction", cp)
	jsonData, err := json.Marshal(rq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}
	data, err := d.PostCall(jsonData)
	if err != nil {
		fmt.Println("Error while broadcasting tx: ", err)
		return nil, err
	}
	var res TxResponse
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}
	return &res, nil
}

func (d *DigiChainClient) SignRawTxAndBroadCast(rawTx RawProposal) (*TxResponse, error) {
	//TODO: sign
	rawTx.signature = Signature{
		R: big.NewInt(1).String(),
		S: big.NewInt(2).String(),
		V: 0,
	}
	return d.BroadcastTx(rawTx)
}
