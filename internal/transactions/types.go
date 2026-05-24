package transactions

import "errors"

const (
	TypeIncome  = "доход"
	TypeExpense = "расход"
)

var (
	ErrInvalidID         = errors.New("некорректный идентификатор транзакции")
	ErrInvalidType       = errors.New("тип операции должен быть: доход или расход")
	ErrInvalidAmount     = errors.New("сумма должна быть больше нуля")
	ErrEmptyCategory     = errors.New("категория обязательна")
	ErrUnknownCategory   = errors.New("неизвестная категория")
	ErrEmptyDate         = errors.New("дата обязательна")
	ErrInvalidDate       = errors.New("дата должна быть в формате DD.MM.YYYY")
	ErrEmptyDescription  = errors.New("описание обязательно")
	ErrTransactionAbsent = errors.New("транзакция не найдена")
)

var allowedCategories = map[string]struct{}{
	"еда":         {},
	"транспорт":   {},
	"зарплата":    {},
	"развлечения": {},
	"здоровье":    {},
	"жилье":       {},
	"другое":      {},
}

type Transaction struct {
	ID          int     `json:"id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

type CreateInput struct {
	Type        string
	Amount      float64
	Category    string
	Date        string
	Description string
}
