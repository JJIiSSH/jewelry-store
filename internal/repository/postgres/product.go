package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductByID(ctx context.Context, ID uuid.UUID) (domain.Product, error) {
	var product domain.Product
	var materials pq.StringArray

	query := `
	SELECT id, category_id, title, slug, description, story, 
       stone, materials, price, weight_g, size, 
       is_unique, stock, status, created_at, updated_at
	FROM products
	WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, ID)

	err := row.Scan(
		&product.ID,
		&product.CategoryID,
		&product.Title,
		&product.Slug,
		&product.Description,
		&product.Story,
		&product.Stone,
		&materials,
		&product.Price,
		&product.WeightG,
		&product.Size,
		&product.IsUnique,
		&product.Stock,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, fmt.Errorf("%w", domain.ErrNotFound)
		}

		return domain.Product{}, fmt.Errorf("error get product by id: %w", err)
	}

	product.Materials = []string(materials)

	return product, nil

}

func (r *ProductRepository) GetProducts(ctx context.Context, queryParams domain.ProductFilter) ([]domain.Product, error) {

	var products []domain.Product
	var args []any

	query := ` SELECT products.id, 
					  products.category_id, 
					  products.title,
					  products.slug,
					  products.description,
					  products.story,
					  products.stone,
					  products.materials,
					  products.price,
					  products.weight_g,
					  products.size, 
					  products.is_unique, 
					  products.stock, 
					  products.status, 
					  products.created_at,  
					  products.updated_at 
			   FROM products 
			   INNER JOIN categories ON products.category_id = categories.id
			   WHERE 1=1
			`
	argIdx := 1

	if queryParams.Stone != "" {
		query += fmt.Sprintf(" AND stone = $%d", argIdx)
		args = append(args, queryParams.Stone)
		argIdx++
	}

	if queryParams.Category != "" {

		query += fmt.Sprintf(" AND categories.slug = $%d", argIdx)

		args = append(args, queryParams.Category)
		argIdx++
	}

	if queryParams.MaxPrice != 0 {
		query += fmt.Sprintf(" AND price <= $%d", argIdx)
		args = append(args, queryParams.MaxPrice)
		argIdx++
	}

	if queryParams.MinPrice != 0 {
		query += fmt.Sprintf(" AND price >= $%d", argIdx)
		args = append(args, queryParams.MinPrice)
		argIdx++
	}

	if queryParams.Sort != "" {

		if queryParams.Sort == "price_asc" {
			query += " ORDER BY price ASC"
		}

		if queryParams.Sort == "price_desc" {
			query += " ORDER BY price DESC"
		}

		if queryParams.Sort == "newest" {
			query += " ORDER BY created_at DESC"
		}

	}

	if queryParams.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", queryParams.Limit)
		offset := (queryParams.Page - 1) * queryParams.Limit
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return []domain.Product{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var product domain.Product
		var materials pq.StringArray

		if err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Title,
			&product.Slug,
			&product.Description,
			&product.Story,
			&product.Stone,
			&materials,
			&product.Price,
			&product.WeightG,
			&product.Size,
			&product.IsUnique,
			&product.Stock,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt); err != nil {
			return nil, err
		}

		product.Materials = []string(materials)
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return []domain.Product{}, err
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, item domain.Product) (uuid.UUID, error) {

	query := `
	INSERT INTO products (id, category_id, title, slug, description, story, 
       stone, materials, price, weight_g, size, 
       is_unique, stock, status, created_at, updated_at)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.CategoryID,
		item.Title,
		item.Slug,
		item.Description,
		item.Story,
		item.Stone,
		pq.StringArray(item.Materials),
		item.Price,
		item.WeightG,
		item.Size,
		item.IsUnique,
		item.Stock,
		item.Status,
		item.CreatedAt,
		item.UpdatedAt)

	if err != nil {
		return uuid.Nil, err
	}

	return item.ID, nil

}

func (r *ProductRepository) UpdateProduct(ctx context.Context, item domain.Product) error {

	query := `
		UPDATE products
		SET 
			category_id = $1, 
			title = $2, 
			slug = $3, 
			description = $4, 
			story = $5,
			stone = $6, 
			materials = $7, 
			price = $8, 			
			weight_g = $9, 
			size = $10,
			is_unique = $11,
			stock = $12, 
			status = $13,
			updated_at = $14

		WHERE id = $15
	`
	res, err := r.db.ExecContext(ctx, query,
		item.CategoryID,
		item.Title,
		item.Slug,
		item.Description,
		item.Story,
		item.Stone,
		pq.StringArray(item.Materials),
		item.Price,
		item.WeightG,
		item.Size,
		item.IsUnique,
		item.Stock,
		item.Status,
		item.UpdatedAt,
		item.ID)

	if err != nil {
		return err
	}
	countRows, _ := res.RowsAffected()

	if countRows == 0 {
		return fmt.Errorf("%w", domain.ErrNotFound)
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {

	query := "DELETE FROM products WHERE id = $1"

	res, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	countRows, _ := res.RowsAffected()

	if countRows == 0 {
		return fmt.Errorf("%w", domain.ErrNotFound)
	}

	return nil
}

func (r *ProductRepository) ChangeProductStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error {

	query := "UPDATE products SET status = $1 WHERE id = $2"

	res, err := r.db.ExecContext(ctx, query,
		status,
		id)

	if err != nil {
		return err
	}

	countRows, _ := res.RowsAffected()

	if countRows == 0 {
		return fmt.Errorf("%w", domain.ErrNotFound)
	}

	return nil
}
