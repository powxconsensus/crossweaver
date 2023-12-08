package digichain

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ContractConfig struct {
	Configs map[string]struct {
		ChainType          int    `json:"chain_type"`
		ContractAddress    string `json:"contract_address"`
		LastProcessedNonce int    `json:"last_proccessed_nonce"`
		LastProcessedBlock int    `json:"last_processed_block"`
		StartBlock         int    `json:"start_block"`
	} `json:"configs"`
	ID string `json:"id"`
}

type CrossChainRequest struct {
	Payload    string      `json:"payload"`
	Sigs       []Signature `json:"sigs"`
	SrcChainId string      `json:"src_chain_id"`
	SrcNonce   *big.Int    `json:"src_nonce"`
	DstChainId string      `json:"dst_chain_id"`
}
type CrossChainRequestsResponse struct {
	CrossChainRequests []CrossChainRequest `json:"crosschain_withdraw_requests"`
	ID                 string              `json:"id"`
}

func (dc *DigiChainClient) FetchContractConfig(chainIds []string) (ContractConfig, error) {
	cp := GetContractParams{
		ChainIDs: chainIds,
	}
	rq := dc.NewRequestBody("get_contracts_config", cp)
	jsonData, err := json.Marshal(rq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ContractConfig{}, err
	}
	data, err := dc.PostCall(jsonData)
	if err != nil {
		fmt.Println("Error while fetching contracts config: ", err)
		return ContractConfig{}, err
	}
	var contractConfig ContractConfig
	err = json.Unmarshal([]byte(data), &contractConfig)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ContractConfig{}, err
	}

	return contractConfig, nil
}

func (dc *DigiChainClient) FetchCrossChainRequest(validator common.Address) ([]CrossChainRequest, error) {
	cp := GetCrossChainRequest{
		Validator: validator.String(),
	}
	rq := dc.NewRequestBody("get_crosschain_requests", cp)
	jsonData, err := json.Marshal(rq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return []CrossChainRequest{}, err
	}
	data, err := dc.PostCall(jsonData)
	if err != nil {
		fmt.Println("Error while fetching contracts config: ", err)
		return []CrossChainRequest{}, err
	}
	var res CrossChainRequestsResponse
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return []CrossChainRequest{}, err
	}
	return res.CrossChainRequests, nil
}

func (dc *DigiChainClient) IsCrosschainRequestBroadcasted(validator common.Address, srcChainId string, srcNonce *big.Int) (bool, error) {
	cp := IsCrossChainRequestBroadcasted{
		Validator:     validator,
		SrcChainId:    srcChainId,
		SrcChainNonce: srcNonce.Uint64(),
	}
	rq := dc.NewRequestBody("is_broadcasted", cp)
	jsonData, err := json.Marshal(rq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return false, err
	}
	data, err := dc.PostCall(jsonData)
	if err != nil {
		fmt.Println("Error while fetching contracts config: ", err)
		return false, err
	}
	var res struct {
		IsBroadCasted bool   `json:"is_broadcasted"`
		ID            string `json:"id"`
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return false, err
	}
	return res.IsBroadCasted, nil
}
