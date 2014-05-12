package main

import (
	"fmt"
	"sync"
	"errors"
	"github.com/stacktic/dropbox"
	"path/filepath"
	"github.com/skratchdot/open-golang/open"
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
	box, err := getBox(config)
	exitIf(err)

	files, err := box.Metadata(config.DropboxPath, true, false, "", "", 0)
	exitIf(err)

	return files.Contents, nil
}

func downloadFile(config *ConfigFile, dropboxPath string, localPath string) error {
	var err error
	var box *dropbox.Dropbox

	if box, err = getBox(config); err == nil {
		err = box.DownloadToFile(dropboxPath, localPath, "")
	}
	return err
}

var box *dropbox.Dropbox

func getBox(config *ConfigFile) (*dropbox.Dropbox, error) {

	if box == nil {
		var err error
		var token string

		if config == nil {
			panic("need config to get the dropbox instance")
		}

		box = dropbox.NewDropbox()
		box.SetAppInfo(config.ClientId, config.ClientSecret)

		if token, err = getAccessToken(config, box); err == nil {
			box.SetAccessToken(token)
		} else {
			panic("an access token is required")
		}
	}

	return box, nil
}

func getAccessToken(config *ConfigFile, box *dropbox.Dropbox) (string, error) {
	var err error

	if len(config.Token) == 0 {
		if err = authorize(box); err == nil {
			config.Token = box.AccessToken()
			config.Write(configFilename)
		}
	}

	return config.Token, err
}

func authorize(box *dropbox.Dropbox) error {
	var code, codeUrl string
	var err error

	codeUrl = box.Session.Config.AuthCodeURL("")
	fmt.Printf("A browser window should have opened at\n%s\n" +
	"Please authorize and enter the code: ", codeUrl)
	if err = open.Start(codeUrl); err == nil {
		fmt.Scanln(&code)
		_, err = box.Session.Exchange(code)
	}

	return err
}

func fixDropboxPath(dropboxPath string) string {
	if strings.Index(dropboxPath, "/") != 0 {
		return "/" + dropboxPath
	} else {
		return dropboxPath
	}
}
