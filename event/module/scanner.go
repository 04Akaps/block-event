package module

import (
	"context"
	"github.com/04Akaps/block-event/init/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"sync/atomic"
	"time"
)

type Scanner struct {
	cfg *config.Config

	chainInfo *ChainInfo

	TopBlock   atomic.Value
	StartBlock atomic.Value

	FilterQuery ethereum.FilterQuery
}

func NewScanner(
	cfg *config.Config,
	info *ChainInfo,
	startBlock int64,
) (chan []types.Log, *Scanner) {
	s := &Scanner{
		cfg:       cfg,
		chainInfo: info,
	}

	logs := make(chan []types.Log, 100)

	go s.scan(startBlock, logs)

	return logs, s
}

func (s *Scanner) scan(
	startBlock int64,
	eventLog chan<- []types.Log,
) {
	s.FilterQuery = ethereum.FilterQuery{
		Addresses: s.chainInfo.scanList,
		Topics:    [][]common.Hash{s.chainInfo.GetEventsToCatch()},
	}

	start, end := startBlock, int64(0)

	ticker := time.NewTicker(1e8)
	stop := make(chan struct{})

	go func() {
		defer close(stop)
		for {
			select {
			case <-stop:
				return
			default:
				if maxBlock, err := s.chainInfo.Client.BlockNumber(context.Background()); err != nil {
					log.Error("Get Current Block", "crit", err)
					return
				} else {
					end = int64(maxBlock)

					if end <= startBlock {
						continue
					}

					s.FilterQuery.FromBlock = big.NewInt(start)
					s.FilterQuery.ToBlock = big.NewInt(end)

					s.TopBlock.Store(big.NewInt(end))
					s.StartBlock.Store(big.NewInt(start))

					tryCount := 1

				BackRetry:
					if logs, err := s.chainInfo.Client.FilterLogs(context.Background(), s.FilterQuery); err != nil {
						// Filter로그를 못가져 온 것이기 때문에,
						// To만 바꿔서 재시도,

						newTo := big.NewInt(int64(end - 1))
						newFrom := big.NewInt(start - 1)

						s.FilterQuery.ToBlock = newTo
						s.FilterQuery.FromBlock = newFrom

						s.TopBlock.Store(newTo)
						s.StartBlock.Store(newFrom)

						log.Info("call FilterLogs Again", "startBlock", start, "end", end)

						tryCount++
						goto BackRetry
					} else if len(logs) > 0 {
						eventLog <- logs
					}

					startBlock = end + 1
				}

				<-ticker.C
			}
		}
	}()
}

func (s *Scanner) StartCatchEvent(eventChan <-chan []types.Log, writerChan chan<- *WriterChan) {
	for events := range eventChan {

		ctx := context.Background()

		blocks := make(map[uint64]*types.Header)
		txs := make(map[common.Hash]*types.Transaction)

		for i, event := range events {
			if _, ok := blocks[event.BlockNumber]; !ok {
				// 읽지 않았다면,
				if header, err := s.chainInfo.Client.HeaderByNumber(ctx, new(big.Int).SetInt64(int64(event.BlockNumber))); err != nil {
					blocks[event.BlockNumber] = header
				}
			}

			if _, ok := txs[event.TxHash]; !ok {
				if tx, pending, err := s.chainInfo.Client.TransactionByHash(ctx, event.TxHash); err == nil {
					if !pending {
						txs[event.TxHash] = tx
					}
				}
			}

			topic := event.Topics[0]

			if eventName, exist := s.chainInfo.CheckEventToCatch(topic); !exist {
				log.Info("Failed To Find Event To Catch")
			} else {
				writerChan <- &WriterChan{
					EventName: eventName,
					Index:     int64(i),
					ChainName: s.chainInfo.ChainName,
					ChainID:   s.chainInfo.ChainID.Int64(),
					Event:     &event,
					Block:     blocks[event.BlockNumber],
					Txs:       txs[event.TxHash],
				}
			}

		}
	}
}
