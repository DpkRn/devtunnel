package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env from current working directory (run the binary from repo root, or set env in Docker/systemd).
	_ = godotenv.Load()
}

type MongoDBCfg struct {
	URI       string
	DBName    string
	ColName   string
	ColPrefix string
}

func (c config) MongoDB() MongoDBCfg {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	db := os.Getenv("MONGODB_DB")
	if db == "" {
		db = "tunnel"
	}
	col := os.Getenv("MONGODB_COLLECTION")
	if col == "" {
		col = "events"
	}
	return MongoDBCfg{
		URI:     uri,
		DBName:  db,
		ColName: col,
	}
}

func (c MongoDBCfg) URIFunc() string {
	return c.URI
}

func (c MongoDBCfg) DBNameFunc() string {
	return c.DBName
}

func (c MongoDBCfg) ColNameFunc() string {
	return c.ColName
}

func (c MongoDBCfg) ColPrefixFunc() string {
	return c.ColPrefix
}
