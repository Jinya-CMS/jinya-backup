package server

import (
	"jinya-backup/server/database"
	"log"
	"net/http"
	"os"
)

func RunServer() {
	err := database.InitializeDatabase()
	if err != nil {
		log.Println("Database seems to exist already")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", http.FileServer(http.Dir("./web")))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
