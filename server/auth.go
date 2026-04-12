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
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func LoginEndpoint(c *echo.Context) error {
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

	// TODO Update saved username with userID instead
	accessToken, refreshToken := GenerateTokens(req.Username)

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "logged in"})
}

func RefreshTokenEndpoint(c *echo.Context) error {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			map[string]string{"error": "refresh token missing"},
		)
	}

	claims := &Claims{}
	_, err = jwt.ParseWithClaims(
		refreshCookie.Value,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtRefreshSecret, nil
		},
	)

	if err != nil || claims.Type != "refresh" {
		return c.JSON(
			http.StatusUnauthorized,
			map[string]string{"error": "invalid refresh token"},
		)
	}

	newAccessToken, _ := GenerateTokens(claims.UserID)

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "token refreshed"})
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		token, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		claims := &Claims{}
		_, err = jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || claims.Type != "access" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		if claims.UserID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token: missing user_id"})
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}

func GenerateTokens(userID string) (accessToken, refreshToken string) {
	accessTokenClaims := Claims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    instanceName,
		},
	}

	refreshTokenClaims := Claims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    instanceName,
		},
	}

	accessTokenInit := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, _ = accessTokenInit.SignedString(jwtSecret)

	refreshTokenInit := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshToken, _ = refreshTokenInit.SignedString(jwtRefreshSecret)

	return accessToken, refreshToken
}
