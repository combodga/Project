package main

import (
	"os"

	"github.com/combodga/Project/internal/app"
)

func main() {
	serverAddr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		serverAddr = "localhost:8080"
	}

	baseURL, ok := os.LookupEnv("BASE_URL")
	if !ok {
		baseURL = "http://" + serverAddr
	}

	err := app.Start(serverAddr, baseURL)
	if err != nil {
		panic(err)
	}
}
