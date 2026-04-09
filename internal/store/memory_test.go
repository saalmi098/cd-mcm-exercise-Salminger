package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID returned unexpected error: %v", err)
	}
	if got.Name != "Widget" || got.Price != 9.99 {
		t.Errorf("got %+v, want Name=Widget Price=9.99", got)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestUpdateProduct(t *testing.T) {
	// table-driven tests for updating products
	cases := []struct {
		name      string
		initial   model.Product
		updated   model.Product
		wantName  string
		wantPrice float64
	}{
		{
			name:      "update name and price",
			initial:   model.Product{Name: "Old", Price: 1.00},
			updated:   model.Product{Name: "New", Price: 2.50},
			wantName:  "New",
			wantPrice: 2.50,
		},
		{
			name:      "update price only",
			initial:   model.Product{Name: "Gadget", Price: 5.00},
			updated:   model.Product{Name: "Gadget", Price: 19.99},
			wantName:  "Gadget",
			wantPrice: 19.99,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewMemoryStore()
			created := s.Create(tc.initial)

			got, err := s.Update(created.ID, tc.updated)
			if err != nil {
				t.Fatalf("Update returned unexpected error: %v", err)
			}
			if got.Name != tc.wantName || got.Price != tc.wantPrice {
				t.Errorf("got %+v, want Name=%s Price=%v", got, tc.wantName, tc.wantPrice)
			}
			if got.ID != created.ID {
				t.Errorf("ID changed after update: got %d, want %d", got.ID, created.ID)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Doomed", Price: 3.00})

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("Delete returned unexpected error: %v", err)
	}

	_, err := s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(42)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for non-existent ID, got %v", err)
	}
}
