package app

import (
	"github.com/combodga/Project/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(serverAddr, baseURL, dbFile string) error {
	h := handler.New(serverAddr, baseURL, dbFile)

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	e.POST("/", h.CreateURL)
	e.GET("/:id", h.RetrieveURL)
	e.POST("/api/shorten", h.CreateURLInJSON)

	e.Logger.Fatal(e.Start(serverAddr))

	return nil
}
