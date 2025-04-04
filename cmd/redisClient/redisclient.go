package redisclient

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisHost = "localhost:6379"
var client = redis.NewClient(&redis.Options{
	Network:    "tcp",
	Addr:       redisHost,
	ClientName: "priceCache",
	Password:   "",
	Username:   "",
	DB:         0,
})

func RedisChecker() bool {
	pong, err := client.Ping(ctx).Result()

	if err != nil {
		log.Fatal("Cannot make pong: ", err)
	}
	if pong == "PONG" {
		return true
	}

	return false
}

func RedisSetOPS(key string, value int) {
	err := client.Set(ctx, key, value, 0)

	if err != nil {
		log.Printf("Cannot set key %s and value %s with error: \n", key, value, err)
	} else {
		return
	}
}

func RedisGetOPS(key string) int {
	startTime := time.Now()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Canntot Find Value Of Key: %s", key, err)
	}
	intValue, _ := strconv.Atoi(value)
	fmt.Printf("Latency to make get key: '%s' in redis client function: ", key, time.Since(startTime))
	return intValue
}
