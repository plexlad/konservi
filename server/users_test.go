package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.POST("/users", CreateUserEndpoint)

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodPost, "/users", CreateUserRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	}, "")

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body)
	}
	var res UserResponse
	json.NewDecoder(rec.Body).Decode(&res)
	if res.Username != "alice" {
		t.Errorf("expected username alice, got %s", res.Username)
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.POST("/users", CreateUserEndpoint)

	seedUser(t, "alice", "alice@example.com", "password123")

	rec := makeRequest(t, e, http.MethodPost, "/users", CreateUserRequest{
		Username: "alice",
		Email:    "other@example.com",
		Password: "password123",
	}, "")

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rec.Code)
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.POST("/users", CreateUserEndpoint)

	rec := makeRequest(t, e, http.MethodPost, "/users", CreateUserRequest{
		Username: "alice",
	}, "")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetUser_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.GET("/users/:id", GetUserEndpoint)

	u := seedUser(t, "bob", "bob@example.com", "pass")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodGet, "/users/"+u.ID.String(), nil, "")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.GET("/users/:id", GetUserEndpoint)

	rec := makeRequest(t, e, http.MethodGet, "/users/00000000-0000-0000-0000-000000000000", nil, "")

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.PATCH("/users/:id", UpdateUserEndpoint)

	u := seedUser(t, "carol", "carol@example.com", "pass")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodPatch, "/users/"+u.ID.String(),
		UpdateUserRequest{Username: "carol_updated"}, u.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body)
	}
	var res UserResponse
	json.NewDecoder(rec.Body).Decode(&res)
	if res.Username != "carol_updated" {
		t.Errorf("expected carol_updated, got %s", res.Username)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.DELETE("/users/:id", DeleteUserEndpoint)

	u := seedUser(t, "dave", "dave@example.com", "pass")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodDelete, "/users/"+u.ID.String(), nil, "")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

