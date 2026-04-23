package main

import (
	"net/http"
	"testing"
)

func TestCreateEntry_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	e.POST("/projects/:project_id/entries", CreateEntryEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "Journal")

	rec := makeRequest(t, e, http.MethodPost, "/projects/"+p.ID.String()+"/entries",
		CreateEntryRequest{Content: "Hello world"}, u.ID.String())

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rec.Code, rec.Body)
	}
}

func TestCreateEntry_GhostProject(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	e.POST("/projects/:project_id/entries", CreateEntryEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")

	rec := makeRequest(t, e, http.MethodPost, "/projects/00000000-0000-0000-0000-000000000000/entries",
		CreateEntryRequest{Content: "Ghost entry"}, u.ID.String())

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestCreateEntry_MissingContent(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route on e
	e.POST("/projects/:project_id/entries", CreateEntryEndpoint)

	// 2. Seed data
	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "Journal")

	// 3. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodPost, "/projects/"+p.ID.String()+"/entries",
		CreateEntryRequest{}, u.ID.String())

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetEntry_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.GET("/entries/:id", GetEntryEndpoint)

	// 2. Seed data
	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "Journal")
	en := seedEntry(t, p.ID.String(), u.ID.String(), "Hello world")

	// 3. Call makeRequest (no handler, no pathParams)
	rec := makeRequest(t, e, http.MethodGet, "/entries/"+en.ID.String(), nil, u.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body)
	}
}

func TestUpdateEntry_Forbidden(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	// 1. Register route
	e.PATCH("/entries/:id", UpdateEntryEndpoint)

	// 2. Seed data
	alice := seedUser(t, "alice", "alice@example.com", "pass")
	bob := seedUser(t, "bob", "bob@example.com", "pass")
	p := seedProject(t, alice.ID.String(), "Journal")
	en := seedEntry(t, p.ID.String(), alice.ID.String(), "Alice's entry")

	// 3. Call makeRequest
	rec := makeRequest(t, e, http.MethodPatch, "/entries/"+en.ID.String(),
		UpdateEntryRequest{Content: "Bob was here"}, bob.ID.String())

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestDeleteEntry_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	e.DELETE("/entries/:id", DeleteEntryEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "Journal")
	en := seedEntry(t, p.ID.String(), u.ID.String(), "To delete")

	rec := makeRequest(t, e, http.MethodDelete, "/entries/"+en.ID.String(), nil, u.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestListEntries_Success(t *testing.T) {
	newTestDB(t)
	e := newEcho()

	e.GET("/projects/:project_id/entries", ListEntriesEndpoint)

	u := seedUser(t, "alice", "alice@example.com", "pass")
	p := seedProject(t, u.ID.String(), "Journal")
	seedEntry(t, p.ID.String(), u.ID.String(), "Entry 1")
	seedEntry(t, p.ID.String(), u.ID.String(), "Entry 2")

	rec := makeRequest(t, e, http.MethodGet, "/projects/"+p.ID.String()+"/entries", nil, u.ID.String())

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
