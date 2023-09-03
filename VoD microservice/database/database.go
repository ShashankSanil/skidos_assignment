package database

import (
	"context"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect() error {
	uri := archaius.GetString("database.mongodb.uri", "")
	poolsize := archaius.GetInt64("database.mongodb.poolsize", 15) // by defaylt pool size is 15.

	var clientOptions *options.ClientOptions

	clientOptions = options.Client().ApplyURI(uri).SetMaxPoolSize(uint64(poolsize))

	clientlocal, err := mongo.Connect(context.TODO(), clientOptions)
	client = clientlocal
	if err != nil {
		return err
	}
	err = clientlocal.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	openlog.Info("Connected to Mongodb")
	return nil
}

//GetClient function
func GetClient() *mongo.Client { return client }
