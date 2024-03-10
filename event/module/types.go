package module

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type WriterChan struct {
	EventName string             `json:"eventName"`
	Block     *types.Header      `json:"blocks"`
	Txs       *types.Transaction `json:"txs"`
}
