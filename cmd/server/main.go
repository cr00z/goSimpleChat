package main

import (
	"github.com/cr00z/goSimpleChat/internal/app"
)

// @title Simple Chat API
// @version 1.0.0
// @description Simple Chat API Backend (Golang)

// @host localhost:5000
// @BasePath /

// @contact.name @imcr00z
// @contact.email netrebinr@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app.Run()
}
