package dto

import "github.com/JJIiSSH/jewelry-store/internal/domain"

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
