package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginRequest) valid() bool {
	return r.Username != "" && r.Password != ""
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func LoginEndpoint(c *echo.Context) error {
	return Login(c, jwtSecret)
}

func Login(c *echo.Context, secret []byte) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil || !req.valid() {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}

	// TODO Database user and password hashing here
	if req.Username != "user" || req.Password != "pass" {
		return c.JSON(
			http.StatusUnauthorized,
			map[string]string{"error": "invalid credentials"},
		)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token generataion failed"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "logged in"})
}
