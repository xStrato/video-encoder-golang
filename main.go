package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func main() {
	// db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// cmd.StartServer(db)

	println("The program is running...")
}
