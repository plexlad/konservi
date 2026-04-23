package main

import (
    "net/http"
    "github.com/labstack/echo/v5"
    //"github.com/plexlad/konservi/ent"
    //"github.com/plexlad/konservi/ent/projectmembership"
)

func CreateProjectMembershipEndpoint(c *echo.Context) error {
    // TODO: Implement based on project_membership schema [5]
    return c.JSON(http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func ListProjectMembershipsEndpoint(c *echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func UpdateProjectMembershipEndpoint(c *echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func DeleteProjectMembershipEndpoint(c *echo.Context) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}
