package main

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	App    *AppConfig
	Valkey *redis.Client
	mux    *BotMux
	Pool   *pgxpool.Pool
)

func main() {
	fmt.Println("Hello world")
}
