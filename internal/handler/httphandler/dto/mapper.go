package dto

import (
	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/google/uuid"
)

func RequestToProduct(r CreateProductRequest) domain.Product {

	return domain.Product{
		Title:       r.Title,
		IsUnique:    r.IsUnique,
		Price:       r.Price,
		Materials:   r.Materials,
		Description: r.Description,
		Story:       r.Story,
		Stone:       r.Stone,
		Size:        r.Size,
		Status:      r.Status,
		CategoryID:  r.CategoryID,
		WeightG:     r.WeightG,
		Stock:       r.Stock,
	}
}

func UpdateRequestToProduct(id uuid.UUID, r UpdateProductRequest) domain.Product {
	return domain.Product{
		ID:          id,
		Title:       r.Title,
		IsUnique:    r.IsUnique,
		Price:       r.Price,
		Materials:   r.Materials,
		Description: r.Description,
		Story:       r.Story,
		Stone:       r.Stone,
		Size:        r.Size,
		CategoryID:  r.CategoryID,
		WeightG:     r.WeightG,
		Stock:       r.Stock,
	}
}

func QueryToProductFilter(q ListProductsQuery) domain.ProductFilter {

	if q.Limit > 100 {
		q.Limit = 100
	}
	sort := domain.Sort(q.Sort)

	return domain.ProductFilter{
		Stone:    q.Stone,
		MinPrice: q.MinPrice,
		MaxPrice: q.MaxPrice,
		Category: q.Category,
		Page:     q.Page,
		Limit:    q.Limit,
		Sort:     sort,
	}
}

func ProductToResponse(p domain.Product) ProductResponse {

	return ProductResponse{
		ID:          p.ID,
		Slug:        p.Slug,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Title:       p.Title,
		IsUnique:    p.IsUnique,
		Price:       p.Price,
		Materials:   p.Materials,
		Description: p.Description,
		Story:       p.Story,
		Stone:       p.Stone,
		Size:        p.Size,
		Status:      p.Status,
		CategoryID:  p.CategoryID,
		WeightG:     p.WeightG,
		Stock:       p.Stock,
	}

}
