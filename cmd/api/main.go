package main

import (
	"github.com/Burakbgmk/go-tbc-bot/internal/server"
)

func main() {

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
