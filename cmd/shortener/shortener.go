package main

import (
	"github.com/combodga/Project/internal/app"
)

func main() {
	err := app.Start("localhost", "8080")
	if err != nil {
		panic(err)
	}
}
