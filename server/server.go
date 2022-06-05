package server

import (
	"jinya-backup/server/database"
	"log"
)

func RunServer() {
	err := database.InitializeDatabase()
	log.Println(err)
}
