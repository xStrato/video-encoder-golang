package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/xStrato/video-encoder-golang/infrastructure/cmd"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	cmd.StartServer(db)
}
