package gateway

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type DigiPayContract struct {
	client         *ethclient.Client
	DigiPayAddress common.Address
	From           common.Address
	PrivateKey     *ecdsa.PrivateKey
	log            *logrus.Entry
}

type Response struct {
	BlockHash   string   `json:"block_hash"`
	BlockHeight uint64   `json:"block_height"`
	Logs        []string `json:"logs"`
	Result      []byte   `json:"result"`
}

func NewDigiPayContract(
	client *ethclient.Client,
	digiPayAddress common.Address,
	privateKey *ecdsa.PrivateKey,
	log *logrus.Entry,
) *DigiPayContract {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	return &DigiPayContract{
		client:         client,
		PrivateKey:     privateKey,
		log:            log,
		DigiPayAddress: digiPayAddress,
		From:           crypto.PubkeyToAddress(*publicKeyECDSA),
	}
}

func (c *DigiPayContract) GetDigiPayContractAddress() common.Address {
	return c.DigiPayAddress
}

func (c *DigiPayContract) EstimateGasPrice() (*big.Int, error) {
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return gasPrice, nil
}

func (c *DigiPayContract) CallContract(data []byte) ([]byte, error) {
	msg := ethereum.CallMsg{
		From:     c.From,
		To:       &c.DigiPayAddress,
		Value:    big.NewInt(0),
		Data:     data,
		GasPrice: big.NewInt(0),
	}
	return c.client.CallContract(context.Background(), msg, nil)
}

func (c *DigiPayContract) EstimateGasLimit(data []byte, gasPrice *big.Int, value *big.Int) (uint64, error) {
	msg := ethereum.CallMsg{
		From:     c.From,
		To:       &c.DigiPayAddress,
		Value:    value,
		Data:     data,
		GasPrice: gasPrice,
	}
	return c.client.EstimateGas(context.Background(), msg)
}

func (c *DigiPayContract) SimulateTransaction(method string, contractAddress common.Address, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{From: c.From, To: &contractAddress, Data: data}
	out, err := c.client.EstimateGas(context.Background(), msg)
	if err != nil {
		c.log.Error("CallContract err: ", err)
		return 0, err
	}
	return out, nil
}

func (c *DigiPayContract) GetTransactionOpts(gasPrice *big.Int, gasLimit uint64, value *big.Int) (*bind.TransactOpts, error) {
	nonce, err := c.client.PendingNonceAt(context.Background(), c.From)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := c.client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	return auth, nil
}
