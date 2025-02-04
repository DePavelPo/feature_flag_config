package main

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func main() {

	_ = context.Background()

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		logrus.Fatalf("cannot connect to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logrus.Fatalf("cannot ping db: %v", err)
	}

	logrus.Info("Pong")

}
