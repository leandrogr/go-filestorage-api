package main

import (
	"github.com/leandrogr/go-filestorage-api/server"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	server := server.NewServer()
	server.Run()
}
