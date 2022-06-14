package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

type BackupJob struct {
	Id         string `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Host       string `db:"host" json:"host"`
	Port       int    `db:"port" json:"port"`
	Type       string `db:"type" json:"type"`
	Username   string `db:"username" json:"username"`
	Password   string `db:"password" json:"-"`
	RemotePath string `db:"remote_path" json:"remotePath"`
	LocalPath  string `db:"local_path" json:"localPath"`
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

	_, err = db.Exec("INSERT INTO backup_job (name, host, port, type, username, password, remote_path, local_path) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", backupJob.Name, backupJob.Host, backupJob.Port, backupJob.Type, backupJob.Username, backupJob.Password, backupJob.RemotePath, backupJob.LocalPath)

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

	_, err = db.Exec("UPDATE backup_job SET name = $1, host = $2, port = $3, type = $4, username = $5, password = $6, remote_path = $7, local_path = $8 WHERE id = $9", backupJob.Name, backupJob.Host, backupJob.Port, backupJob.Type, backupJob.Username, backupJob.Password, backupJob.RemotePath, backupJob.LocalPath, backupJob.Id)

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
	decodedSecret, err := base64.StdEncoding.DecodeString(os.Getenv("DB_SECRET_KEY"))
	if err != nil {
		return err
	}

	cphr, err := aes.NewCipher(decodedSecret)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	out := gcm.Seal(nonce, nonce, []byte(password), nil)

	encryptedPassword := base64.StdEncoding.EncodeToString(out)
	backupJob.Password = encryptedPassword

	return nil
}

func (backupJob *BackupJob) decryptPasswordAES() (string, error) {
	dbSecret := os.Getenv("DB_SECRET_KEY")
	decodedSecret, err := base64.StdEncoding.DecodeString(dbSecret)
	if err != nil {
		return "", err
	}

	cphr, err := aes.NewCipher(decodedSecret)
	if err != nil {
		return "", err
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(backupJob.Password)
	if err != nil {
		return "", err
	}

	decryptedBytes := make([]byte, len(decodedPassword))
	cphr.Decrypt(decryptedBytes, decodedPassword)

	return string(decryptedBytes), nil
}

func (backupJob *BackupJob) GetPassword() (string, error) {
	dbSecret := os.Getenv("DB_SECRET_KEY")
	decodedSecret, err := base64.StdEncoding.DecodeString(dbSecret)
	if err != nil {
		return "", err
	}

	cphr, err := aes.NewCipher(decodedSecret)
	if err != nil {
		return "", err
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(backupJob.Password)
	if err != nil {
		return "", err
	}

	gcmDecrypt, err := cipher.NewGCM(cphr)
	if err != nil {
		fmt.Println(err)
	}
	nonceSize := gcmDecrypt.NonceSize()
	if len(decodedPassword) < nonceSize {
		fmt.Println(err)
	}
	nonce, encryptedMessage := decodedPassword[:nonceSize], decodedPassword[nonceSize:]
	decryptedBytes, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return backupJob.decryptPasswordAES()
	}

	return string(decryptedBytes), nil
}

func (backupJob *BackupJob) GetStoredBackups() ([]StoredBackup, error) {
	db, err := ConnectToDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	storedBackups := make([]StoredBackup, 0)

	err = db.Select(&storedBackups, "SELECT id, full_path, name, backup_date, backup_job_id FROM stored_backup WHERE backup_job_id = $1", backupJob.Id)

	return storedBackups, err
}
