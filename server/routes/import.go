package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"jinya-backup/server/database"
	"net/http"
)

type importData struct {
	Jobs    []database.BackupJob    `json:"backups"`
	Backups []database.StoredBackup `json:"storedBackups"`
}

func ImportDatabase(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := new(importData)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.ImportDatabase(data.Jobs, data.Backups)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
