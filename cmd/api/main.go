package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/JJIiSSH/jewelry-store/internal/repository/postgres"
	"github.com/JJIiSSH/jewelry-store/internal/service"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var testItem domain.Product
	var testMaterials []string
	var newCategoryID uuid.UUID

	queryCategory := "INSERT INTO categories (name, slug) VALUES ('Pendants', 'pendants') RETURNING id;"

	testMaterials = append(testMaterials, "gold")

	testItem.Title = "Noble Stone Pendant"
	testItem.Description = "Premium piece"
	testItem.Stone = "diamond"
	testItem.Status = domain.ProductStatusDraft
	testItem.Materials = testMaterials

	dsn := "postgres://IvanDev:1111@localhost:5432/mydb?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = db.QueryRow(queryCategory).Scan(&newCategoryID)

	if err != nil {
		log.Fatal(err)
	}

	testItem.CategoryID = newCategoryID

	productRepo := postgres.NewProductRepository(db)
	productService := service.NewProductService(productRepo)

	id, err := productService.CreateProduct(context.Background(), testItem)

	if err != nil {
		log.Fatal(err)
	}

	productByID, err := productService.GetProductByID(context.Background(), id)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(productByID.Title)
}
