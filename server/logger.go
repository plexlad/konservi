package main

import (
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Formatter: log.JSONFormatter,
	ReportTimestamp: true,
	ReportCaller: true,
	TimeFormat: time.RFC822Z,
})

func HTTPErrorHandler(c *echo.Context, err error) {
	logger.Error("HTTP error", "error", err)
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.JSON(code, map[string]string{"error": err.Error()})
}

func EchoLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()
		err := next(c)
		duration := time.Since(start)

		status := 0
		if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
			status = resp.Status
		}

		logger.Info("Request processed",
			"method", c.Request().Method,
			"path", c.Request().RequestURI,
			"status", status,
			"duration", duration,
			"client_ip", c.RealIP(),
		)
		return err
	}
}
