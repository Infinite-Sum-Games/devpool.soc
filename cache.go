package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() error {
	host := Config.Redis.Host
	port := Config.Redis.Port
	uname := Config.Redis.Username
	passwd := Config.Redis.Password
	resp := 3

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: uname,
		Password: passwd,
		DB:       0,
		Protocol: resp,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		Log.Fatal("[FAIL]: Heal-check failed for Redis", err)
		return err
	}
	Log.Info(
		fmt.Sprintf("[INFO]: Health-check successful for Redis. Response: %s", pong),
	)
	Rdb = rdb

	return nil
}

func GetToken(ctx context.Context, key string) (string, bool) {
	return "", true
}

func SetToken(ctx context.Context, key, token string,
	expiry time.Duration) {
	err := Rdb.Set(ctx, key, token, expiry).Err()
	if err != nil {
	}
}

func GetInstallationToken(ctx context.Context, appId, installationId int64,
	privateKeyPath, cacheKey string) (string, error) {

	if Rdb != nil {
		if token, found := GetToken(ctx, cacheKey); found {
			return token, nil
		}
	}
	return "", nil
}
