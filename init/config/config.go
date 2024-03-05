package config

import (
	"github.com/naoina/toml"
	"os"
)

type Node struct {
	Dial       string
	StartBlock uint64
	ChainName  string
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
