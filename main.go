package main

import (
	"crypto/rsa"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	App        *AppConfig
	Valkey     *redis.Client
	mux        *BotMux
	Pool       *pgxpool.Pool
	PrivateKey *rsa.PrivateKey
)

func main() {
	fmt.Println("Hello world")
}
