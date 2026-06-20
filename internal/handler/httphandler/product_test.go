package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JJIiSSH/jewelry-store/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockProductService struct {
	createProductID  uuid.UUID
	createProductErr error
	capturedProduct  domain.Product
}

func (m *mockProductService) CreateProduct(ctx context.Context, item domain.Product) (uuid.UUID, error) {
	m.capturedProduct = item
	return m.createProductID, m.createProductErr
}

func (m *mockProductService) GetProductByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	return domain.Product{}, nil
}

func (m *mockProductService) GetProducts(ctx context.Context, queryParams domain.ProductFilter) (domain.ProductList, error) {
	return domain.ProductList{}, nil
}

func (m *mockProductService) UpdateProduct(ctx context.Context, item domain.Product) error {
	return nil
}

func (m *mockProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockProductService) ChangeProductStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error {
	return nil
}

func setupProductTestRouter(service ProductService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewProductHandler(service)
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1)

	return router
}

func validCreateProductReqData(categoryID uuid.UUID, materials []string) map[string]any {

	return map[string]any{
		"title":       "Gold chain",
		"price":       3250,
		"stone":       "obsidian",
		"materials":   materials,
		"weight_g":    100,
		"category_id": categoryID.String(),
		"stock":       1,
		"status":      domain.ProductStatusDraft,
		"description": "test description about chain",
		"story":       "test story",
		"size":        "12",
	}

}

func TestProductHandler_CreateProduct_Success(t *testing.T) {

	expectedID := uuid.New()

	service := &mockProductService{createProductID: expectedID}

	router := setupProductTestRouter(service)

	categoryID := uuid.New()
	materials := []string{"Rose gold", "Silver"}

	reqData := validCreateProductReqData(categoryID, materials)

	jsonBytes, err := json.Marshal(reqData)

	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(jsonBytes)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Data.ID != expectedID.String() {
		t.Fatalf("expected id %q, got %q", expectedID, resp.Data.ID)
	}

	if service.capturedProduct.Title != "Gold chain" {
		t.Fatalf("expected title %q, got %q", "Gold chain", service.capturedProduct.Title)
	}
	if service.capturedProduct.Price != 3250 {
		t.Fatalf("expected price %v, got %v", 3250, service.capturedProduct.Price)
	}
	if service.capturedProduct.Status != domain.ProductStatusDraft {
		t.Fatalf("expected Status %q, got %q", domain.ProductStatusDraft, service.capturedProduct.Status)
	}

	assert.Equal(t, materials, service.capturedProduct.Materials, "unexpected materials")
	assert.Equal(t, categoryID, service.capturedProduct.CategoryID, "unexpected categoryID")
}

func TestProductHandler_CreateProduct_ConflictPath(t *testing.T) {

	service := &mockProductService{createProductErr: domain.ErrConflict}

	router := setupProductTestRouter(service)

	categoryID := uuid.New()
	materials := []string{"Rose gold", "Silver"}

	reqData := validCreateProductReqData(categoryID, materials)

	jsonBytes, err := json.Marshal(reqData)

	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(jsonBytes)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var resp map[string]string

	if err = json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected only error field, got %v", resp)
	}

	if resp["error"] != "product already exist" {
		t.Fatalf("unexpected public error: %q", resp["error"])
	}
}

// TestProductHandler_CreateProduct_BadRequest asserts that an invalid body
// returns 400 and never leaks the internal Go type name (CreateProductRequest)
// from gin's validator error. It pins the public error contract.
func TestProductHandler_CreateProduct_BadRequest(t *testing.T) {

	service := &mockProductService{}

	router := setupProductTestRouter(service)

	categoryID := uuid.New()
	materials := []string{"Rose gold", "Silver"}

	reqData := validCreateProductReqData(categoryID, materials)

	delete(reqData, "story")

	jsonBytes, err := json.Marshal(reqData)

	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(jsonBytes)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected error %v, got %v", http.StatusBadRequest, w.Code)
	}

	bodyText := w.Body.String()

	if strings.Contains(bodyText, "Story") {
		t.Errorf("response leaked internal error details: %v", bodyText)
	}

}
