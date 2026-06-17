package httphandler

import "github.com/gin-gonic/gin"

func (p *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/admin/products", p.CreateProduct)
	router.GET("/products/:id", p.GetProductByID)

	router.GET("/products", p.ListProducts)
	router.PUT("/admin/products/:id", p.UpdateProduct)
	router.DELETE("/admin/products/:id", p.DeleteProduct)
	router.PATCH("/admin/products/:id/status", p.ChangeProductStatus)
	// Note: ListProducts is public (no /admin) — the catalog is visible to everyone.
	// Once auth middleware lands, /admin/* will sit behind a JWT admin-role check;
	// at that point consider splitting admin routes into a dedicated RouterGroup.
}
