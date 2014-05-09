package main

import (
	"fmt"
	"github.com/stacktic/dropbox"
	"path/filepath"
)

func archiveFolder(dropboxPath string, localPath string) error {
	var files, err = getFiles(dropboxPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		var fileName = filepath.Base(file.Path)
		fmt.Println("downloading " + fileName)
		var downloadError = downloadFile(file.Path, localPath+`\`+fileName)
		if downloadError != nil {
			fmt.Println("error downloading: " + )
			continue
		}
	}
	return nil
}

func getFiles(path string) ([]dropbox.Entry, error) {
	var box = getBox()

	var files, err = box.Metadata(path, true, false, "", "", 0)

	if err != nil {
		return nil, err
	} else {
		return files.Contents, nil
	}
}

func downloadFile(dropboxPath string, localPath string) error {
	var box = getBox()
	var err = box.DownloadToFile(dropboxPath, localPath, "")
	if err != nil {
		return err
	}
}

var box *dropbox.Dropbox

func getBox() *dropbox.Dropbox {
	if box == nil {
		var clientid, clientsecret, token string

		//TODO: pull this info a file or sumthin
		clientid = "uduoyfcovd614d3"
		clientsecret = "9jruxmosec72ko1"
		token = "TOKEN" //don't push this to github.

		box = dropbox.NewDropbox()
		box.SetAppInfo(clientid, clientsecret)
		box.SetAccessToken(token)
	}

	return box
}
