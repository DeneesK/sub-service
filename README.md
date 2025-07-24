# Sub-Service
[![Tests](https://github.com/DeneesK/sub-service/actions/workflows/sub-service-tests.yml/badge.svg)](https://github.com/DeneesK/sub-service/actions/workflows/sub-service-tests.yml)

![Go version](https://img.shields.io/badge/go-1.23-blue)

REST-сервис для агрегации данных об онлайн-подписках пользователей.

---

## Описание

Сервис предоставляет CRUDL операции для подписок пользователей и возможность агрегировать суммарную стоимость подписок за указанный период с фильтрацией по пользователю и названию сервиса.

---

## Запуск

```bash
cp .env.example .env
docker-compose up
```

### Swagger документация доступна после заруска по адресу

```
http://localhost:8000/swagger/index.html
```
---

## Функциональность

- **CRUDL для подписок:**
  - Название сервиса (`service_name`)
  - Стоимость месячной подписки в рублях (`price`)
  - ID пользователя (`user_id`), UUID
  - Дата начала подписки (`start_date`, формат `MM-YYYY`)
  - Опционально дата окончания подписки (`end_date`), может быть `null`

- **Агрегация стоимости** подписок за период с фильтрами по `user_id` и `service_name`

- Используется PostgreSQL с миграциями для инициализации базы данных

- Логирование всех операций с уровнями логов

- Конфигурация через `.env`

- Документация API через Swagger

- Запуск через Docker Compose
- Автоматический запуск тестов через GitHub Actions

---

## API Endpoints

| Метод | Путь                                  | Описание                                   |
|-------|--------------------------------------|--------------------------------------------|
| POST  | `/api/v1/subs`               | Создать новую подписку                      |
| GET   | `/api/v1/subs/{id}`          | Получить подписку по ID                     |
| GET   | `/api/v1/subs?user_id=...`  | Список подписок, опциональный фильтр по пользователю |
| PATCH   | `/api/v1/subs/{id}`          | Обновить подписку                           |
| DELETE| `/api/v1/subs/{id}`          | Удалить подписку                            |
| GET   | `/api/v1/subs/aggregate`     | Получить сумму стоимости подписок за период с фильтрами |

---

## Пример запроса на создание подписки

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
