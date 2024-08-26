package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/opplieam/bb-core-api/internal/utils"
)

func main() {
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	if err := utils.SeedUsers(db); err != nil {
		log.Fatal(err)
	}
	if err := utils.SeedSellers(db); err != nil {
		log.Fatal(err)
	}
	if err := utils.SeedProducts(db); err != nil {
		log.Fatal(err)
	}
	if err := utils.SeedSubscribeProduct(db); err != nil {
		log.Fatal(err)
	}

}
