package app

import (
	"fmt"

	"github.com/combodga/Project/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(serverAddr, baseURL, dbFile string) error {
	h, err := handler.New(serverAddr, baseURL, dbFile)
	if err != nil {
		return fmt.Errorf("handler: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	e.POST("/", h.CreateURL)
	e.GET("/:id", h.RetrieveURL)
	e.POST("/api/shorten", h.CreateURLInJSON)
	e.GET("/api/user/urls", h.ListURL)

	e.Logger.Fatal(e.Start(serverAddr))

	return nil
}
