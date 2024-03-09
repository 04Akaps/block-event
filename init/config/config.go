package config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
	"os"
)

type Node struct {
	Dial         string
	StartBlock   int64
	ChainName    string
	TokenAddress []common.Address
}

type Config struct {
	DB struct {
		Uri string
		DB  string
	}

	Nodes map[string]*Node

	Log struct {
		LogName string
	}
}

func NewConfig(f string) *Config {
	c := new(Config)

	if f, err := os.Open(f); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(f).Decode(c); err != nil {
		panic(err)
	} else {
		return c
	}

}
