package main

import (
	"fmt"
	"sync"
	"errors"
	"github.com/stacktic/dropbox"
	"path/filepath"
)

func archiveFolder(config *ConfigFile) error {
	const numDownloaders = 5

	done := make(chan struct{})
	defer close(done)

	files, errc := walkFiles(done, config)

	// spawn up a bunch of concurrent downloaders
	var wg sync.WaitGroup
	wg.Add(numDownloaders)
	c := make(chan download)
	for i := 0; i < numDownloaders; i++ {
		go func() {
			downloader(config, done, files, c)
			wg.Done()
		}()
	}

	// wait for all the downloaders to finish
	go func() {
		wg.Wait()
		close(c)
	}()

	// examine all the download results looking for an error
	for r := range c {
		if r.err != nil {
			return r.err
		}
	}

	// finally check to see if the original file walk failed
	if err := <-errc; err != nil {
		return err
	}

	return nil
}

type download struct {
	path string
	revision string
	err error
}

func walkFiles(done <-chan struct{}, config *ConfigFile) (<-chan dropbox.Entry, <-chan error) {
	files := make(chan dropbox.Entry)
	errc := make(chan error, 1)

	go func() {
		defer close(files)

		var fileEntries []dropbox.Entry
		var err error

		fileEntries, err = getFiles(config)
		if err != nil {
			errc <- err
		}

		for _, file := range fileEntries {
			if file.IsDir {
				fmt.Printf("%v is a directory; skipping\n", file.Path)
				continue
			}

			select {
			case files <- file:
			case <-done:
				errc <- errors.New("walk cancelled")
			}
		}
		errc <- nil
	}()

	return files, errc
}

func downloader(config *ConfigFile, done <-chan struct{}, files <-chan dropbox.Entry, c chan<- download) {
	for file := range files {
		fileName := filepath.Base(file.Path)
		fmt.Printf("downloading file: %v - rev = %v\n", file.Path, file.Revision)
		err := downloadFile(config, file.Path, filepath.Join(config.LocalPath, fileName))
		select {
		case c <- download{file.Path, file.Revision, err}:
		case <-done:
			return
		}
	}
}

func getFiles(config *ConfigFile) ([]dropbox.Entry, error) {
	var box = getBox(config)

	var files, err = box.Metadata(config.DropboxPath, true, false, "", "", 0)

	if err != nil {
		return nil, err
	} else {
		return files.Contents, nil
	}
}

func downloadFile(config *ConfigFile, dropboxPath string, localPath string) error {
	var box = getBox(config)
	var err = box.DownloadToFile(dropboxPath, localPath, "")
	if err != nil {
		return err
	}
	return nil
}

var box *dropbox.Dropbox

func getBox(config *ConfigFile) *dropbox.Dropbox {

	if box == nil {
		if config == nil {
			panic("need config to get the dropbox instance")
		}
		box = dropbox.NewDropbox()
		box.SetAppInfo(config.ClientId, config.ClientSecret)
		box.SetAccessToken(config.Token)
	}

	return box
}
