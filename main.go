package main

import (
	"log"
	"os"
	"time"

	"github.com/daugminas/kyc/app/adapter"
	"github.com/daugminas/kyc/app/delivery"
	"github.com/daugminas/kyc/lib/db"
	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
)

func main() {

	// env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// db connection
	dbObj := db.New(&db.Options{
		MongoURI:     os.Getenv("MONGO_URI"),
		DatabaseName: os.Getenv("MONGO_DB"),
		DialTimeout:  10 * time.Second,
	})

	// user adapter
	userAdapter := adapter.NewUserAdapter(dbObj, os.Getenv("USERS_COLLECTION"))

	// server init, start
	e := echo.New()
	server := delivery.NewServer(e, userAdapter)
	server.RegisterRouter()
	server.Start(os.Getenv("SERVER_URI"))

}
