package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	//"fmt"
)

func DownloadQueue(dlQueue map[string]string, maxConcurrentJobs int) {
	concurrentJobs := make(chan struct{}, maxConcurrentJobs)

	for i := 0; i < maxConcurrentJobs; i++ {
		concurrentJobs <- struct{}{}
	}

	done := make(chan bool)
	allJobsDone := make(chan bool)

	go func() {
		for i := 0; i < len(dlQueue); i++ {
			<-done
			concurrentJobs <- struct{}{}
		}
		allJobsDone <- true
	}()

	tmpDirPath, err := ioutil.TempDir("", "shopify2linesheet_dl")

	if err != nil {
		log.Fatal("Failed to create temporary download directory", err)
	} else {
		log.Println("Created temporary download directory", tmpDirPath)
	}

	for url, path := range dlQueue {
		<-concurrentJobs
		go func(url string, path string) {
			err := DownloadFile(url, path, tmpDirPath)
			if err != nil {
				//log.Println(fmt.Sprintf("ERR dl picture to %s: %s", path, err))
			} else {
				//log.Println("Downloaded picture to", path)
			}
			done <- true
		}(url, path)
	}

	<-allJobsDone

	log.Println("Removing temporary download directory", tmpDirPath)
	os.RemoveAll(tmpDirPath)
}

func SuggestFilename(url string) string {
	pathComponents := strings.Split(strings.Split(url, "?")[0], "/")
	return pathComponents[len(pathComponents)-1]
}

func DownloadFile(fromUrl string, toPath string, tmpDirPath string) (err error) {

	fileinfo, err := os.Stat(toPath)
	_ = fileinfo

	if !os.IsNotExist(err) {
		return errors.New("File exists")
	}

	tmpFile, err := ioutil.TempFile(tmpDirPath, "")
	tmpFilePath := tmpFile.Name()

	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(fromUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	os.Rename(tmpFilePath, toPath)

	return nil
}
