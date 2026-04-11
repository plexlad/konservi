package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var jwtSecret []byte
var jwtRefreshSecret []byte
var frontendAddress string
var instanceName string
var appPort string

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system env vars")
	}
	jwtSecret = []byte(os.Getenv("KONSERVI_JWT_SECRET"))
	jwtRefreshSecret = []byte(os.Getenv("KONSERVI_JWT_REFRESH_SECRET"))
	frontendAddress = os.Getenv("KONSERVI_FRONTEND_ADDRESS")
	instanceName = os.Getenv("KONSERVI_INSTANCE_NAME")
	appPort = os.Getenv("KONSERVI_PORT")
}

func isEnvValid() bool {
	return len(jwtSecret) > 0 &&
				 len(jwtRefreshSecret) > 0 &&
				 len(frontendAddress) > 0 &&
				 len(instanceName) > 0 &&
				 len(appPort) > 0
}

func main() {
	if !isEnvValid() {
		logger.Fatal("Check your .env files or environment variables.",
            "jwtSecret", len(jwtSecret),
            "refreshSecret", len(jwtRefreshSecret),
            "frontendAddress", len(frontendAddress),
            "instanceName", len(instanceName),
        )
	}
	if len(jwtSecret) < 32 {
		logger.Warn("KONSERVI_JWT_SECRET should be at least 32 characters long.")
	}

	e := echo.New()
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Use(middleware.Recover())
	e.Use(EchoLogger)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{frontendAddress}, // Allow frontend
		AllowCredentials: true,
	}))

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Status 200 OK")
	})

	token := e.Group("/token")
	token.POST("/login", LoginEndpoint)
	token.POST("/refresh", RefreshTokenEndpoint)

	api := e.Group("/api")
	api.Use(JWTMiddleware)
	api.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Status 200 OK")
	})

	if err := e.Start(":" + appPort); err != nil {
		logger.Error("failed to start server", "error", err)
	}
}
