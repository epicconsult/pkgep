package pkgep

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisStore interface {
	SaveJson(string, string) error
	LoadJson(string, string) (string, error)
}

type redisClient struct{}

func NewRedisStore() RedisStore {
	return &redisClient{}
}

func newConn() (*redis.Client, error) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New("redis connection was refused")
	}

	return rdb, nil
}

func (r *redisClient) SaveJson(keyName, jsond string) error {
	ctx := context.Background()

	rdb, err := newConn()
	if err != nil {
		return err
	}

	defer rdb.Close()

	statusCmd := rdb.JSONSet(ctx, keyName, "$", jsond)

	_, err = statusCmd.Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisClient) LoadJson(keyName string, accessor string) (string, error) {
	ctx := context.Background()

	rdb, err := newConn()
	if err != nil {
		return "", err
	}

	defer rdb.Close()

	jsonCmd := rdb.JSONGet(ctx, keyName, "$."+accessor)

	result, err := jsonCmd.Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
