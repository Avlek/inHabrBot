package main

import (
	"log"
	"os"
	"strings"

	"github.com/avlek/inHabrBot/internal/impl"
)

func main() {
	toInit := false
	fileName := "configs/dev.yaml"
	if len(os.Args) > 1 {
		for i := range os.Args {
			if os.Args[i] == "init" {
				toInit = true
			}
			if strings.HasSuffix(os.Args[i], "yaml") {
				fileName = os.Args[i]
			}
		}
	}

	server := impl.NewServer(fileName)
	err := server.Run(toInit)
	if err != nil {
		log.Println(err)
	}
}
