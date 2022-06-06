package server

import (
	"github.com/julienschmidt/httprouter"
	"jinya-backup/server/database"
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

	router.NotFound = http.FileServer(http.Dir("./web"))

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
