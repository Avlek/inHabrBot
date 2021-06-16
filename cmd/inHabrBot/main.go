package main

import (
	"log"
	"os"

	"github.com/avlek/inHabrBot/internal/impl"
)

func main() {
	toInit := false
	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			toInit = true
		}
	}
	server := impl.NewServer()
	err := server.Run(toInit)
	if err != nil {
		log.Println(err)
	}
}
