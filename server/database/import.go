package database

import (
	"context"
	"time"
)

func ImportDatabase(jobs []BackupJob, backups []StoredBackup) error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		_, err = tx.Exec("INSERT INTO backup_job (id, name, host, port, type, username, password, remote_path, local_path) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", job.Id, job.Name, job.Host, job.Port, job.Type, job.Username, job.Password, job.RemotePath, job.LocalPath)
		if err != nil {
			err = tx.Rollback()
			return err
		}
	}

	for _, backup := range backups {
		_, err = tx.Exec("INSERT INTO stored_backup (id, name, full_path, backup_date, backup_job_id) VALUES ($1, $2, $3, $4, $5)", backup.Id, backup.Name, backup.FullPath, backup.BackupDate.Format(time.RFC3339), backup.JobId)
		if err != nil {
			err = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
