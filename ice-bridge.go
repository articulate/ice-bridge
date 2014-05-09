package main

import (
	"fmt"
	_ "strconv"
)

func main() {
	fmt.Println("Hello from the ice-box")

	var config ConfigFile
	_ = config.Read(".icebridge")

	var err = archiveFolder(&config)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		panic(err)
		return
	}
}
