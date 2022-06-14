package download

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"jinya-backup/server/database"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var downloadChannel = make(chan database.BackupJob)

func QueueJob(job database.BackupJob) {
	downloadChannel <- job
}

func downloadFromFtpServer(job database.BackupJob) {
	log.Printf("Job %s: Downloading from ftp for backup job %s", job.Id, job.Name)
	log.Printf("Job %s: Decrypting password", job.Id)
	password, err := job.GetPassword()
	if err != nil {
		log.Printf("Job %s: Failed to decrypt password %s", job.Id, err.Error())
		return
	}

	log.Printf("Job %s: Dial connection to %s:%d", job.Id, job.Host, job.Port)
	connection, err := ftp.Dial(job.Host + ":" + strconv.Itoa(job.Port))
	if err != nil {
		log.Printf("Job %s: Failed to connect to server %s", job.Id, err.Error())
		return
	}
	defer connection.Quit()

	log.Printf("Job %s: Login to ftp server %s:%d", job.Id, job.Host, job.Port)
	err = connection.Login(job.Username, password)
	if err != nil {
		log.Printf("Job %s: Failed to login to sever %s", job.Id, err.Error())
		return
	}

	log.Printf("Job %s: Set connection type to binary", job.Id)
	err = connection.Type(ftp.TransferTypeBinary)
	if err != nil {
		log.Printf("Job %s: Failed to set tranfer type to binary %s", job.Id, err.Error())
		return
	}

	log.Printf("Job %s: Retrieve file %s", job.Id, job.RemotePath)
	response, err := connection.Retr(job.RemotePath)
	if err != nil {
		log.Printf("Job %s: Failed to download file from %s", job.Id, err.Error())
		return
	}
	defer response.Close()

	saveRemoteFileToLocalFile(job, err, response)
}

func downloadFromSftpServer(job database.BackupJob) {
	log.Printf("Job %s: Downloading from ftp for backup job %s", job.Id, job.Name)
	log.Printf("Job %s: Decrypting password", job.Id)
	password, err := job.GetPassword()
	if err != nil {
		log.Printf("Job %s: Failed to decrypt password %s", job.Id, err.Error())
		return
	}

	log.Printf("Job %s: Dial SSH connection to %s:%d", job.Id, job.Host, job.Port)
	connection, err := ssh.Dial("tcp", job.Host+":"+strconv.Itoa(job.Port), &ssh.ClientConfig{
		User: job.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Printf("Job %s: Failed to dial connection to %s:%d", job.Id, job.Host, job.Port)
		return
	}

	defer connection.Close()

	log.Printf("Job %s: Create new SFTP client", job.Id)
	client, err := sftp.NewClient(connection)
	if err != nil {
		log.Printf("Job %s: Failed to create SFTP client on SSH connection %s", job.Id, err.Error())
		return
	}

	log.Printf("Job %s: Create target file on local server", job.Id)
	response, err := client.Open(job.RemotePath)
	if err != nil {
		log.Printf("Job %s: Failed to open file on remote server %s:%d %s", job.Id, job.Host, job.Port, err.Error())
		return
	}

	defer response.Close()
	saveRemoteFileToLocalFile(job, err, response)
}

func saveRemoteFileToLocalFile(job database.BackupJob, err error, response io.Reader) {
	log.Printf("Job %s: Read all content from response", job.Id)
	content, err := ioutil.ReadAll(response)
	if err != nil {
		log.Printf("Job %s: Failed read downloaded file %s", job.Id, err.Error())
		return
	}

	log.Printf("Job %s: Hash file content", job.Id)
	hash := sha512.New()
	hash.Write(content)
	sum := hash.Sum(nil)

	path := filepath.Join(job.LocalPath, hex.EncodeToString(sum))

	log.Printf("Job %s: Write content to file %s", job.Id, path)
	err = ioutil.WriteFile(path, content, 0777)
	if err != nil {
		log.Printf("Job %s: Failed write downloaded file %s", job.Id, err.Error())
		return
	}

	backup := database.StoredBackup{
		Name:       job.Name,
		FullPath:   path,
		BackupDate: time.Now(),
		JobId:      job.Id,
	}

	log.Printf("Job %s: Save stored backup to database", job.Id)
	err = backup.Create()
	if err != nil {
		log.Printf("Job %s: Failed to store backup file %s", job.Id, err.Error())
		_ = os.Remove(path)
	}
}

func StartJobWorker() {
	for job := range downloadChannel {
		if job.Type == "ftp" {
			downloadFromFtpServer(job)
		} else if job.Type == "sftp" {
			downloadFromSftpServer(job)
		}
	}
}
