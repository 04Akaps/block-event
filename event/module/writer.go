package module

import (
	"github.com/04Akaps/block-event/log"
	"github.com/04Akaps/block-event/repository"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"time"
)

type Writer struct {
	repository *repository.Repository
	writerChan <-chan *WriterChan
}

func NewWriter(
	repository *repository.Repository,
	writerChan <-chan *WriterChan,
) *Writer {
	w := &Writer{repository: repository, writerChan: writerChan}

	go w.lookingEvent()

	return w
}

func (w *Writer) lookingEvent() {
	for {
		event := <-w.writerChan

		if event.EventName == Transfer {
			w.Transfer(event)
		}

	}
}

// Transfer Transfer(address indexed from, address indexed to, uint256 value)
func (w Writer) Transfer(data *WriterChan) {
	//event * types.Log, block * types.Header, txs * types.Transaction

	eventName := data.EventName
	event := data.Event
	tx := event.TxHash.Hex()

	if len(data.Event.Topics) == 1 {
		log.InfoLog("Not Existed Topic", "event", eventName)
	} else {
		from := common.BytesToAddress(event.Topics[1][:])
		to := common.BytesToAddress(event.Topics[2][:])
		value := new(big.Int).SetBytes(event.Data[:0x20])
		collection := data.Txs.To()

		go func() {
			if v, err := ToJSON(&Tx{
				Tx:        tx,
				Time:      time.Now().Unix(),
				Index:     data.Index,
				ChainName: data.ChainName,
				EventName: eventName,
				ChainID:   data.ChainID,
				Address:   event.Address.Hex(),
			}); err != nil {
				log.ErrLog("Failed To ToJSON In Transfer")
			} else {
				w.repository.Mongo.SaveTx(tx, eventName, v)
			}
		}()

		go func() {
			w.repository.Mongo.Transfer(
				hexutil.Encode(from[:]),
				hexutil.Encode(to[:]),
				hexutil.Encode(collection[:]),
				value.String(),
			)
		}()

		// TODO MySql

	}

}
