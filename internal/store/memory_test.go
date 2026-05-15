package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newProduct(name string, price float64) model.Product {
	return model.Product{Name: name, Price: price}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetAllWithProducts(t *testing.T) {
	s := NewMemoryStore()
	s.Create(newProduct("A", 1.0))
	s.Create(newProduct("B", 2.0))
	products := s.GetAll()
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestCreateAssignsID(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(newProduct("Widget", 9.99))
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
	p2 := s.Create(newProduct("Gadget", 4.99))
	if p2.ID != 2 {
		t.Errorf("expected ID 2, got %d", p2.ID)
	}
}

func TestGetByIDFound(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(newProduct("Widget", 9.99))
	p, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Widget" {
		t.Errorf("expected name Widget, got %s", p.Name)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateExisting(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(newProduct("Widget", 9.99))
	updated, err := s.Update(created.ID, newProduct("Updated", 19.99))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected name Updated, got %s", updated.Name)
	}
	if updated.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, updated.ID)
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Update(999, newProduct("X", 1.0))
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExisting(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(newProduct("Widget", 9.99))
	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Error("expected product to be deleted")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}
