package main

import (
	"log"

	"tracker/internal/helpers"
	"tracker/internal/server"
)

func main() {
	if err := helpers.ChangeDirProjectRoot(); err != nil {
		log.Fatal(err)
	}
	server.Init(link, port)
}

const (
	port = ":8080"
	link = "http://localhost"
)
