package app

import (
	"github.com/combodga/Project/internal/handler"

	"github.com/labstack/echo/v4"
)

func Start(serverAddr, baseURL string) error {
	handler.ServerAddr = serverAddr
	handler.BaseURL = baseURL

	e := echo.New()
	e.POST("/", handler.CreateURL)
	e.GET("/:id", handler.RetrieveURL)
	e.POST("/api/shorten", handler.CreateURLInJSON)

	e.Logger.Fatal(e.Start(serverAddr))

	return nil
}
