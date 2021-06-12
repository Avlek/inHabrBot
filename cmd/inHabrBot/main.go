package main

import "github.com/avlek/inHabrBot/internal/impl"

func main() {
	server := impl.NewServer()
	server.Run()
}
