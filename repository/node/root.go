package node

import (
	"context"
	"github.com/04Akaps/block-event/init/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Node struct {
	config *config.Node

	client  *ethclient.Client
	chain   string
	chainID *big.Int
}

func NewNode(config *config.Node) (*Node, error) {
	n := &Node{
		config: config,
		chain:  config.ChainName,
	}

	var err error
	ctx := context.Background()

	if n.client, err = ethclient.DialContext(ctx, config.Dial); err != nil {
		return nil, err
	} else if n.chainID, err = n.client.ChainID(ctx); err != nil {
		return nil, err
	} else {
		return n, nil
	}

}
