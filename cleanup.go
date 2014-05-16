package main

import (
	"fmt"
	"github.com/stacktic/dropbox"
	"path/filepath"
	"strings"
	"time"
)

func cleanupOldFiles(config *ConfigFile) error {

	fmt.Println("Attempting to clean up old files")

	if config.MaxFileAge <= 0 {
		fmt.Println("MaxFileAge not configured; skipping cleanup")
		return nil
	}

	var files, err = getAllFiles(config)
	if err != nil {
		return err
	}

	for _, file := range files {
		var daysOld = getFileAgeInDays(config, &file)
		if daysOld > config.MaxFileAge {
			fmt.Printf("File '%v' is %v days old; deleting...\n", file.Path, daysOld)
			//TODO: only delete after double checking that we have a copy somewhere else
			var dErr = deleteOldFile(config, &file)
			if dErr != nil {
				fmt.Println("non-fatal error deleting a file. continuing")
				continue
			}
		}
	}
	return nil
}

func deleteOldFile(config *ConfigFile, file *dropbox.Entry) error {
	var box, boxErr = getBox(config)
	exitIf(boxErr)

	var _, err = box.Delete(file.Path)
	return err
}

func getFileAgeInDays(config *ConfigFile, file *dropbox.Entry) int {
	var fileModified time.Time
	var parseErr error

	if config.FileAgeMethod == "FileNameDate" {
		var fileNameWithExt = filepath.Base(file.Path)
		var fileName = fileNameWithExt[0:strings.LastIndex(fileNameWithExt, ".")]
		fileModified, parseErr = time.Parse("2006-01-02 15.04.05", fileName)
		exitIf(parseErr)
	} else {
		fileModified = time.Time(file.Modified)
	}

	return int(time.Now().Sub(fileModified).Hours() / float64(24))
}
