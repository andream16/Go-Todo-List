package redismanager

import (
	"fmt"
	"github.com/go-redis/redis"
)

func InitRedisClient() redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr    : "localhost:6379",
		Password: "",
		DB      : 0, //default
	})

	pong, err := client.Ping().Result()
	if( err != nil ){
		fmt.Println("Cannot Initialize Redis Client ", err)
	}
	fmt.Println("Redis Client Successfully Initialized . . .", pong)

	return *client
}
