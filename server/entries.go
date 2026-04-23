package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/plexlad/konservi/ent"
	"github.com/plexlad/konservi/ent/entry"
)

type CreateEntryRequest struct {
	Content   string `json:"content"`
	EntryDate string `json:"entry_date"` // RFC3339
}

func (r CreateEntryRequest) valid() bool {
	return r.Content != ""
}

type UpdateEntryRequest struct {
	Content string `json:"content"`
}

type EntryResponse struct {
	ID           string `json:"id"`
	ProjectID    string `json:"project_id"`
	AuthorID     string `json:"author_id"`
	Content      string `json:"content"`
	EntryDate    string `json:"entry_date"`
	CreatedAt    string `json:"created_at"`
	EditedLastAt string `json:"edited_last_at"`
}

func toEntryResponse(e *ent.Entry) EntryResponse {
	return EntryResponse{
		ID:           e.ID.String(),
		ProjectID:    e.ProjectID.String(),
		AuthorID:     e.AuthorID.String(),
		Content:      e.Content,
		EntryDate:    e.EntryDate.Format(time.RFC3339),
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
		EditedLastAt: e.EditedLastAt.Format(time.RFC3339),
	}
}

func CreateEntryEndpoint(c *echo.Context) error {
	callerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	projectID, err := parseUUID(c.Param("project_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project_id"})
	}

	var req CreateEntryRequest
	if err := c.Bind(&req); err != nil || !req.valid() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	entryDate := time.Now()
	if req.EntryDate != "" {
		entryDate, err = time.Parse(time.RFC3339, req.EntryDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid entry_date format, use RFC3339"})
		}
	}

	_, err = db.Project.Get(c.Request().Context(), projectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to verify project"})
	}

	e, err := db.Entry.Create().
		SetContent(req.Content).
		SetProjectID(projectID).
		SetAuthorID(callerID).
		SetEntryDate(entryDate).
		Save(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create entry"})
	}

	return c.JSON(http.StatusCreated, toEntryResponse(e))
}

func GetEntryEndpoint(c *echo.Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	e, err := db.Entry.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "entry not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get entry"})
	}

	return c.JSON(http.StatusOK, toEntryResponse(e))
}

func UpdateEntryEndpoint(c *echo.Context) error {
	callerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	e, err := db.Entry.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "entry not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get entry"})
	}
	if e.AuthorID != callerID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	var req UpdateEntryRequest
	if err := c.Bind(&req); err != nil || req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	e, err = db.Entry.UpdateOneID(id).
		SetContent(req.Content).
		SetEditedLastAt(time.Now()).
		Save(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update entry"})
	}

	return c.JSON(http.StatusOK, toEntryResponse(e))
}

func DeleteEntryEndpoint(c *echo.Context) error {
	callerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	e, err := db.Entry.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "entry not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get entry"})
	}
	if e.AuthorID != callerID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	if err := db.Entry.DeleteOneID(id).Exec(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete entry"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "entry deleted"})
}

func ListEntriesEndpoint(c *echo.Context) error {
	projectID, err := parseUUID(c.Param("project_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project_id"})
	}

	entries, err := db.Entry.Query().
		Where(entry.ProjectID(projectID)).
		Order(ent.Desc(entry.FieldEntryDate)).
		All(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list entries"})
	}

	res := make([]EntryResponse, len(entries))
	for i, e := range entries {
		res[i] = toEntryResponse(e)
	}
	return c.JSON(http.StatusOK, res)
}

