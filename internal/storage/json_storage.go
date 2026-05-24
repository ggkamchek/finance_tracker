package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"finance_tracker/internal/transactions"
)

type JSONStorage struct {
	path string
}

func NewJSONStorage(path string) *JSONStorage {
	return &JSONStorage{path: path}
}

func (s *JSONStorage) Load() ([]transactions.Transaction, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []transactions.Transaction{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []transactions.Transaction{}, nil
	}

	var items []transactions.Transaction
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *JSONStorage) Save(items []transactions.Transaction) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0o644)
}
