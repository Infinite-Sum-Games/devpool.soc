package main

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	App    *AppConfig
	Valkey *redis.Client
	bs     *BotMux
)

func main() {
	fmt.Println("Hello world")
}
