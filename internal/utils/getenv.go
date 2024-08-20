package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetEnvForTesting() {
	path := ".env"
	for i := 0; i < 5; i++ {
		err := godotenv.Load(path)
		if err == nil {
			break
		}
		path = "../" + path
	}
}
