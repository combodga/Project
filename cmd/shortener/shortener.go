package main

import (
	"flag"
	"os"

	"github.com/combodga/Project/internal/app"
)

func main() {
	var serverAddr string
	var baseURL string
	var dbFile string

	flag.StringVar(&serverAddr, "a", "", "server address")
	flag.StringVar(&baseURL, "b", "", "base URL")
	flag.StringVar(&dbFile, "f", "", "file storage path")
	flag.Parse()

	if serverAddr == "" {
		serverAddr = os.Getenv("SERVER_ADDRESS")
		if serverAddr == "" {
			serverAddr = "localhost:8080"
		}
	}

	if baseURL == "" {
		baseURL = os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://" + serverAddr
		}
	}

	if dbFile == "" {
		dbFile = os.Getenv("FILE_STORAGE_PATH")
	}

	err := app.Start(serverAddr, baseURL, dbFile)
	if err != nil {
		panic(err)
	}
}
