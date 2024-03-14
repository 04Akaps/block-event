package mongo

import (
	"context"
	"github.com/04Akaps/block-event/event/module"
	"github.com/04Akaps/block-event/init/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	cfg *config.Config

	client *mongo.Client

	transfer *mongo.Collection
	tx       *mongo.Collection
}

func NewMongo(cfg *config.Config) (*Mongo, error) {
	m := &Mongo{
		cfg: cfg,
	}

	var err error
	ctx := context.Background()
	dbConfig := cfg.Mongo

	if m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.Uri)); err != nil {
		log.Error("Failed to connect mongo", "config", dbConfig.Uri, "err", err)
		return nil, err
	} else if err = m.client.Ping(ctx, nil); err != nil {
		log.Error("Failed mongo ping", "err", err)
		return nil, err
	} else {
		db := m.client.Database(dbConfig.DB, nil)

		m.transfer = db.Collection("transfer")
		m.tx = db.Collection("tx")

		return m, nil
	}
}

func (m *Mongo) Transfer(from, to, collection, value string) {
	ctx := context.Background()
	now := time.Now().Unix()

	baseFilter := bson.M{"user": from, "collection": collection}

	valueDecimal, _ := decimal.NewFromString(value)

	var fromTransfer module.TransferType
	err := m.transfer.FindOne(ctx, baseFilter).Decode(&fromTransfer)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("First User", "from", from, "collection", collection)
		} else {
			log.Error("Failed To Get FromUser", "from", from, "collection", collection)
		}
	} else {
		beforeFromBalance, _ := decimal.NewFromString(fromTransfer.Balance)

		fromTransfer.Balance = beforeFromBalance.Sub(valueDecimal).String()
		fromTransfer.UpdatedTime = now

		if v, err := module.ToJSON(fromTransfer); err != nil {
			log.Error("Failed To ToJson At FromTransfer Struct", "collection", collection, "from", from)
		} else if r, err := m.transfer.UpdateOne(ctx, baseFilter, bson.M{"$set": v}); err != nil {
			log.Error("Failed To Update FromTransfer Struct", "collection", collection, "from", from)
		} else {
			log.Info("Success To Update FromTransfer Struct", "collection", collection, "from", from, "modified", r.ModifiedCount)
		}
	}

	baseFilter["user"] = to

	var toTransfer module.TransferType
	err = m.transfer.FindOne(ctx, baseFilter).Decode(&toTransfer)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create New Document
			toTransfer = module.TransferType{
				UpdatedTime: now,
				CreatedTime: now,
				Collection:  collection,
				User:        to,
			}
		} else {
			log.Error("Failed To Get ToUser", "to", to, "collection", collection)
		}
	} else {
		beforeToBalance, _ := decimal.NewFromString(toTransfer.Balance)
		toTransfer.Balance = beforeToBalance.Add(valueDecimal).String()
		toTransfer.UpdatedTime = now
	}

	if v, err := module.ToJSON(&toTransfer); err != nil {
		log.Error("Failed To ToJson At ToTransfer Struct", "collection", collection, "from", from)
	} else if r, err := m.transfer.UpdateOne(ctx, baseFilter, bson.M{"$set": v}); err != nil {
		log.Error("Failed To Update ToTransfer Struct", "collection", collection, "from", from)
	} else {
		log.Info("Success To Update ToTransfer Struct", "collection", collection, "from", from, "modified", r.ModifiedCount)
	}

}

func (m *Mongo) SaveTx(tx, eventName string, value interface{}) {

	filter := bson.M{"tx": tx, "eventName": eventName}
	update := bson.M{"$set": value}

	if r, err := m.tx.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true)); err != nil {
		log.Error("Failed To Update Tx", tx)
	} else {
		log.Info("Success To Save Tx", "tx", tx, "modified", r.ModifiedCount)
	}

}
