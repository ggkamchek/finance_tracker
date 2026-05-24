package transactions

import "testing"

func validInput() CreateInput {
	return CreateInput{
		Type:        TypeIncome,
		Amount:      1000,
		Category:    "зарплата",
		Date:        "24.05.2026",
		Description: "майская зарплата",
	}
}

func TestAddAndList(t *testing.T) {
	svc := NewService(nil)
	_, err := svc.Add(validInput())
	if err != nil {
		t.Fatalf("expected add success, got error: %v", err)
	}

	items := svc.List()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ID != 1 {
		t.Fatalf("expected id 1, got %d", items[0].ID)
	}
}

func TestBalance(t *testing.T) {
	svc := NewService(nil)
	_, _ = svc.Add(validInput())
	_, _ = svc.Add(CreateInput{
		Type:        TypeExpense,
		Amount:      250,
		Category:    "еда",
		Date:        "24.05.2026",
		Description: "покупки",
	})

	if got, want := svc.Balance(), 750.0; got != want {
		t.Fatalf("expected balance %.2f, got %.2f", want, got)
	}
}

func TestFilterByCategory(t *testing.T) {
	svc := NewService(nil)
	_, _ = svc.Add(CreateInput{
		Type:        TypeExpense,
		Amount:      100,
		Category:    "еда",
		Date:        "24.05.2026",
		Description: "обед",
	})
	_, _ = svc.Add(CreateInput{
		Type:        TypeExpense,
		Amount:      50,
		Category:    "транспорт",
		Date:        "24.05.2026",
		Description: "автобус",
	})

	filtered := svc.FilterByCategory("еда")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered item, got %d", len(filtered))
	}
	if filtered[0].Category != "еда" {
		t.Fatalf("expected category еда, got %s", filtered[0].Category)
	}
}

func TestUpdateTransaction(t *testing.T) {
	svc := NewService(nil)
	_, _ = svc.Add(validInput())

	updated, err := svc.Update(1, CreateInput{
		Type:        TypeExpense,
		Amount:      300,
		Category:    "развлечения",
		Date:        "25.05.2026",
		Description: "кино",
	})
	if err != nil {
		t.Fatalf("expected update success, got error: %v", err)
	}

	if updated.Type != TypeExpense || updated.Amount != 300 {
		t.Fatalf("unexpected updated values: %+v", updated)
	}
}

func TestDeleteTransaction(t *testing.T) {
	svc := NewService(nil)
	_, _ = svc.Add(validInput())

	if err := svc.Delete(1); err != nil {
		t.Fatalf("expected delete success, got error: %v", err)
	}

	if len(svc.List()) != 0 {
		t.Fatalf("expected empty list after delete")
	}
}

func TestRejectInvalidDate(t *testing.T) {
	svc := NewService(nil)
	_, err := svc.Add(CreateInput{
		Type:        TypeIncome,
		Amount:      10,
		Category:    "зарплата",
		Date:        "2026-05-24",
		Description: "ошибка даты",
	})

	if err != ErrInvalidDate {
		t.Fatalf("expected ErrInvalidDate, got %v", err)
	}
}

func TestRejectUnknownCategory(t *testing.T) {
	svc := NewService(nil)
	_, err := svc.Add(CreateInput{
		Type:        TypeIncome,
		Amount:      10,
		Category:    "инвестиции",
		Date:        "24.05.2026",
		Description: "дивиденды",
	})

	if err != ErrUnknownCategory {
		t.Fatalf("expected ErrUnknownCategory, got %v", err)
	}
}
