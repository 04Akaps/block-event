package module

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type ChainInfo struct {
	Client    *ethclient.Client
	ChainName string
	ChainID   *big.Int

	Topics   map[string]common.Hash
	scanList []common.Address
}

func NewChainInfo(
	client *ethclient.Client,
	chainName string,
	chainID *big.Int,
	scanList []common.Address,
) *ChainInfo {
	c := &ChainInfo{
		Client:    client,
		ChainName: chainName,
		ChainID:   chainID,
		scanList:  scanList,
		Topics:    make(map[string]common.Hash),
	}

	c.Topics["Transfer"] = common.BytesToHash(crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Bytes())

	return c
}

func (c *ChainInfo) GetEventsToCatch() []common.Hash {
	return func() []common.Hash {
		eventToCatchList := make([]common.Hash, 0)
		for _, event := range c.Topics {
			eventToCatchList = append(eventToCatchList, event)
		}
		return eventToCatchList
	}()
}

func (c *ChainInfo) CheckEventToCatch(hash common.Hash) (string, bool) {
	return func() (string, bool) {
		for topic, h := range c.Topics {
			if h == hash {
				return topic, true
			}
		}
		return "", false
	}()
}
