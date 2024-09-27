package redis

import (
	"context"
	"strconv"

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

func CacheLatestArticleID(articleID uint32) error {
	ctx := context.Background()
	err := client.LPush(ctx, "latest_article_id", articleID).Err()
	if err != nil {
		return err
	}
	err = client.LTrim(ctx, "latest_article_id", 0, 2-1).Err()
	return err

}

func GetLatestArticleIDs() ([]uint32, error) {
	ctx := context.Background()
	ids, err := client.LRange(ctx, "latest_article_id", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var articleIDs []uint32
	for _, idStr := range ids {
		id, _ := strconv.ParseUint(idStr, 10, 32)
		articleIDs = append(articleIDs, uint32(id))
	}
	return articleIDs, nil
}
