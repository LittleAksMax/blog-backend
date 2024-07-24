package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Config struct {
	Client *mongo.Client
	db     *mongo.Database
	Posts  *mongo.Collection
}

func InitDb(ctx context.Context, host string, port int, user string, passwd string, dbName string) *Config {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", user, passwd, host, port, dbName)

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	creds := options.Credential{Username: user, Password: passwd}
	opts := options.Client().ApplyURI(connStr).SetServerAPIOptions(serverAPI).SetAuth(creds)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// database
	db := client.Database(dbName)

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &Config{
		Client: client,
		db:     db,
		Posts:  db.Collection("posts"),
	}

}

func (cfg *Config) CloseDb() {
	if err := cfg.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
