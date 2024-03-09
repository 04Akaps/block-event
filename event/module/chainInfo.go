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

	Topics   []common.Hash
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
	}

	c.Topics = []common.Hash{
		common.BytesToHash(crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Bytes()),
	}

	return c
}
