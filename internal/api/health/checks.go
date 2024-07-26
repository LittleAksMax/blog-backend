package health

import (
	"context"
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

// TODO: redis client parameter
func checkCacheHealth(ctx context.Context) (string, bool) {
	// TODO: implement for Redis
	return "ok", true
}
