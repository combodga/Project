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
	var dbCredentials string

	flag.StringVar(&serverAddr, "a", os.Getenv("SERVER_ADDRESS"), "server address")
	flag.StringVar(&baseURL, "b", os.Getenv("BASE_URL"), "base URL")
	flag.StringVar(&dbFile, "f", os.Getenv("FILE_STORAGE_PATH"), "file storage path")
	flag.StringVar(&dbCredentials, "d", os.Getenv("DATABASE_DSN"), "database credentials")
	flag.Parse()

	if serverAddr == "" {
		serverAddr = "localhost:8080"
	}

	if baseURL == "" {
		baseURL = "http://" + serverAddr
	}

	err := app.Start(serverAddr, baseURL, dbFile, dbCredentials)
	if err != nil {
		panic(err)
	}
}
