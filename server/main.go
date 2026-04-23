package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var (
	appPort          string
	frontendAddress  string
	instanceName     string
	jwtSecret        []byte
	jwtRefreshSecret []byte
	sqlDbAddress     string
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system env vars")
	}
	jwtSecret = []byte(os.Getenv("KONSERVI_JWT_SECRET"))
	jwtRefreshSecret = []byte(os.Getenv("KONSERVI_JWT_REFRESH_SECRET"))
	frontendAddress = os.Getenv("KONSERVI_FRONTEND_ADDRESS")
	instanceName = os.Getenv("KONSERVI_INSTANCE_NAME")
	appPort = os.Getenv("KONSERVI_PORT")
	sqlDbAddress = os.Getenv("KONSERVI_DATABASE_URL")
}

func isEnvValid() bool {
	return len(jwtSecret) > 0 &&
		len(jwtRefreshSecret) > 0 &&
		len(frontendAddress) > 0 &&
		len(instanceName) > 0 &&
		len(appPort) > 0 &&
		len(sqlDbAddress) > 0
}

func main() {
	if !isEnvValid() {
		logger.Fatal("Check your .env files or environment variables.",
			"KONSERVI_JWT_SECRET", len(jwtSecret),
			"KONSERVI_REFRESH_SECRET", len(jwtRefreshSecret),
			"KONSERVI_FRONTEND_ADDRESS", len(frontendAddress),
			"KONSERVI_INSTANCE_NAME", len(instanceName),
			"KONSERVI_PORT", len(appPort),
			"KONSERVI_DB_URL", len(sqlDbAddress),
		)
	}
	if len(jwtSecret) < 32 {
		logger.Warn("KONSERVI_JWT_SECRET should be at least 32 characters long.")
	}

	if err := InitDB(sqlDbAddress); err != nil {
		logger.Fatal("Failed to initialize databse", "error", err)
	}
	defer db.Close()

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

	api.POST("/users", CreateUserEndpoint)
	api.GET("/users", ListUsersEndpoint)
	api.GET("/users/:id", GetUserEndpoint)
	api.PUT("/users/:id", UpdateUserEndpoint)
	api.DELETE("/users/:id", DeleteUserEndpoint)

	api.POST("/projects", CreateProjectEndpoint)
	api.GET("/projects", ListProjectsEndpoint)
	api.GET("/projects/:id", GetProjectEndpoint)
	api.PUT("/projects/:id", UpdateProjectEndpoint)
	api.DELETE("/projects/:id", DeleteProjectEndpoint)

	api.POST("/projects/:project_id/memberships", CreateProjectMembershipEndpoint)
	api.GET("/projects/:project_id/memberships", ListProjectMembershipsEndpoint)
	api.PUT("/projects/:project_id/memberships/:user_id", UpdateProjectMembershipEndpoint)
	api.DELETE("/projects/:project_id/memberships/:user_id", DeleteProjectMembershipEndpoint)

	api.POST("/projects/:project_id/entries", CreateEntryEndpoint)
	api.GET("/projects/:project_id/entries", ListEntriesEndpoint)
	api.GET("/projects/:project_id/entries/:id", GetEntryEndpoint)
	api.PUT("/projects/:project_id/entries/:id", UpdateEntryEndpoint)
	api.DELETE("/projects/:project_id/entries/:id", DeleteEntryEndpoint)

	api.POST("/projects/:project_id/people", CreatePersonEndpoint)
	api.GET("/projects/:project_id/people", ListPeopleEndpoint)
	api.GET("/projects/:project_id/people/:id", GetPersonEndpoint)
	api.PUT("/projects/:project_id/people/:id", UpdatePersonEndpoint)
	api.DELETE("/projects/:project_id/people/:id", DeletePersonEndpoint)

	api.POST("/projects/:project_id/families", CreateFamilyEndpoint)
	api.GET("/projects/:project_id/families", ListFamiliesEndpoint)
	api.GET("/projects/:project_id/families/:id", GetFamilyEndpoint)
	api.PUT("/projects/:project_id/families/:id", UpdateFamilyEndpoint)
	api.DELETE("/projects/:project_id/families/:id", DeleteFamilyEndpoint)

	if err := e.Start(":" + appPort); err != nil {
		logger.Error("failed to start server", "error", err)
	}
}
