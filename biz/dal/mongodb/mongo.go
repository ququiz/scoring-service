package mongodb

import (
	"context"
	"fmt"
	"ququiz/lintang/scoring-service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Mongodb struct {
	Conn      *mongo.Database
	FakerConn *mongo.Database
}

func NewMongo(cfg *config.Config) *Mongodb {
	zap.L().Info(fmt.Sprintf("url mongo: %s", cfg.Mongodb.MongoURL))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongodb.MongoURL))
	if err != nil {
		zap.L().Fatal("mongo.Connect()", zap.Error(err))
	}

	db := client.Database(cfg.Mongodb.Database)

	writeClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongodb.MongoWriteURL))
	if err != nil {
		zap.L().Fatal("mongo.Connect() (write db)", zap.Error(err))
	}
	writeDb := writeClient.Database(cfg.Mongodb.Database)

	return &Mongodb{db, writeDb}
}

func (db *Mongodb) Close(ctx context.Context) {
	db.Conn.Client().Disconnect(ctx)
}
