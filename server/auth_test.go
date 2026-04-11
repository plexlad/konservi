package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestLogin_Success(t *testing.T) {
	e := echo.New()

	reqBody := LoginRequest{
		Username: "user",
		Password: "pass",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/token/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	LoginEndpoint(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) < 2 {
		t.Error("expected both access_token and refresh_token cookies")
	}

	accessCookie := cookies[0]
	if accessCookie.Name != "access_token" {
		t.Errorf("expected cookie name 'access_token', got '%s'", accessCookie.Name)
	}

	refreshCookie := cookies[1]
	if refreshCookie.Name != "refresh_token" {
		t.Errorf("expected cookie name 'refresh_token', got '%s'", refreshCookie.Name)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	e := echo.New()

	reqBody := LoginRequest{
		Username: "wrong",
		Password: "wrong",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/token/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	LoginEndpoint(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestLogin_MissingBody(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewReader([]byte{}),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	LoginEndpoint(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	e := echo.New()

	loginReq := LoginRequest{Username: "user", Password: "pass"}
	loginBody, _ := json.Marshal(loginReq)
	loginReqTest := httptest.NewRequest(
		http.MethodPost, "/token/login",
		bytes.NewReader(loginBody),
	)
	loginReqTest.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	loginC := e.NewContext(loginReqTest, loginRec)

	LoginEndpoint(loginC)

	refreshCookie := loginRec.Result().Cookies()[1]

	req := httptest.NewRequest(http.MethodPost, "/token/refresh", nil)
	req.AddCookie(refreshCookie)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	RefreshTokenEndpoint(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("expected new access token cookie")
	}
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	e := echo.New()

	LoginReq := LoginRequest{Username: "user", Password: "pass"}
	loginBody, _ := json.Marshal(LoginReq)
	loginReqTest := httptest.NewRequest(
		http.MethodPost,
		"/token/login",
		bytes.NewReader(loginBody),
	)
	loginReqTest.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	loginC := e.NewContext(loginReqTest, loginRec)

	LoginEndpoint(loginC)

	accessCookie := loginRec.Result().Cookies()[0]

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.AddCookie(accessCookie)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	JWTMiddleware(func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestJWTMiddleware_NoCookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	JWTMiddleware(func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}
