package main

import (
	"log"

	"github.com/avlek/inHabrBot/internal/impl"
)

func main() {
	server := impl.NewServer()
	err := server.Run()
	if err != nil {
		log.Println(err)
	}
}
