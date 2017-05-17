package redismanager

import (
	"github.com/go-redis/redis"
	"errors"
)

func InitRedisClient() (redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr    : "localhost:6379",
		Password: "",
		DB      : 0, //default
	})

	_, err := client.Ping().Result()
	if( err != nil ){
		return *client, errors.New("Cannot Initialize Redis Client : " + err.Error())
	}

	return *client, err
}
