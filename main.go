package main

import (
	"SongLibrary/di"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Panic("Ошибка при загрузке .env файла", err)
	}

	log.Println("INFO: Environment variables loaded")

	di.Migrate()
	di.InitAppModule()
}
