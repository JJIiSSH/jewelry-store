package httphandler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/JJIiSSH/jewelry-store/internal/handler/httphandler/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, user domain.User, password string) (uuid.UUID, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/auth/register", h.Register)
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErrorMessage(err),
		})
		return
	}
	user := dto.RegisterRequestToUser(req)

	id, err := h.service.Register(c.Request.Context(), user, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		case errors.Is(err, domain.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password too long"})
		default:
			log.Printf("register user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}})

}
