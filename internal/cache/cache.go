package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

type Config struct {
	Client *redis.Client
}

func InitCache(ctx context.Context, host string, port int, password string) *Config {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to Redis!")

	return &Config{
		Client: rdb,
	}
}

func (cfg *Config) CloseCache() {
	err := cfg.Client.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
