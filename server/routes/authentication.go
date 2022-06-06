package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"jinya-backup/server/database"
	"net/http"
)

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	loginData := new(userLogin)
	err := decoder.Decode(loginData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	apiKey, _, err := database.Login(loginData.Username, loginData.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	responseBody, err := json.Marshal(apiKey)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	authCookie := new(http.Cookie)
	authCookie.Name = "Jinya-Auth"
	authCookie.Value = apiKey.Token
	authCookie.HttpOnly = true

	http.SetCookie(w, authCookie)
	_, _ = w.Write(responseBody)
}

func UserLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authCookie, err := r.Cookie("Jinya-Auth")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	apiKey, err := database.FindApiKeyByToken(authCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_ = apiKey.Delete()
	w.WriteHeader(http.StatusNoContent)
}
