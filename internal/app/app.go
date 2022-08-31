package app

import (
	"github.com/combodga/Project/internal/handler"

	"github.com/labstack/echo/v4"
)

func Start(host, port string) error {
	handler.Host = host
	handler.Port = port

	e := echo.New()
	e.POST("/", handler.CreateURL)
	e.GET("/:id", handler.RetrieveURL)

	e.Logger.Fatal(e.Start(host + ":" + port))

	return nil
}
