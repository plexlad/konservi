package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/plexlad/konservi/ent"
	"github.com/plexlad/konservi/ent/enttest"
	"golang.org/x/crypto/bcrypt"
)

func newTestDB(t *testing.T) *ent.Client {
	t.Helper()
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	db = client
	t.Cleanup(func() { client.Close() })
	return client
}

func newEcho() *echo.Echo {
	e := echo.New()

	// Register middleware once
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID := c.Request().Header.Get("X-User-ID")
			if userID != "" {
				c.Set("user_id", userID)
			}
			return next(c)
		}
	})

	return e
}

func makeRequest(
	t *testing.T,
	e *echo.Echo,
	method, path string,
	body interface{},
	userID string,
) *httptest.ResponseRecorder {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()

	// Set user_id before routing
	req.Header.Set("X-User-ID", userID) // Or use middleware on e

	e.ServeHTTP(rec, req)
	return rec
}
func seedUser(t *testing.T, username, email, password string) *ent.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("seedUser: %v", err)
	}
	u, err := db.User.Create().
		SetUsername(username).
		SetEmail(email).
		SetPasswordHash(string(hash)).
		Save(context.Background())
	if err != nil {
		t.Fatalf("seedUser: %v", err)
	}
	return u
}

func seedProject(t *testing.T, ownerID, name string) *ent.Project {
	t.Helper()
	oid, err := parseUUID(ownerID)
	if err != nil {
		t.Fatalf("seedProject: bad ownerID: %v", err)
	}
	p, err := db.Project.Create().
		SetName(name).
		SetOwnerID(oid).
		Save(context.Background())
	if err != nil {
		t.Fatalf("seedProject: %v", err)
	}
	return p
}

func seedEntry(t *testing.T, projectID, authorID, content string) *ent.Entry {
	t.Helper()
	pid, err := parseUUID(projectID)
	if err != nil {
		t.Fatalf("seedEntry: bad projectID: %v", err)
	}
	aid, err := parseUUID(authorID)
	if err != nil {
		t.Fatalf("seedEntry: bad authorID: %v", err)
	}
	e, err := db.Entry.Create().
		SetContent(content).
		SetProjectID(pid).
		SetAuthorID(aid).
		SetEntryDate(time.Now()).
		Save(context.Background())
	if err != nil {
		t.Fatalf("seedEntry: %v", err)
	}
	return e
}

// satisfy the compiler — net/http is used via http.MethodGet etc. in test files
var _ = http.MethodGet
