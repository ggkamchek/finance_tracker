# 03. UML/BPMN-диаграммы проекта

В документе собраны диаграммы в формате `mermaid`, описывающие поведение и структуру `finance_tracker`.

## 1) Use Case (варианты использования)

```mermaid
flowchart LR
    U[Пользователь] --> UC1((Добавить транзакцию))
    U --> UC2((Просмотреть список))
    U --> UC3((Изменить транзакцию))
    U --> UC4((Удалить транзакцию))
    U --> UC5((Фильтрация по категории))
    U --> UC6((Расчет баланса))
    U --> UC7((Сохранение в JSON))
    U --> UC8((Загрузка при старте))
```

## 2) Activity/BPMN (процесс выполнения команды)

```mermaid
flowchart TD
    A([Старт приложения]) --> B[Загрузить transactions.json]
    B --> C{Файл существует?}
    C -- Нет --> D[Создать пустой список]
    C -- Да --> E[Прочитать JSON]
    E --> F{JSON корректный?}
    F -- Нет --> G[Показать ошибку загрузки]
    F -- Да --> H[Показать меню команд]
    D --> H
    G --> H

    H --> I{Выбор команды}
    I -->|добавить| J[Ввести поля и выполнить валидацию]
    J --> K[Добавить запись]
    I -->|изменить| L[Изменить запись по ID]
    I -->|удалить| M[Удалить запись по ID]
    I -->|список| N[Вывести все записи]
    I -->|фильтр| O[Вывести выбранную категорию]
    I -->|баланс| P[Рассчитать итог]
    I -->|сохранить| Q[Сохранить список в JSON]
    I -->|выход| R([Завершение])

    K --> H
    L --> H
    M --> H
    N --> H
    O --> H
    P --> H
    Q --> H
```

## 3) Component (компонентная диаграмма)

```mermaid
flowchart LR
    CLI[CLI слой: main.go]
    SVC[Сервис: internal/transactions]
    STG[Хранилище: internal/storage]
    JSON[(data/transactions.json)]
    DOM[Модели: Transaction/Type]

    CLI --> SVC
    CLI --> STG
    SVC --> DOM
    STG --> DOM
    STG --> JSON
```

## 4) Sequence (последовательность: add + save)

```mermaid
sequenceDiagram
    actor U as Пользователь
    participant C as CLI
    participant S as Service
    participant J as JSONStorage
    participant F as transactions.json

    U->>C: add
    C->>U: Запрос type/amount/category/date/description
    U->>C: Ввод данных
    C->>S: Add(transactionInput)
    S->>S: validate + nextID
    S-->>C: OK
    C-->>U: Транзакция добавлена

    U->>C: save
    C->>S: List()
    S-->>C: []Transaction
    C->>J: Save(items)
    J->>F: write(json)
    F-->>J: OK
    J-->>C: OK
    C-->>U: Данные сохранены
```
