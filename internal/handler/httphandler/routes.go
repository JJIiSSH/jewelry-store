package httphandler

import "github.com/gin-gonic/gin"

func (p *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/admin/products", p.CreateProduct)
	router.GET("/products/:id", p.GetProductByID)
	// TODO: зарегистрировать 4 роута:
	//   GET    /products                          → ListProducts          (публичный)
	//   PUT    /admin/products/:id                → UpdateProduct         (админский)
	//   DELETE /admin/products/:id                → DeleteProduct         (админский)
	//   PATCH  /admin/products/:id/status         → ChangeProductStatus   (админский)
	//

	router.GET("/products", p.ListProducts)
	router.PUT("/admin/products/:id", p.UpdateProduct)
	router.DELETE("/admin/products/:id", p.DeleteProduct)
	router.PATCH("/admin/products/:id/status", p.ChangeProductStatus)
	// Заметь: ListProducts — публичный (без /admin), потому что каталог видят все.
	// Когда подключим auth middleware, /admin/* будет за JWT-проверкой роли admin —
	// тогда стоит пересмотреть структуру и сделать отдельный admin-RouterGroup.
}
