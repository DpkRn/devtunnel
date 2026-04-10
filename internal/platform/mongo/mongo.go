package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Client interface {
	InsertTunnelLog(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	InsertRequestLog(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
}

type mongoDB struct {
	TunnelLogCollection  *mongo.Collection
	RequestLogCollection *mongo.Collection
}

type Config interface {
	URIFunc() string
	DBNameFunc() string
	ColNameFunc() string
	ColPrefixFunc() string
}

func NewMongoClient(mongoConfig Config) (Client, error) {
	fmt.Println("Connecting to MongoDB", mongoConfig.URIFunc())
	opts := options.Client().
		ApplyURI(mongoConfig.URIFunc()).
		SetConnectTimeout(10 * time.Second)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB")
	return &mongoDB{
		TunnelLogCollection:  client.Database(mongoConfig.DBNameFunc()).Collection("tunnel_logs"),
		RequestLogCollection: client.Database(mongoConfig.DBNameFunc()).Collection("request_logs"),
	}, nil
}

func (m *mongoDB) InsertTunnelLog(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return m.TunnelLogCollection.InsertOne(ctx, document)
}

func (m *mongoDB) InsertRequestLog(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return m.RequestLogCollection.InsertOne(ctx, document)
}
