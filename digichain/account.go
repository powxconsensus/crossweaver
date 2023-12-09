package digichain

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type AccountRes struct {
	Account struct {
		Address        common.Address `json:"address"`
		TxNonce        string         `json:"tx_nonce"`
		ProposalNonce  string         `json:"proposal_nonce"`
		IsKYCDone      bool           `json:"is_kyc_done"`
		Name           string         `json:"name"`
		Country        string         `json:"country"`
		Mobile         string         `json:"mobile"`
		AadharNo       string         `json:"aadhar_no"`
		KYCCompletedAt uint64         `json:"kyc_completed_at"`
		UPIID          string         `json:"upi_id"`
		Transactions   []interface{}  `json:"transactions"`
	} `json:"account"`
	ID string `json:"id"`
}

func (dc *DigiChainClient) FetchAccount(address common.Address) (AccountRes, error) {
	cp := GetAccountParams{
		Address: address,
	}
	rq := dc.NewRequestBody("get_account", cp)
	jsonData, err := json.Marshal(rq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return AccountRes{}, err
	}
	data, err := dc.PostCall(jsonData)
	if err != nil {
		fmt.Println("Error while fetching account: ", err)
		return AccountRes{}, err
	}
	var account AccountRes
	err = json.Unmarshal([]byte(data), &account)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return AccountRes{}, err
	}

	return account, nil
}
