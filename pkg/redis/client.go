package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func NewClient(addr, password string, db int) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return client
}

func BanUser(username string) error {
	ctx := context.Background()
	return client.Set(ctx, "banned:"+username, "true", 0).Err()
}

func IsUserBanned(username string) (bool, error) {
	ctx := context.Background()
	banned, err := client.Get(ctx, "banned:"+username).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return banned == "true", nil
}
