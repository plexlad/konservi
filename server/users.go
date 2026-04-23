package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/plexlad/konservi/ent"
	"github.com/plexlad/konservi/ent/user"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r CreateUserRequest) valid() bool {
	return r.Username != "" && r.Email != "" && r.Password != ""
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func toUserResponse(u *ent.User) UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func CreateUserEndpoint(c *echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil || !req.valid() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	u, err := db.User.Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPasswordHash(string(hash)).
		Save(c.Request().Context())
	if err != nil {
		if ent.IsConstraintError(err) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "username or email already taken"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	return c.JSON(http.StatusCreated, toUserResponse(u))
}

func GetUserEndpoint(c *echo.Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	u, err := db.User.Get(c.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
	}

	return c.JSON(http.StatusOK, toUserResponse(u))
}

func UpdateUserEndpoint(c *echo.Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	q := db.User.UpdateOneID(id)
	if req.Username != "" {
		q.SetUsername(req.Username)
	}
	if req.Email != "" {
		q.SetEmail(req.Email)
	}

	u, err := q.Save(c.Request().Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		if ent.IsConstraintError(err) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "username or email already taken"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update user"})
	}

	return c.JSON(http.StatusOK, toUserResponse(u))
}

func DeleteUserEndpoint(c *echo.Context) error {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	err = db.User.DeleteOneID(id).Exec(c.Request().Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}

func ListUsersEndpoint(c *echo.Context) error {
	users, err := db.User.Query().
		Order(ent.Asc(user.FieldCreatedAt)).
		All(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
	}

	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = toUserResponse(u)
	}
	return c.JSON(http.StatusOK, res)
}

