package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

type mockPostgresStore struct {
	pingErr    error
	products   []model.Product
	getAllErr  error
	getByIDErr error
	createErr  error
	updateErr  error
	deleteErr  error
}

func (m *mockPostgresStore) Ping() error { return m.pingErr }

func (m *mockPostgresStore) GetAll() ([]model.Product, error) {
	return m.products, m.getAllErr
}

func (m *mockPostgresStore) GetByID(id int) (model.Product, error) {
	if m.getByIDErr != nil {
		return model.Product{}, m.getByIDErr
	}
	for _, p := range m.products {
		if p.ID == id {
			return p, nil
		}
	}
	return model.Product{}, store.ErrNotFound
}

func (m *mockPostgresStore) Create(p model.Product) (model.Product, error) {
	if m.createErr != nil {
		return model.Product{}, m.createErr
	}
	p.ID = len(m.products) + 1
	m.products = append(m.products, p)
	return p, nil
}

func (m *mockPostgresStore) Update(id int, p model.Product) (model.Product, error) {
	if m.updateErr != nil {
		return model.Product{}, m.updateErr
	}
	p.ID = id
	return p, nil
}

func (m *mockPostgresStore) Delete(id int) error {
	return m.deleteErr
}

func setupPostgresRouter(mock *mockPostgresStore) *mux.Router {
	h := &PostgresHandler{Store: mock}
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func TestPostgresHealth(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHealthDBDown(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{pingErr: errors.New("down")})
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
}

func TestPostgresGetProducts(t *testing.T) {
	mock := &mockPostgresStore{products: []model.Product{{ID: 1, Name: "Widget", Price: 9.99}}}
	r := setupPostgresRouter(mock)
	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductsError(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{getAllErr: errors.New("db error")})
	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresGetProductFound(t *testing.T) {
	mock := &mockPostgresStore{products: []model.Product{{ID: 1, Name: "Widget", Price: 9.99}}}
	r := setupPostgresRouter(mock)
	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductNotFound(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresCreateProduct(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidJSON(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	req := httptest.NewRequest("POST", "/products", strings.NewReader("bad"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidValidation(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	body := `{"name":"","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductStoreError(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{createErr: errors.New("db error")})
	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresUpdateProduct(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	body := `{"name":"Updated","price":19.99}`
	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductInvalidPayload(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader("bad"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductNotFound(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{updateErr: store.ErrNotFound})
	body := `{"name":"X","price":1.0}`
	req := httptest.NewRequest("PUT", "/products/999", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresDeleteProduct(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{})
	req := httptest.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresDeleteProductNotFound(t *testing.T) {
	r := setupPostgresRouter(&mockPostgresStore{deleteErr: store.ErrNotFound})
	req := httptest.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
