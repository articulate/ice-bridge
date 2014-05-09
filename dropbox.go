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
