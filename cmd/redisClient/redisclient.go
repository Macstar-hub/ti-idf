package redisclient

import (
	"context"
	"fmt"
	"tf-idf/cmd/logger"
	"time"

	// "fmt"
	"log"
	"strconv"

	// "tf-idf/cmd/logger"
	// "time"

	"github.com/redis/go-redis/v9"
)

const (
	logFilePath = "/Users/Shared/codes.dir/go.dir/git.dir/ti-idf/logs/redis/"
	logPrefix   = ".log"
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
	status := client.Set(ctx, key, value, 0)

	if status != nil {
		log.Printf("%s\n", key, value, status)
	} else {
		return
	}
	logger.Logger(logFilePath, logPrefix, fmt.Sprintf("%s", status), "debug")
}

func RedisGetOPS(key string) int {
	startTime := time.Now()
	value, status := client.Get(ctx, key).Result()
	if status != nil {
		log.Printf(key)
		logger.Logger(logFilePath, logPrefix, fmt.Sprintf("%s", status), "error")
	}
	intValue, _ := strconv.Atoi(value)
	// fmt.Printf("Latency to make get key: '%s' in redis client function: ", key, time.Since(startTime))
	logger.Logger(logFilePath, logPrefix, fmt.Sprintf("Latency to make get key: '%s' in redis client function: ", key, time.Since(startTime)), "debug")
	return intValue
}
