package node

import (
	"context"
	"github.com/04Akaps/block-event/init/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Node struct {
	config *config.Node

	Client  *ethclient.Client
	Chain   string
	ChainID *big.Int
}

func NewNode(config *config.Node) (*Node, error) {
	n := &Node{
		config: config,
		Chain:  config.ChainName,
	}

	var err error
	ctx := context.Background()

	if n.Client, err = ethclient.DialContext(ctx, config.Dial); err != nil {
		return nil, err
	} else if n.ChainID, err = n.Client.ChainID(ctx); err != nil {
		return nil, err
	} else {
		return n, nil
	}

}
