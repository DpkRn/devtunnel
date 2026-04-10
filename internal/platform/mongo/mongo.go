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
	GetLogs(ctx context.Context, tunnelID string, limit int64) ([]map[string]any, error)
	GetLogByID(ctx context.Context, id string) (map[string]any, error)
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

func (m *mongoDB) GetLogs(ctx context.Context, tunnelID string, limit int64) ([]map[string]any, error) {
	opts := options.Find().SetSort(map[string]int{"created_at": -1}).SetLimit(limit)

	cursor, err := m.RequestLogCollection.Find(ctx, map[string]any{
		"tunnel_id": tunnelID,
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]any
	err = cursor.All(ctx, &results)
	return results, err
}

func (m *mongoDB) GetLogByID(ctx context.Context, id string) (map[string]any, error) {
	var result map[string]any
	err := m.RequestLogCollection.FindOne(ctx, map[string]any{
		"_id": id,
	}).Decode(&result)
	return result, err
}
