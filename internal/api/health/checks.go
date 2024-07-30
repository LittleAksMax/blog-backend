package health

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func checkDbHealth(ctx context.Context, db *mongo.Database) (string, bool) {
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := db.RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return "database couldn't be reached (ping)", false
	}
	return "ok", true
}

func checkCacheHealth(ctx context.Context, rdb *redis.Client) (string, bool) {
	// returned result is just PONG on success
	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		return "not ok", false
	}
	return "ok", true
}
