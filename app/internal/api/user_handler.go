package api

import (
	"app/internal/models"
	"app/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepository repository.UserRepository
}

func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	user := &models.User{
		ID:       uuid.New().String(),
		Username: ctx.PostForm("username"),
		Email:    ctx.PostForm("email"),
	}
	if err := h.userRepository.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.userRepository.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	users, err := h.userRepository.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user := &models.User{
		ID:       id,
		Username: ctx.PostForm("username"),
		Email:    ctx.PostForm("email"),
	}
	if err := h.userRepository.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.userRepository.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
