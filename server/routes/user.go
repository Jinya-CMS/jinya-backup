package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"jinya-backup/server/database"
	"net/http"
)

func isMyself(r *http.Request, userId string) bool {
	authCookie, err := r.Cookie("Jinya-Auth")
	if err != nil {
		return false
	}

	apiKey, err := database.FindApiKeyByToken(authCookie.Value)
	if err != nil {
		return false
	}

	user, err := apiKey.GetUser()
	if err != nil {
		return false
	}

	return user.Id == userId
}

func ListUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	users, err := database.FindAllUsers()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(users)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(result)
}

func GetUserById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	user, err := database.FindUserById(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	result, err := json.Marshal(user)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(result)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	user, err := database.FindUserById(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	type userPostData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	postBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postData := new(userPostData)
	err = json.Unmarshal(postBody, postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if postData.Username != "" {
		user.Name = postData.Username
	}

	if postData.Password != "" {
		err = user.SetPassword(postData.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err = user.Update()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if isMyself(r, params.ByName("id")) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	user, err := database.FindUserById(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = user.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	type userPostData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	postBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postData := new(userPostData)
	err = json.Unmarshal(postBody, postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := database.User{Name: postData.Username}
	err = user.SetPassword(postData.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
