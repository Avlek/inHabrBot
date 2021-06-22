package main

import (
	"log"
	"os"

	"github.com/avlek/inHabrBot/internal/impl"
)

func main() {
	toInit := false
	fileName := "configs/dev.yaml"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	if len(os.Args) > 2 {
		if os.Args[2] == "init" {
			toInit = true
		}
	}
	server := impl.NewServer(fileName)
	err := server.Run(toInit)
	if err != nil {
		log.Println(err)
	}
}
