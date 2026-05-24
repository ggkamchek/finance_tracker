package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"finance_tracker/internal/storage"
	"finance_tracker/internal/transactions"
)

const dataFilePath = "data/transactions.json"

func main() {
	store := storage.NewJSONStorage(dataFilePath)

	initial, err := store.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки данных: %v\n", err)
		return
	}

	service := transactions.NewService(initial)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Финансовый трекер (Go CLI)")
	fmt.Println("Введите `помощь`, чтобы увидеть команды.")

	for {
		fmt.Print("\n> ")
		command, err := readLine(reader)
		if err != nil {
			fmt.Printf("Ошибка чтения команды: %v\n", err)
			continue
		}

		switch strings.ToLower(command) {
		case "help", "помощь":
			printHelp()
		case "add", "добавить":
			handleAdd(service, reader)
		case "list", "список":
			printTransactions(service.List())
		case "edit", "изменить":
			handleEdit(service, reader)
		case "delete", "удалить":
			handleDelete(service, reader)
		case "filter", "фильтр":
			handleFilter(service, reader)
		case "balance", "баланс":
			fmt.Printf("Текущий баланс: %.2f\n", service.Balance())
		case "save", "сохранить":
			if err := store.Save(service.List()); err != nil {
				fmt.Printf("Ошибка сохранения: %v\n", err)
				continue
			}
			fmt.Println("Данные сохранены.")
		case "exit", "выход":
			if err := store.Save(service.List()); err != nil {
				fmt.Printf("Ошибка автосохранения: %v\n", err)
				return
			}
			fmt.Println("Данные сохранены. До свидания!")
			return
		case "":
			fmt.Println("Введите команду или `помощь`.")
		default:
			fmt.Println("Неизвестная команда. Используйте `помощь`.")
		}
	}
}

func printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  добавить  - добавить транзакцию")
	fmt.Println("  список    - показать все транзакции")
	fmt.Println("  изменить  - редактировать транзакцию по ID")
	fmt.Println("  удалить   - удалить транзакцию по ID")
	fmt.Println("  фильтр    - показать транзакции по категории")
	fmt.Println("  баланс    - показать текущий баланс")
	fmt.Println("  сохранить - сохранить данные в JSON")
	fmt.Println("  выход     - сохранить и выйти")
	fmt.Println("Также поддерживаются английские команды: help, add, list, edit, delete, filter, balance, save, exit")
	fmt.Println("Допустимые категории: еда, транспорт, зарплата, развлечения, здоровье, жилье, другое")
}

func handleAdd(service *transactions.Service, reader *bufio.Reader) {
	input, err := readTransactionInput(reader)
	if err != nil {
		fmt.Printf("Ошибка ввода: %v\n", err)
		return
	}

	tx, err := service.Add(input)
	if err != nil {
		fmt.Printf("Ошибка добавления: %v\n", err)
		return
	}

	fmt.Printf("Транзакция добавлена, ID: %d\n", tx.ID)
}

func handleEdit(service *transactions.Service, reader *bufio.Reader) {
	id, err := readID(reader)
	if err != nil {
		fmt.Printf("Ошибка ввода ID: %v\n", err)
		return
	}

	input, err := readTransactionInput(reader)
	if err != nil {
		fmt.Printf("Ошибка ввода: %v\n", err)
		return
	}

	tx, err := service.Update(id, input)
	if err != nil {
		fmt.Printf("Ошибка редактирования: %v\n", err)
		return
	}

	fmt.Printf("Транзакция %d обновлена.\n", tx.ID)
}

func handleDelete(service *transactions.Service, reader *bufio.Reader) {
	id, err := readID(reader)
	if err != nil {
		fmt.Printf("Ошибка ввода ID: %v\n", err)
		return
	}

	if err := service.Delete(id); err != nil {
		fmt.Printf("Ошибка удаления: %v\n", err)
		return
	}

	fmt.Printf("Транзакция %d удалена.\n", id)
}

func handleFilter(service *transactions.Service, reader *bufio.Reader) {
	fmt.Print("Категория: ")
	category, err := readLine(reader)
	if err != nil {
		fmt.Printf("Ошибка чтения категории: %v\n", err)
		return
	}

	filtered := service.FilterByCategory(category)
	if len(filtered) == 0 {
		fmt.Println("По этой категории транзакции не найдены.")
		return
	}

	printTransactions(filtered)
}

func readTransactionInput(reader *bufio.Reader) (transactions.CreateInput, error) {
	fmt.Print("Тип (доход/расход): ")
	t, err := readLine(reader)
	if err != nil {
		return transactions.CreateInput{}, err
	}

	fmt.Print("Сумма: ")
	amountRaw, err := readLine(reader)
	if err != nil {
		return transactions.CreateInput{}, err
	}
	amount, err := strconv.ParseFloat(amountRaw, 64)
	if err != nil {
		return transactions.CreateInput{}, fmt.Errorf("сумма должна быть числом")
	}

	fmt.Print("Категория: ")
	category, err := readLine(reader)
	if err != nil {
		return transactions.CreateInput{}, err
	}

	fmt.Print("Дата (ДД.ММ.ГГГГ): ")
	date, err := readLine(reader)
	if err != nil {
		return transactions.CreateInput{}, err
	}

	fmt.Print("Описание: ")
	description, err := readLine(reader)
	if err != nil {
		return transactions.CreateInput{}, err
	}

	return transactions.CreateInput{
		Type:        t,
		Amount:      amount,
		Category:    category,
		Date:        date,
		Description: description,
	}, nil
}

func readID(reader *bufio.Reader) (int, error) {
	fmt.Print("ID: ")
	raw, err := readLine(reader)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("ID должен быть целым числом")
	}

	return id, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func printTransactions(items []transactions.Transaction) {
	if len(items) == 0 {
		fmt.Println("Список транзакций пуст.")
		return
	}

	fmt.Println("Транзакции:")
	for _, tx := range items {
		fmt.Printf("ID=%d | Тип=%s | Сумма=%.2f | Категория=%s | Дата=%s | Описание=%s\n",
			tx.ID, tx.Type, tx.Amount, tx.Category, tx.Date, tx.Description)
	}
}
