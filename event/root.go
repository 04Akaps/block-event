package event

import (
	"github.com/04Akaps/block-event/event/module"
	"github.com/04Akaps/block-event/init/config"
	"github.com/04Akaps/block-event/repository"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
)

type TokenTransferEvent struct {
	cfg *config.Config

	chainInfos []*module.ChainInfo
	scanner    []*module.Scanner
	writer     *module.Writer

	repository *repository.Repository
}

func NewTokenTransferScanner(
	cfg *config.Config,
	repository *repository.Repository,
) *TokenTransferEvent {
	t := &TokenTransferEvent{cfg: cfg, repository: repository}

	var writerLog chan *module.WriterChan

	t.writer = module.NewWriter(repository, writerLog)

	for name, node := range cfg.Nodes {
		r := repository.NodeMap[name]

		chainInfo := module.NewChainInfo(r.Client, r.Chain, r.ChainID, node.TokenAddress)
		t.chainInfos = append(t.chainInfos, chainInfo)

		var scannerLog chan []types.Log
		var scanner *module.Scanner

		scannerLog, scanner = module.NewScanner(cfg, chainInfo, node.StartBlock)

		scanner.StartCatchEvent(scannerLog, writerLog)

		go func() {
			ticker := time.NewTicker(5e9)

			for range ticker.C {
				// 주기마다 최신 블록에 대한 로그를 설정
				eb := scanner.TopBlock.Load().(*big.Int)
				sb := scanner.StartBlock.Load().(*big.Int)

				log.Info(
					"chain", chainInfo.ChainName,
					"tokenAddress", node.TokenAddress,
					"startBlock", sb,
					"endBlock", eb,
				)
			}
		}()
	}

	return t
}
