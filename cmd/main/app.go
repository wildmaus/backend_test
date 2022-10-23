package main

import (
	"fmt"
	"log"
	"net/http"

	"backend_test/internal/config"
	"backend_test/internal/server"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет!")
}

func main() {
	cfg := config.GetConfig()
	if err := server.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
