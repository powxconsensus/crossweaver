package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/digilabs/crossweaver/digichain"
	"github.com/digilabs/crossweaver/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.json"
const DefaultBlockTimeout = int64(180) // 3 minutes
const DefaultChainLimit = int(10)
const DefaultRetryLimit = int(10)
const DefaultBlockRetryInterval = time.Second * 5
const DefaultSleeptimeforBatchRequest = time.Second * 5
const DefaultSleeptimetoStartbatchTx = time.Second * 10

type Config struct {
	Chains       []Chain      `json:"chains"`
	ChainSpecs   []ChainSpecs `json:"chainSpecs"`
	GlobalConfig GlobalCfg    `json:"globalConfig,omitempty"`
}

// GlobalCfg for non chain specific
type GlobalCfg struct {
	DbPath        string `json:"dbPath,omitempty"`
	EthPrivateKey string `json:"ethPrivateKey,omitempty"`
	DigiChainRPC  string `json:"digirpc,omitempty"`
	From          string `json:"from"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type Chain struct {
	ChainId               string `json:"chainId"`
	ChainName             string `json:"chainName"`
	ChainType             string `json:"chainType"`
	ChainRpc              string `json:"chainRpc"`
	BlocksToSearch        uint64 `json:"blocksToSearch,omitempty"`
	BlockTime             string `json:"blockTime,omitempty"`
	ChainApi              string `json:"chainApi,omitempty"`
	StartBlock            uint64 `json:"startBlock"`
	From                  string `json:"from"`
	KeyPath               string `json:"keyPath"`
	ConfirmationsRequired uint64 `json:"confirmationsRequired"`
	NumOfThreads          int    `json:"numOfThreads"`
}

type ChainSpecs struct {
	ChainId                 string          `json:"chain_id"`
	ChainName               string          `json:"chainName"`
	Symbol                  string          `json:"symbol"`
	ChainType               types.ChainType `json:"chainType"`
	ConfirmationsRequired   uint64          `json:"confirmationsRequired"`
	StartBlock              uint64          `json:"startBlock"`
	StartEventNonce         uint64          `json:"startEventNonce"`
	LastObservedValsetNonce uint64          `json:"lastObservedValsetNonce"`
	ChainRpc                string          `json:"chainRpc,omitempty"`
	ChainApi                string          `json:"chainApi,omitempty"`
	BlocksToSearch          uint64          `json:"blocksToSearch,omitempty"`
	BlockTime               string          `json:"blockTime,omitempty"`
	ContractAddress         string          `json:"contractAddress"`
	From                    string          `json:"from"`
	KeyPath                 string          `json:"keyPath"`
	NumOfThreads            int             `json:"numOfThreads"`
}

type Response struct {
	ChainConfig []ChainSpecs `json:"chainConfig"`
}

func NewConfig() *Config {
	return &Config{
		Chains:     []Chain{},
		ChainSpecs: []ChainSpecs{},
	}
}

func (c *Config) ToJSON(file string) *os.File {
	var (
		newFile *os.File
		err     error
	)

	var raw []byte
	if raw, err = json.Marshal(*c); err != nil {
		//log.Warn("error marshalling json", "err", err)
		os.Exit(1)
	}

	newFile, err = os.Create(file)
	if err != nil {
		//log.Warn("error creating config file", "err", err)
	}
	_, err = newFile.Write(raw)
	if err != nil {
		//log.Warn("error writing to config file", "err", err)
	}
	//if err := newFile.Close(); err != nil {
	//	//log.Warn("error closing file", "err", err)
	//}
	return newFile
}

func (c *Config) validate() error {
	for _, chain := range c.Chains {
		if chain.ChainId == "" {
			//log.Error("Required filed ID in Supported Chains ", "chains", chain)
			return fmt.Errorf("required field chain.ID empty for chain %s", chain.ChainId)
		}
		if chain.ChainName == "" {
			//log.Error("Required filed Name in Supported Chains ", "chains", chain)
			return fmt.Errorf("required field chain.Name empty for chain %s", chain.ChainName)
		}
		if chain.ChainType == "" {
			//log.Error("Required filed Type in Supported Chains ", "chains", chain)
			return fmt.Errorf("required field chain.Type empty for chain %s", chain.ChainType)
		}
		if chain.ChainRpc == "" {
			//log.Error("Required filed Type in Supported Chains ", "chains", chain)
			return fmt.Errorf("required field chain.RPC empty for chain %s", chain.ChainRpc)
		}
	}
	return nil
}

func (c *Config) FetchSpecs(digiChainClient digichain.DigiChainClient) error {
	ctx := context.Background()
	config, err := digiChainClient.FetchContractConfig(make([]string, 0))
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Error in Fetching values from router client ")
	}
	for chainId, v1 := range config.Configs {
		lnonce, err := strconv.ParseUint(v1.LastProcessedNonce, 10, 64)
		if err != nil {
			continue
		}
		for _, v2 := range c.Chains {
			if chainId == v2.ChainId {
				Cx := ChainSpecs{
					ChainId:                 chainId,
					ChainName:               v2.ChainName,
					ChainType:               types.EVM_CHAIN, //TODO: hardcoded to evm only
					LastObservedValsetNonce: lnonce,
					StartBlock:              uint64(v1.StartBlock),
					ChainRpc:                v2.ChainRpc,
					From:                    v2.From,
					BlocksToSearch:          v2.BlocksToSearch,
					BlockTime:               v2.BlockTime,
					ConfirmationsRequired:   v2.ConfirmationsRequired,
					NumOfThreads:            v2.NumOfThreads,
					ContractAddress:         v1.ContractAddress,
				}
				c.ChainSpecs = append(c.ChainSpecs, Cx)
			}
		}
	}
	ctx.Done()
	return nil
}

func GetConfig(ctx *cli.Context) (*Config, error) {
	var fig Config
	path := DefaultConfigPath
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, &fig)
	if err != nil {
		return &fig, err
	}
	// GET .ENV
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	// GET .ENV
	err = fig.validate()
	if err != nil {
		return nil, err
	}
	return &fig, nil
}

func loadConfig(file string, config *Config) error {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	if ext == ".json" {
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}
	return nil
}
