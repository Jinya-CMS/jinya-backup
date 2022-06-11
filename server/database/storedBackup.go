package database

import (
	"time"
)

type StoredBackup struct {
	Id         string    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	FullPath   string    `db:"full_path" json:"fullPath"`
	BackupDate time.Time `db:"backup_date" json:"backupDate"`
	JobId      string    `db:"backup_job_id" json:"jobId"`
}

func FindAllStoredBackups() ([]StoredBackup, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	storedBackups := make([]StoredBackup, 0)

	err = db.Select(&storedBackups, "SELECT id, full_path, name, backup_date, backup_job_id FROM stored_backup")

	return storedBackups, err
}

func FindStoredBackupById(id string) (*StoredBackup, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	storedBackup := new(StoredBackup)

	err = db.Get(storedBackup, "SELECT id, full_path, name, backup_date, backup_job_id FROM stored_backup WHERE id = $1", id)

	return storedBackup, err
}

func (storedBackup *StoredBackup) Create() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO stored_backup (name, full_path, backup_date, backup_job_id) VALUES ($1, $2, $3, $4)", storedBackup.Name, storedBackup.FullPath, storedBackup.BackupDate.Format(time.RFC3339), storedBackup.JobId)

	if err != nil {
		return err
	}

	return nil
}

func (storedBackup *StoredBackup) Delete() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM stored_backup WHERE id = $1", storedBackup.Id)

	return err
}

func (storedBackup *StoredBackup) GetBackupJob() (*BackupJob, error) {
	return FindBackupJobById(storedBackup.JobId)
}
