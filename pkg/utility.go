package pkg

import (
	 "log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load("../.env")

	if err != nil {
		 log.Fatalf("Error loading .env file")
		// panic(err)
	}

	return os.Getenv(key)
}