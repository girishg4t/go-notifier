package main

import (
	"log"

	c "github.com/girishg4t/go-notifier/consumer"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	c.StartConsumer()
}
