package main

import (
	"os"
)

func main() {

	initConfig()

	puppetRequest := readcsr(os.Stdin)
	// fmt.Printf("%+v\n", puppetRequest)
	if puppetRequest.validate() {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
