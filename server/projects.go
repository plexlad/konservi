package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/plexlad/konservi/ent"
	"github.com/plexlad/konservi/ent/project"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r CreateProjectRequest) valid() bool {
	return r.Name != ""
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	CreatedAt   string `json:"created_at"`
}

func toProjectResponse(p *ent.Project) ProjectResponse {
	return ProjectResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID.String(),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	}
}

func callerUserID(c *echo.Context) (uuid.UUID, error) {
	return parseUUID(c.Get("user_id").(string))
}

func CreateProjectEndpoint(c *echo.Context) error {
	ownerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	var req CreateProjectRequest
	if err := c.Bind(&req); err != nil || !req.valid() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	p, err := db.Project.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetOwnerID(ownerID).
		Save(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create project"})
	}

	return c.JSON(http.StatusCreated, toProjectResponse(p))
}

func GetProjectEndpoint(c *echo.Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	p, err := db.Project.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get project"})
	}

	return c.JSON(http.StatusOK, toProjectResponse(p))
}

func UpdateProjectEndpoint(c *echo.Context) error {
	callerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	p, err := db.Project.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get project"})
	}
	if p.OwnerID != callerID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	var req UpdateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	q := db.Project.UpdateOneID(id)
	if req.Name != "" {
		q.SetName(req.Name)
	}
	if req.Description != "" {
		q.SetDescription(req.Description)
	}

	p, err = q.Save(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update project"})
	}

	return c.JSON(http.StatusOK, toProjectResponse(p))
}

func DeleteProjectEndpoint(c *echo.Context) error {
	callerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	p, err := db.Project.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get project"})
	}
	if p.OwnerID != callerID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	if err := db.Project.DeleteOneID(id).Exec(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete project"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "project deleted"})
}

func ListProjectsEndpoint(c *echo.Context) error {
	callerID, err := callerUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
	}

	projects, err := db.Project.Query().
		Where(project.OwnerID(callerID)).
		Order(ent.Asc(project.FieldCreatedAt)).
		All(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list projects"})
	}

	res := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		res[i] = toProjectResponse(p)
	}
	return c.JSON(http.StatusOK, res)
}

