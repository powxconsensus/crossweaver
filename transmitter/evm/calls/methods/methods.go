package methods

import (
	"strings"

	eAbi "github.com/ethereum/go-ethereum/accounts/abi"
)

type Method string

var HandleRequest Method = "handleRequest"
var Executed Method = "executed"

func (m *Method) GetAbiPackBytes(abi string, args ...interface{}) ([]byte, error) {
	abiParsed, err := eAbi.JSON(strings.NewReader(abi))
	if err != nil {
		return nil, err
	}
	// Pack the arguments
	packedArguments, err := abiParsed.Pack(string(*m), args...)
	if err != nil {
		return nil, err
	}
	return packedArguments, nil
}
