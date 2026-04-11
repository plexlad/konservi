package main

import (
	"os"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system env vars")
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func main() {
	e := echo.New()
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Use(EchoLogger)

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)

	}
}
