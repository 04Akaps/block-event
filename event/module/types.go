package module

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"reflect"
)

type WriterChan struct {
	EventName string             `json:"eventName"`
	ChainName string             `json:"chainName"`
	Index     int64              `json:"index"`
	ChainID   int64              `json:"chainID"`
	Block     *types.Header      `json:"blocks"`
	Txs       *types.Transaction `json:"txs"`
	Event     *types.Log         `json:"event"`
}

type Tx struct {
	Tx        string `json:"tx"`
	Time      int64  `json:"time"`
	Index     int64  `json:"index"`
	ChainName string `json:"chainName"`
	EventName string `json:"eventName"`
	ChainID   int64  `json:"chainID"`
	Address   string `json:"address"`
}

type TransferType struct {
	User        string `json:"user"`
	Balance     string `json:"balance"`
	Collection  string `json:"collection"`
	UpdatedTime int64  `json:"updatedTime"`
	CreatedTime int64  `json:"createdTime"`
}

func ToJSON(t interface{}) (interface{}, error) {
	var v interface{}
	if bytes, err := json.Marshal(t); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bytes, &v); err != nil {
		return nil, err
	} else {
		jsonMap := v.(map[string]interface{})
		for key, value := range jsonMap {
			if reflect.TypeOf(value) == reflect.TypeOf(float64(0)) {
				jsonMap[key] = int64(value.(float64))
			}
		}

		return jsonMap, nil
	}
}

const (
	Transfer = "Transfer"
)
