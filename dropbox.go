package main

import (
	"fmt"
	"github.com/stacktic/dropbox"
	"path/filepath"
)

func archiveFolder(config *ConfigFile) error {
	var files, err = getFiles(config)
	if err != nil {
		return err
	}

	for _, file := range files {
		var fileName = filepath.Base(file.Path)
		fmt.Println("downloading " + fileName)
		var downloadError = downloadFile(config, file.Path, filepath.Join(config.LocalPath, fileName))
		if downloadError != nil {
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
		fmt.Println(config.DropboxPath)
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
		box = dropbox.NewDropbox()
		box.SetAppInfo(config.ClientId, config.ClientSecret)
		box.SetAccessToken(config.Token)
	}

	return box
}
