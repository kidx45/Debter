package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/kidx45/Debter/internal/adapter/outbound/postgres"
	"github.com/kidx45/Debter/internal/service"
	"github.com/lib/pq"
)

type UserAdapter struct {
	UserService *service.UserService
}

func NewUserAdapter(UserService service.UserService) *UserAdapter {
	return &UserAdapter{
		UserService: &UserService,
	}
}

type userResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (a *UserAdapter) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := db.User{
		Username:       req.Username,
		HashedPassword: req.Password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	res, err := a.UserService.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, newUserResponse(res))
}

func (a *UserAdapter) GetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	res, err := a.UserService.GetUserByUsername(ctx.Request.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(res))
}

type updateFullNameRequest struct {
	FullName string `json:"fullName" binding:"required"`
}

func (a *UserAdapter) UpdateFullNameByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	var req updateFullNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.UserService.UpdateFullNameByUsername(ctx.Request.Context(), username, req.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(res))
}

type updateUsernameRequest struct {
	NewUsername string `json:"newUsername" binding:"required"`
}

func (a *UserAdapter) UpdateUserNameByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	var req updateUsernameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.UserService.UpdateUserNameByUsername(ctx.Request.Context(), username, req.NewUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(res))
}

func (a *UserAdapter) DeleteUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	err := a.UserService.DeleteUserByUsername(ctx.Request.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
