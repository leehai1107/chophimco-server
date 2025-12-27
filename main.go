package main

import (
	"log"

	"github.com/leehai1107/chophimco-server/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Printf("error while execute: %s", err.Error())
	}
}
