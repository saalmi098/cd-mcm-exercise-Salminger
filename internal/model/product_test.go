package model

import (
	"testing"
)

func TestDiscountNormal(t *testing.T) {
	p := Product{Name: "Widget", Price: 100.0}
	got := p.Discount(20)
	if got != 80.0 {
		t.Errorf("Discount(20) = %v, want 80.0", got)
	}
}

func TestDiscountInvalidPct(t *testing.T) {
	p := Product{Name: "Widget", Price: 100.0}
	if p.Discount(-1) != 100.0 {
		t.Error("negative pct should return original price")
	}
	if p.Discount(101) != 100.0 {
		t.Error("pct > 100 should return original price")
	}
}

func TestIsExpensive(t *testing.T) {
	p := Product{Name: "Widget", Price: 50.0}
	if !p.IsExpensive(49.0) {
		t.Error("expected IsExpensive true when price > threshold")
	}
	if p.IsExpensive(50.0) {
		t.Error("expected IsExpensive false when price == threshold")
	}
	if p.IsExpensive(51.0) {
		t.Error("expected IsExpensive false when price < threshold")
	}
}

func TestString(t *testing.T) {
	p := Product{ID: 1, Name: "Widget", Price: 9.99}
	got := p.String()
	want := `Product{ID: 1, Name: "Widget", Price: 9.99}`
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestValidateEmptyName(t *testing.T) {
	p := Product{Name: "", Price: 10.0}
	if p.Validate() {
		t.Error("expected validation to fail for empty name")
	}
}

func TestValidateNegativePrice(t *testing.T) {
	p := Product{Name: "Widget", Price: -5.0}
	if p.Validate() {
		t.Error("expected validation to fail for negative price")
	}
}

func TestValidateValidProduct(t *testing.T) {
	p := Product{Name: "Widget", Price: 9.99}
	if !p.Validate() {
		t.Error("expected validation to pass for valid product")
	}
}
