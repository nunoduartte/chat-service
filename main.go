package main

import (
	"chat-service/chat"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	server := &chat.Chat{redisClient}
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("não foi possível ouvir na porta 8080 %v", err)
	}
}
