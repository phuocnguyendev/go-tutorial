package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() *gorm.DB {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || pass == "" || dbName == "" {
		log.Fatal("missing DB environment variables")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		user, pass, port, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to Postgres:", err)
	}

	return db
}
