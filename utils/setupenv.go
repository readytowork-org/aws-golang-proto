package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	fmt.Println(".env file loaded successfully")
	if err != nil {
		log.Fatal("Failure loading .env file")
	}
}

func GetEnvWithKey(key string) string {
	value := os.Getenv(key)
	return value
}