package repository

import (
	"fmt"
	"strconv"

	config "check-files/pkg/config"

	"github.com/go-redis/redis/v7"
)

type RedisRepository struct {
	Client *redis.Client
	Key    string
}

// NewRedisRepository creates a new instance of RedisRepository with a Redis client
func NewRedisRepository(forKey string) *RedisRepository {

	config := config.NewConfig()

	db, _ := strconv.ParseInt(config.GetWithDefault("redis.db", "0"), 10, 0)

	client := redis.NewClient(&redis.Options{
		Addr:     config.GetWithDefault("redis.host", "localhost") + ":" + config.GetWithDefault("redis.port", "6379"),
		Password: config.GetWithDefault("redis.password", ""), // no password set
		DB:       int(db),                                     // use default DB
	})
	return &RedisRepository{Client: client, Key: forKey}
}

// SaveValue saves a value to Redis
func (r *RedisRepository) SaveValue(hashCode string, value interface{}) {
	err := r.Client.HSet(r.Key, hashCode, value).Err()
	if err != nil {
		fmt.Println("Error setting hash value:", err)
	} else {
	}
}

func (r *RedisRepository) GetValue(hash string) (string, error) {
	return r.Client.HGet(r.Key, hash).Result()
}
