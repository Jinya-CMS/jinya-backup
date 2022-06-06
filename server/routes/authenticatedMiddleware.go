package routes

import (
	"github.com/julienschmidt/httprouter"
	"jinya-backup/server/database"
	"net/http"
)

func AuthenticatedMiddleware(action func(http.ResponseWriter, *http.Request, httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		authCookie, err := r.Cookie("Jinya-Auth")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err = database.FindApiKeyByToken(authCookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		action(w, r, params)
	}
}
