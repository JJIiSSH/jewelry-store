package httphandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/JJIiSSH/jewelry-store/internal/handler/httphandler/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, item domain.Product) (uuid.UUID, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (p *ProductHandler) CreateProduct(c *gin.Context) {

	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	product := dto.RequestToProduct(req)

	id, err := p.service.CreateProduct(c.Request.Context(), product)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": "product already exist"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

func (p *ProductHandler) GetProductByID(c *gin.Context) {

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	product, err := p.service.GetProductByID(c.Request.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found by ID"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.ProductToResponse(product)

	c.JSON(http.StatusOK, resp)
}
