package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	return &PostgresStore{DB: db}, mock
}

func TestPostgresPing(t *testing.T) {
	s, _ := newMockStore(t)
	// sqlmock v1 accepts Ping without expectation; just verify no panic
	_ = s.Ping()
}

func TestPostgresEnsureTable(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))
	if err := s.EnsureTable(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPostgresGetAll(t *testing.T) {
	s, mock := newMockStore(t)
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99).
		AddRow(2, "Gadget", 4.99)
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestPostgresGetAllEmpty(t *testing.T) {
	s, mock := newMockStore(t)
	rows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestPostgresGetByIDFound(t *testing.T) {
	s, mock := newMockStore(t)
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(1).
		WillReturnRows(rows)

	p, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Widget" {
		t.Errorf("expected Widget, got %s", p.Name)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	rows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").
		WithArgs(999).
		WillReturnRows(rows)

	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Widget", 9.99).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	p, err := s.Create(newProduct("Widget", 9.99))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresUpdateFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("UPDATE products").
		WithArgs("Updated", 19.99, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	p, err := s.Update(1, newProduct("Updated", 19.99))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("UPDATE products").
		WithArgs("X", 1.0, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := s.Update(999, newProduct("X", 1.0))
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDeleteFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("DELETE FROM products WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := s.Delete(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	s, mock := newMockStore(t)
	mock.ExpectExec("DELETE FROM products WHERE id").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.Delete(999); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
