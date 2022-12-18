package main

import (
	"log"

	"github.com/therealnemethmate/wordlesolver-cli/cmd"
)

func main() {
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
