package main

import (
	"fmt"
	"github.com/stacktic/dropbox"
	"path/filepath"
)

func archiveFolder(config *ConfigFile) error {

	var fileEntries []dropbox.Entry
	var err error

	fileEntries, err = getFiles(config)
	if err != nil {
		return err
	}

	for _, file := range fileEntries {
		var fileName string

		fileName = filepath.Base(file.Path)
		if file.IsDir {
			fmt.Printf("%v is a directory; skipping\n", file.Path)
			continue
		}

		fmt.Printf("downloading file: %v - rev = %v\n", fileName, file.Revision)

		err = downloadFile(config, file.Path, filepath.Join(config.LocalPath, fileName))
		if err != nil {
			fmt.Println("error downloading: " + file.Path)
			continue
		}
	}
	return err
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

		if token, err = getAccessToken(config); err == nil {
			box.SetAccessToken(token)
		} else {
			panic("an access token is required")
		}
	}

	return box, nil
}

func getAccessToken(config *ConfigFile) (string, error) {
	var err error

	if len(config.Token) == 0 {
		if err = box.Auth(); err == nil {
			config.Token = box.AccessToken()
			config.Write(configFilename)
		}
	}

	return config.Token, err
}
