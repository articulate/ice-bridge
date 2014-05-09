package main

import (
	"fmt"
	_ "strconv"
)

func main() {
	fmt.Println("Hello from the ice-box")

	//TODO: these should maybe be command-line args
	var localPath = `C:\dev\dropbox_archive`
	var dropboxPath = "images/wallpapers"
	var err = archiveFolder(dropboxPath, localPath)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

}
