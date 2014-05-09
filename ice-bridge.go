package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	fmt.Println("Hello from the ice-box")

	var config ConfigFile
	err = config.Read(".icebridge")
	exitIf(err)

	err = archiveFolder(&config)
	exitIf(err)
}

func exitIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
