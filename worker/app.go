package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func notifyJinyaBackupServer(server string, id string) {
	log.Printf("Send request to get backup to %s", server)
	response, err := http.Post(fmt.Sprintf("%sapi/backup-job/%s/backup", server, id), "text/plain", nil)
	if err != nil {
		log.Println("Failed to trigger server")
		log.Println(err.Error())
		return
	}

	log.Printf("Got response code %d", response.StatusCode)
}

func runDumpJob(job Job) {
	log.Printf("Start mysql dump for job %s on server %s", job.Id, job.Host)
	cmd := exec.Command("mysqldump", "--host="+job.Host, "--password="+job.Password, "--port="+strconv.Itoa(job.Port), "--user="+job.User, job.Database)
	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = os.Stderr.WriteString("Failed to execute mysqldump\n")
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		_, _ = os.Stderr.WriteString(string(output))
		return
	}

	log.Println("Got output from command dumped below")
	log.Println(string(output))

	err = ioutil.WriteFile(job.Output, output, 0755)
	if err != nil {
		_, _ = os.Stderr.WriteString("Failed to write mysql dump\n")
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		return
	}
}

func archiveDirectory(source string, target string) error {
	log.Printf("Create target file %s", target)
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	log.Println("Create new tar writer")
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	log.Printf("Get info from source directory %s", source)
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		log.Printf("Get file info header from %s", info.Name())
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		log.Printf("Write tarball header for %s", header.Name)
		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		log.Printf("Open file %s", path)
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		log.Printf("Copy %s to tarball", file.Name())
		_, err = io.Copy(tarball, file)

		return err
	})
}

func compressFile(source string, target string) error {
	log.Printf("Compress %s using gzip", source)
	log.Printf("Open %s as reader", source)
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(target)
	log.Printf("Create target file %s", target)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	log.Printf("Create new gzip writer")
	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	log.Printf("Copy source to new gzip writer")
	_, err = io.Copy(archiver, reader)

	return err
}

func runCompressJob(job Job) {
	log.Printf("Start archive job of directory %s", job.Input)
	pwd, _ := os.Getwd()
	tmpFile := filepath.Join(pwd, uuid.New().String())

	err := archiveDirectory(job.Input, tmpFile)
	if err != nil {
		_ = os.Remove(tmpFile)
		_, _ = os.Stderr.WriteString("Failed to write tar file\n")
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	err = compressFile(tmpFile, job.Output)
	if err != nil {
		_ = os.Remove(tmpFile)
		_, _ = os.Stderr.WriteString("Failed to write gzip file\n")
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		return
	}

	_ = os.Remove(tmpFile)
}

func main() {
	log.Println("Start worker...")
	pwd, _ := os.Getwd()
	configFilePath := ""
	flag.StringVar(&configFilePath, "config-file", fmt.Sprintf("%s/jinya-backup.yaml", pwd), "Specifies the config file to use")
	flag.Parse()

	log.Printf("Using config file %s", configFilePath)
	log.Println("Parsing config file")
	configFileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to open config file %s", configFilePath)
	}

	config := ConfigStructure{}
	err = yaml.Unmarshal(configFileContent, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file %s", configFilePath)
	}

	for _, job := range config.Jobs {
		if job.Type == JobTypeMySqlDump {
			runDumpJob(job)
		} else if job.Type == JobTypeArchive {
			runCompressJob(job)
		}
		notifyJinyaBackupServer(config.Server, job.Id)
	}
}
