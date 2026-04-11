package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system env vars")
	}
}

var jwtSecret = []byte(os.Getenv("KONSERVI_JWT_SECRET"))
var frontendAddress = []byte(os.Getenv("KONSERVI_FRONTEND_URL"))

func main() {
	if len(jwtSecret) == 0 || len(frontendAddress) == 0 {
		logger.Fatal("KONSERVI_JWT_SECRET and KONSERVI_FRONTEND_URL are necessary.")
	}
	if len(jwtSecret) < 32 {
		logger.Warn("KONSERVI_JWT_SECRET should be at least 32 characters long.")
	}

	e := echo.New()
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Use(EchoLogger)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // Allow frontend
		AllowCredentials: true,
	}))

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/login", LoginEndpoint)

	if err := e.Start(":8080"); err != nil {
		logger.Error("failed to start server", "error", err)
	}
}
