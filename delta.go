package main

import (
	"fmt"
	"github.com/stacktic/dropbox"
)

func getFilesToDownload(config *ConfigFile) ([]dropbox.Entry, error) {
	var delta *dropbox.DeltaPage
	var err error

	fmt.Printf("Fetching delta with cursor %v\n", config.Cursor)
	delta, err = getDelta(config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("updating local cursor value to: %v\n", delta.Cursor)
	config.Cursor = delta.Cursor
	err = config.Write(".icebridge")
	if err != nil {
		fmt.Println("non-fatal error updating cursor value. Continuing.")
	}

	if delta.Reset {
		fmt.Printf("Delta indicates working copy needs to be reset. Downloading all files from %v\n", config.DropboxPath)
		return getAllFiles(config)
	} else {
		var entries = make([]dropbox.Entry, len(delta.Entries))
		for i, deltaEntry := range delta.Entries {
			entries[i] = *deltaEntry.Entry
		}
		return entries, nil
	}
}

func getDelta(config *ConfigFile) (*dropbox.DeltaPage, error) {
	var box = getBox(config)
	var delta, err = box.Delta(config.Cursor, "/"+config.DropboxPath)
	if err != nil {
		return nil, err
	}
	return delta, nil
}

func getAllFiles(config *ConfigFile) ([]dropbox.Entry, error) {
	var box = getBox(config)

	var files, err = box.Metadata(config.DropboxPath, true, false, "", "", 0)

	if err != nil {
		return nil, err
	} else {
		return files.Contents, nil
	}
}
