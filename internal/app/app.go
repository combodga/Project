package app

import (
	"github.com/combodga/Project/internal/handler"
	"github.com/combodga/Project/internal/storage"

	"github.com/labstack/echo/v4"
)

func Start(serverAddr, baseURL, dbFile string) error {
	handler.ServerAddr = serverAddr
	handler.BaseURL = baseURL
	storage.DBFile = dbFile

	e := echo.New()
	e.POST("/", handler.CreateURL)
	e.GET("/:id", handler.RetrieveURL)
	e.POST("/api/shorten", handler.CreateURLInJSON)

	e.Logger.Fatal(e.Start(serverAddr))

	return nil
}
