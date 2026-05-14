package httphandler

import "github.com/gin-gonic/gin"

func (p *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/admin/products", p.CreateProduct)
	router.GET("/products/:id", p.GetProductByID)
}
