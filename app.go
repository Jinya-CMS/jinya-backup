package main

import (
	"jinya-backup/server"
	"jinya-backup/worker"
	"os"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "serve" {
		server.RunServer()
	} else {
		worker.RunBackupWorker()
	}
}
