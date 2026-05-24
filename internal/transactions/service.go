package transactions

import (
	"strings"
	"time"
)

type Service struct {
	items  []Transaction
	nextID int
}

func NewService(initial []Transaction) *Service {
	maxID := 0
	items := make([]Transaction, len(initial))
	copy(items, initial)
	for i, tx := range items {
		items[i].Type = normalizeType(tx.Type)
		if tx.ID > maxID {
			maxID = tx.ID
		}
	}

	return &Service{
		items:  items,
		nextID: maxID + 1,
	}
}

func (s *Service) Add(input CreateInput) (Transaction, error) {
	clean, err := validateInput(input)
	if err != nil {
		return Transaction{}, err
	}

	tx := Transaction{
		ID:          s.nextID,
		Type:        clean.Type,
		Amount:      clean.Amount,
		Category:    clean.Category,
		Date:        clean.Date,
		Description: clean.Description,
	}

	s.nextID++
	s.items = append(s.items, tx)

	return tx, nil
}

func (s *Service) Update(id int, input CreateInput) (Transaction, error) {
	if id <= 0 {
		return Transaction{}, ErrInvalidID
	}

	clean, err := validateInput(input)
	if err != nil {
		return Transaction{}, err
	}

	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Type = clean.Type
			s.items[i].Amount = clean.Amount
			s.items[i].Category = clean.Category
			s.items[i].Date = clean.Date
			s.items[i].Description = clean.Description
			return s.items[i], nil
		}
	}

	return Transaction{}, ErrTransactionAbsent
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	for i := range s.items {
		if s.items[i].ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}

	return ErrTransactionAbsent
}

func (s *Service) List() []Transaction {
	result := make([]Transaction, len(s.items))
	copy(result, s.items)
	return result
}

func (s *Service) FilterByCategory(category string) []Transaction {
	cleanCategory := strings.ToLower(strings.TrimSpace(category))
	if cleanCategory == "" {
		return []Transaction{}
	}

	filtered := make([]Transaction, 0)
	for _, tx := range s.items {
		if strings.ToLower(tx.Category) == cleanCategory {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

func (s *Service) Balance() float64 {
	var balance float64
	for _, tx := range s.items {
		if tx.Type == TypeIncome {
			balance += tx.Amount
		} else {
			balance -= tx.Amount
		}
	}

	return balance
}

func validateInput(input CreateInput) (CreateInput, error) {
	clean := CreateInput{
		Type:        normalizeType(input.Type),
		Amount:      input.Amount,
		Category:    strings.ToLower(strings.TrimSpace(input.Category)),
		Date:        strings.TrimSpace(input.Date),
		Description: strings.TrimSpace(input.Description),
	}

	if clean.Type != TypeIncome && clean.Type != TypeExpense {
		return CreateInput{}, ErrInvalidType
	}
	if clean.Amount <= 0 {
		return CreateInput{}, ErrInvalidAmount
	}
	if clean.Category == "" {
		return CreateInput{}, ErrEmptyCategory
	}
	if _, ok := allowedCategories[clean.Category]; !ok {
		return CreateInput{}, ErrUnknownCategory
	}
	if clean.Date == "" {
		return CreateInput{}, ErrEmptyDate
	}
	if _, err := time.Parse("02.01.2006", clean.Date); err != nil {
		return CreateInput{}, ErrInvalidDate
	}
	if clean.Description == "" {
		return CreateInput{}, ErrEmptyDescription
	}

	return clean, nil
}

func normalizeType(raw string) string {
	clean := strings.ToLower(strings.TrimSpace(raw))

	switch clean {
	case "income", "доход":
		return TypeIncome
	case "expense", "расход":
		return TypeExpense
	default:
		return clean
	}
}
