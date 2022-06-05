package database

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
	"io/fs"
	"os"
)

type BackupJob struct {
	Id         string `db:"id"`
	Name       string `db:"name"`
	Host       string `db:"host"`
	Port       int    `db:"port"`
	Type       string `db:"type"`
	Username   string `db:"username"`
	password   string `db:"password"`
	RemotePath string `db:"remote_path"`
	LocalPath  string `db:"local_path"`
}

func FindAllBackupJobs() ([]BackupJob, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	backupJobs := make([]BackupJob, 0)

	err = db.Select(&backupJobs, "SELECT id, name, host, port, type, username, password, remote_path, local_path FROM backup_job")

	return backupJobs, err
}

func FindBackupJobById(id string) (*BackupJob, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	backupJob := new(BackupJob)

	err = db.Get(backupJob, "SELECT id, name, host, port, type, username, password, remote_path, local_path FROM backup_job WHERE id = $1", id)

	return backupJob, err
}

func (backupJob *BackupJob) Create() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO backup_job (name, host, port, type, username, password, remote_path, local_path) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", backupJob.Name, backupJob.Host, backupJob.Port, backupJob.Type, backupJob.Username, backupJob.password, backupJob.RemotePath, backupJob.LocalPath)

	if err != nil {
		return err
	}

	if _, err := os.Stat(backupJob.LocalPath); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(backupJob.LocalPath, 0775)
		if err != nil {
			return err
		}
	}

	return nil
}

func (backupJob *BackupJob) Update() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE backup_job SET name = $1, host = $2, port = $3, type = $4, username = $5, password = $6, remote_path = $6, local_path = $7 WHERE id = $8", backupJob.Id, backupJob.Name, backupJob.Host, backupJob.Port, backupJob.Type, backupJob.Username, backupJob.password, backupJob.RemotePath, backupJob.LocalPath)

	if _, err := os.Stat(backupJob.LocalPath); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(backupJob.LocalPath, 0775)
		if err != nil {
			return err
		}
	}

	return err
}

func (backupJob *BackupJob) Delete() error {
	db, err := ConnectToDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM backup_job WHERE id = $1", backupJob.Id)

	return err
}

func (backupJob *BackupJob) SetPassword(password string) error {
	dbSecret := os.Getenv("DB_SECRET_KEY")
	decodedSecret, err := base64.StdEncoding.DecodeString(dbSecret)
	if err != nil {
		return err
	}

	cipher, err := aes.NewCipher(decodedSecret)
	if err != nil {
		return err
	}

	encryptedBytes := make([]byte, 0)
	cipher.Encrypt(encryptedBytes, []byte(password))

	encryptedPassword := base64.StdEncoding.EncodeToString(encryptedBytes)
	backupJob.password = encryptedPassword

	return nil
}

func (backupJob *BackupJob) GetPassword() (string, error) {
	dbSecret := os.Getenv("DB_SECRET_KEY")
	decodedSecret, err := base64.StdEncoding.DecodeString(dbSecret)
	if err != nil {
		return "", err
	}

	cipher, err := aes.NewCipher(decodedSecret)
	if err != nil {
		return "", err
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(backupJob.password)
	if err != nil {
		return "", err
	}

	decryptedBytes := make([]byte, 0)
	cipher.Decrypt(decryptedBytes, decodedPassword)

	return string(decryptedBytes), nil
}

func (backupJob *BackupJob) GetStoredBackups() ([]StoredBackup, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	storedBackups := make([]StoredBackup, 0)

	err = db.Get(storedBackups, "SELECT id, full_path, name, backup_date, backup_job_id FROM stored_backup WHERE backup_job_id = $1", backupJob.Id)

	return storedBackups, err
}
