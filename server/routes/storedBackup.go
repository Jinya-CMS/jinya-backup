package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"jinya-backup/server/database"
	"net/http"
	"os"
)

func GetStoredBackupsByJob(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	job, err := database.FindBackupJobById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	backups, err := job.GetStoredBackups()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(backups)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func DownloadStoredBackup(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	id := params.ByName("backupId")
	backup, err := database.FindStoredBackupById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := os.Open(backup.FullPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+backup.Name+"\"")
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func DeleteStoredBackup(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	backup, err := database.FindStoredBackupById(params.ByName("backupId"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = backup.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
