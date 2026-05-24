package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"finance_tracker/internal/transactions"
)

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transactions.json")
	store := NewJSONStorage(path)

	items := []transactions.Transaction{
		{
			ID:          1,
			Type:        transactions.TypeIncome,
			Amount:      1000,
			Category:    "зарплата",
			Date:        "24.05.2026",
			Description: "зарплата",
		},
		{
			ID:          2,
			Type:        transactions.TypeExpense,
			Amount:      350,
			Category:    "еда",
			Date:        "25.05.2026",
			Description: "покупка продуктов",
		},
	}

	if err := store.Save(items); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !reflect.DeepEqual(items, loaded) {
		t.Fatalf("roundtrip mismatch\nexpected: %#v\ngot: %#v", items, loaded)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")
	store := NewJSONStorage(path)

	items, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %d items", len(items))
	}
}

func TestLoadInvalidJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	store := NewJSONStorage(path)

	if err := os.WriteFile(path, []byte("{invalid_json"), 0o644); err != nil {
		t.Fatalf("failed to write invalid json fixture: %v", err)
	}

	items, err := store.Load()
	if err == nil {
		t.Fatalf("expected json parse error, got nil")
	}
	if items != nil {
		t.Fatalf("expected nil items on parse error, got: %#v", items)
	}
}
