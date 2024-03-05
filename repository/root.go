package repository

import (
	"github.com/04Akaps/block-event/init/config"
	"github.com/04Akaps/block-event/log"
	"github.com/04Akaps/block-event/repository/mongo"
	"github.com/04Akaps/block-event/repository/mysql"
	"github.com/04Akaps/block-event/repository/node"
)

type Repository struct {
	cfg *config.Config

	MySQL   *mysql.MySql
	Mongo   *mongo.Mongo
	NodeMap map[string]*node.Node
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg, NodeMap: make(map[string]*node.Node)}

	var err error

	if r.MySQL, err = mysql.NewMySql(cfg); err != nil {
		return nil, err
	} else if r.Mongo, err = mongo.NewMongo(cfg); err != nil {
		return nil, err
	} else {
		for chainName, nodeInfo := range cfg.Nodes {
			if node, err := node.NewNode(nodeInfo); err != nil {
				log.CritLog("Failed To Connect Node", err.Error())
			} else {
				r.NodeMap[chainName] = node
			}
		}

		log.InfoLog("success to connect repository")
		return r, nil
	}
}
