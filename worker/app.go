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
	"jinya-backup/runner"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func notifyJinyaBackupServer(server string, id string, logf func(message string, data ...interface{})) {
	logf("Send request to get backup to %s", server)
	response, err := http.Post(fmt.Sprintf("%sapi/backup-job/%s/backup", server, id), "text/plain", nil)
	if err != nil {
		logf("Failed to trigger server")
		logf(err.Error())
		return
	}

	logf("Got response code %d", response.StatusCode)
}

func runDumpJob(job Job, logf func(message string, data ...interface{})) {
	logf("Start mysql dump for job %s on server %s", job.Id, job.Host)
	cmd := exec.Command("mysqldump", "--host="+job.Host, "--password="+job.Password, "--port="+strconv.Itoa(job.Port), "--user="+job.User, job.Database)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logf("Failed to execute mysqldump\n")
		logf(err.Error() + "\n")
		logf(string(output))
		return
	}

	logf("Got output from command dumped below")
	logf(string(output))

	err = ioutil.WriteFile(job.Output, output, 0755)
	if err != nil {
		logf("Failed to write mysql dump\n")
		logf(err.Error() + "\n")
		return
	}
}

func archiveDirectory(source string, target string, logf func(message string, data ...interface{})) error {
	logf("Create target file %s", target)
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	logf("Create new tar writer")
	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	logf("Get info from source directory %s", source)
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

		logf("Get file info header from %s", info.Name())
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

		logf("Write tarball header for %s", header.Name)
		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		logf("Open file %s", path)
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		logf("Copy %s to tarball", file.Name())
		_, err = io.Copy(tarball, file)

		return err
	})
}

func compressFile(source string, target string, logf func(message string, data ...interface{})) error {
	logf("Compress %s using gzip", source)
	logf("Open %s as reader", source)
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(target)
	logf("Create target file %s", target)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	logf("Create new gzip writer")
	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	logf("Copy source to new gzip writer")
	_, err = io.Copy(archiver, reader)

	return err
}

func runCompressJob(job Job, logf func(message string, data ...interface{})) {
	logf("Start archive job of directory %s", job.Input)
	pwd, _ := os.Getwd()
	tmpFile := filepath.Join(pwd, uuid.New().String())

	err := archiveDirectory(job.Input, tmpFile, logf)
	if err != nil {
		_ = os.Remove(tmpFile)
		logf("Failed to write tar file\n")
		logf(err.Error() + "\n")
		return
	}

	err = compressFile(tmpFile, job.Output, logf)
	if err != nil {
		_ = os.Remove(tmpFile)
		logf("Failed to write gzip file\n")
		logf(err.Error() + "\n")
		return
	}

	_ = os.Remove(tmpFile)
}

func processor(wg *sync.WaitGroup, jobChan chan Job, logChan chan string, id int, server string) {
	logf := func(message string, data ...interface{}) {
		logChan <- fmt.Sprintf("CPU "+strconv.Itoa(id)+": "+message, data...)
	}
	for job := range jobChan {
		if job.Type == JobTypeMySqlDump {
			runDumpJob(job, logf)
		} else if job.Type == JobTypeArchive {
			runCompressJob(job, logf)
		}
		notifyJinyaBackupServer(server, job.Id, logf)
	}
	wg.Done()
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

	var (
		jobChan = make(chan Job, len(config.Jobs))
		logChan = make(chan string)
		wg      = &sync.WaitGroup{}
	)
	wg.Add(runner.CpuCount)

	go func(logChan chan string) {
		for logEntry := range logChan {
			log.Println(logEntry)
		}
	}(logChan)

	for i := 0; i < runner.CpuCount; i++ {
		go processor(wg, jobChan, logChan, i, config.Server)
	}

	for _, job := range config.Jobs {
		jobChan <- job
	}

	close(jobChan)

	wg.Wait()
	close(logChan)
}
