package server

import (
	"github.com/julienschmidt/httprouter"
	"jinya-backup/helper"
	"jinya-backup/server/database"
	"jinya-backup/server/download"
	"jinya-backup/server/routes"
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

	router := httprouter.New()

	router.POST("/api/login", routes.UserLogin)
	router.DELETE("/api/login", routes.AuthenticatedMiddleware(routes.UserLogout))
	router.HEAD("/api/login", routes.AuthenticatedMiddleware(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(http.StatusNoContent)
	}))

	router.GET("/api/user", routes.AuthenticatedMiddleware(routes.ListUsers))
	router.POST("/api/user", routes.AuthenticatedMiddleware(routes.CreateUser))
	router.GET("/api/user/:id", routes.AuthenticatedMiddleware(routes.GetUserById))
	router.PUT("/api/user/:id", routes.AuthenticatedMiddleware(routes.UpdateUser))
	router.DELETE("/api/user/:id", routes.AuthenticatedMiddleware(routes.DeleteUser))

	router.GET("/api/backup-job", routes.AuthenticatedMiddleware(routes.ListBackupJobs))
	router.POST("/api/backup-job", routes.AuthenticatedMiddleware(routes.CreateBackupJob))
	router.PUT("/api/backup-job/:id", routes.AuthenticatedMiddleware(routes.UpdateBackupJob))
	router.DELETE("/api/backup-job/:id", routes.AuthenticatedMiddleware(routes.DeleteBackupJob))

	router.GET("/api/backup-job/:id/backup", routes.AuthenticatedMiddleware(routes.GetStoredBackupsByJob))
	router.POST("/api/backup-job/:id/backup", routes.TriggerDownload)
	router.GET("/api/backup-job/:id/backup/:backupId", routes.AuthenticatedMiddleware(routes.DownloadStoredBackup))
	router.DELETE("/api/backup-job/:id/backup/:backupId", routes.AuthenticatedMiddleware(routes.DeleteStoredBackup))

	router.NotFound = http.FileServer(http.Dir("./web"))

	for i := 0; i < helper.CpuCount; i++ {
		go download.StartJobWorker()
	}

	log.Printf("Serving on localhost:%s", port)
	log.Printf("Serving on 0.0.0.0:%s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
