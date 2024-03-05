package main

import (
	"flag"
	"github.com/04Akaps/block-event/init"
	"github.com/04Akaps/block-event/init/config"
	"github.com/04Akaps/block-event/log"
)

var envFlag = flag.String("config", "./config.toml", "config not found")

func main() {
	flag.Parse()
	cfg := config.NewConfig(*envFlag)
	log.SetLog(cfg.Log.LogName)
	init.StartApp(cfg)
}
