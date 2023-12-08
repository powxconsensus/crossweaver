package digichain

import (
	"crypto/ecdsa"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

type DigiChainClient struct {
	chainId    string
	rpc        string
	privateKey string
	from       string
}

func privateKeyToAddress(privateKeyString string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		err = errors.Wrap(err, "failed to hex-decode Ethereum ECDSA Private Key")
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	// address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	// fmt.Println(address)
	hash := crypto.NewKeccakState()
	hash.Write(publicKeyBytes[1:])
	return hexutil.Encode(hash.Sum(nil)[12:]), nil

}

func NewDigiChainClient(rpc string, privateKey string) DigiChainClient {
	from, err := privateKeyToAddress(privateKey)
	if err != nil {
		panic(err)
	}
	voyagerEventProcessor := DigiChainClient{
		chainId:    "11",
		rpc:        rpc,
		from:       from,
		privateKey: privateKey,
	}
	return voyagerEventProcessor
}

func (dc *DigiChainClient) FromAddress() string {
	return dc.from
}

func (dc *DigiChainClient) GetChainId() string {
	return dc.chainId
}
