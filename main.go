package main

import (
	"fmt"

	"github.com/google/go-github/v62/github"
	"github.com/redis/go-redis/v9"
)

var (
	App    *AppConfig
	Valkey *redis.Client
	gh     *github.Client
)

func main() {
	fmt.Println("Hello world")
}
