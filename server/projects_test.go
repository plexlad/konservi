package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateProject_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route FIRST
	e.POST("/projects", CreateProjectEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodPost, "/projects",
		CreateProjectRequest{Name: "My Project", Description: "A test project"},
		u.ID.String())

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body)
	}
	var res ProjectResponse
	json.NewDecoder(rec.Body).Decode(&res)
	if res.Name != "My Project" {
		t.Errorf("expected 'My Project', got %s", res.Name)
	}
	if res.OwnerID != u.ID.String() {
		t.Errorf("expected ownerID %s, got %s", u.ID, res.OwnerID)
	}
}

func TestCreateProject_MissingName(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route FIRST on e
	e.POST("/projects", CreateProjectEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")

	// 2. Call makeRequest (no handler argument, no pathParams)
	rec := makeRequest(t, e, http.MethodPost, "/projects",
		CreateProjectRequest{}, u.ID.String())

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetProject_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.GET("/projects/:id", GetProjectEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "Test Project")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodGet, "/projects/"+p.ID.String(), nil, u.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body)
	}
}

func TestUpdateProject_Forbidden(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.PATCH("/projects/:id", UpdateProjectEndpoint)

	owner := seedUser(t, "owner", "owner@example.com", "pass")
	other := seedUser(t, "other", "other@example.com", "pass")
	p := seedProject(t, owner.ID.String(), "Owner Project")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodPatch, "/projects/"+p.ID.String(),
		UpdateProjectRequest{Name: "Hijacked"}, other.ID.String())

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestDeleteProject_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.DELETE("/projects/:id", DeleteProjectEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "To Delete")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodDelete, "/projects/"+p.ID.String(), nil, u.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestListProjects_OnlyOwned(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.GET("/projects", ListProjectsEndpoint)

	alice := seedUser(t, "alice", "alice@example.com", "pass")
	bob := seedUser(t, "bob", "bob@example.com", "pass")
	seedProject(t, alice.ID.String(), "Alice Project 1")
	seedProject(t, alice.ID.String(), "Alice Project 2")
	seedProject(t, bob.ID.String(), "Bob Project")

	// 2. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodGet, "/projects", nil, alice.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body)
	}
}
