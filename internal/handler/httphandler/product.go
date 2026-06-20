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

type ProductService interface {
	CreateProduct(ctx context.Context, item domain.Product) (uuid.UUID, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (domain.Product, error)
	GetProducts(ctx context.Context, queryParams domain.ProductFilter) (domain.ProductList, error)
	UpdateProduct(ctx context.Context, item domain.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ChangeProductStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error
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
			"error": bindErrorMessage(err),
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
			log.Printf("create product %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}})
}

func (p *ProductHandler) GetProductByID(c *gin.Context) {

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := p.service.GetProductByID(c.Request.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found by ID"})
		default:
			log.Printf("get product by id %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.ProductToResponse(product)

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (p *ProductHandler) ListProducts(c *gin.Context) {

	var queryParams dto.ListProductsQuery

	if err := c.ShouldBindQuery(&queryParams); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": bindErrorMessage(err)})
		return

	}
	filter := dto.QueryToProductFilter(queryParams)

	productsStruct, err := p.service.GetProducts(c.Request.Context(), filter)

	if err != nil {
		log.Printf("list products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := []dto.ProductResponse{}
	for _, item := range productsStruct.Items {
		resp = append(resp, dto.ProductToResponse(item))
	}

	c.JSON(http.StatusOK, gin.H{"data": resp, "total": productsStruct.Total, "page": queryParams.Page, "limit": queryParams.Limit})
}

func (p *ProductHandler) UpdateProduct(c *gin.Context) {

	var req dto.UpdateProductRequest

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErrorMessage(err)})

		return
	}

	product := dto.UpdateRequestToProduct(id, req)

	err = p.service.UpdateProduct(c.Request.Context(), product)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			log.Printf("update product %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (p *ProductHandler) DeleteProduct(c *gin.Context) {

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = p.service.DeleteProduct(c.Request.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			log.Printf("delete product %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}

func (p *ProductHandler) ChangeProductStatus(c *gin.Context) {

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var req dto.ChangeStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": bindErrorMessage(err),
		})

		return
	}

	err = p.service.ChangeProductStatus(c.Request.Context(), id, req.Status)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			log.Printf("change product status %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)

}
