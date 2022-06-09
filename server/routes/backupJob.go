package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"jinya-backup/server/database"
	"jinya-backup/server/download"
	"log"
	"net/http"
)

type backupJobPostData struct {
	Host       string `json:"host"`
	LocalPath  string `json:"localPath"`
	RemotePath string `json:"remotePath"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Port       int    `json:"port"`
}

func ListBackupJobs(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	jobs, err := database.FindAllBackupJobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func CreateBackupJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postData := new(backupJobPostData)
	err = json.Unmarshal(body, postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	backupJob := database.BackupJob{
		Name:       postData.Name,
		Host:       postData.Host,
		Port:       postData.Port,
		Type:       "ftp",
		Username:   postData.Username,
		RemotePath: postData.RemotePath,
		LocalPath:  postData.LocalPath,
	}
	err = backupJob.SetPassword(postData.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = backupJob.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateBackupJob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	backupJob, err := database.FindBackupJobById(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postData := new(backupJobPostData)
	err = json.Unmarshal(body, postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if postData.Name != "" {
		backupJob.Name = postData.Name
	}
	if postData.Host != "" {
		backupJob.Host = postData.Host
	}
	if postData.Port != 0 {
		backupJob.Port = postData.Port
	}
	if postData.Username != "" {
		backupJob.Username = postData.Username
	}
	if postData.RemotePath != "" {
		backupJob.RemotePath = postData.RemotePath
	}
	if postData.LocalPath != "" {
		backupJob.LocalPath = postData.LocalPath
	}
	if postData.Password != "" {
		err = backupJob.SetPassword(postData.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = backupJob.Update()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteBackupJob(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	backupJob, err := database.FindBackupJobById(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = backupJob.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func TriggerDownload(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	backupJob, err := database.FindBackupJobById(id)
	if err != nil {
		log.Printf("Failed to find backup job by id %s\n", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	download.QueueJob(*backupJob)

	w.WriteHeader(http.StatusNoContent)
}
