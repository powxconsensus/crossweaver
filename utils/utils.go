package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/digilabs/crossweaver/digichain"
)

type ByteSlice []byte

func (b ByteSlice) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte("[]"), nil
	}
	encoded := make([]int, len(b))
	for i, val := range b {
		encoded[i] = int(val)
	}
	return json.Marshal(encoded)
}

func CreateCrosschainDBId(crosschainRequest digichain.CrossChainRequest) (string, error) {
	if len(crosschainRequest.SrcChainId) == 0 {
		return "", fmt.Errorf("invalid ID creation due to data missing")
	}
	return fmt.Sprintf("%v:%v", crosschainRequest.SrcChainId, crosschainRequest.SrcNonce), nil
}

func Remove0xPrefix(input string) string {
	if strings.HasPrefix(input, "0x") {
		return strings.TrimPrefix(input, "0x")
	}
	return input
}
