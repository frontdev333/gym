# API Трекера Тренировок

Backend-сервис для учёта выполненных упражнений и получения агрегированной статистики.

## Быстрый старт

```bash
# 1. Создайте .env файл
echo 'POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=gym' > .env

# 2. Запустите приложение
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

## API

### Health check

```bash
curl http://localhost:8080/health
```

### Пользователи

```bash
# Создать
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Иван"}'

# Список
curl http://localhost:8080/api/v1/users

# Получить по ID
curl http://localhost:8080/api/v1/users/{user_id}

# Обновить
curl -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Новое имя"}'

# Удалить (каскадно удаляет workouts)
curl -X DELETE http://localhost:8080/api/v1/users/{user_id}
```

### Упражнения

```bash
# Создать
curl -X POST http://localhost:8080/api/v1/exercises \
  -H "Content-Type: application/json" \
  -d '{"title": "Жим лёжа"}'

# Список / получить / обновить / удалить
curl http://localhost:8080/api/v1/exercises
curl http://localhost:8080/api/v1/exercises/{exercise_id}
curl -X PUT http://localhost:8080/api/v1/exercises/{exercise_id} \
  -H "Content-Type: application/json" \
  -d '{"title": "Приседания"}'
curl -X DELETE http://localhost:8080/api/v1/exercises/{exercise_id}
```

### Тренировки

```bash
# Зафиксировать выполнение
curl -X POST http://localhost:8080/api/v1/users/{user_id}/workouts \
  -H "Content-Type: application/json" \
  -d '{"exercise_id": "{exercise_id}", "amount": 10}'

# Список тренировок пользователя
curl http://localhost:8080/api/v1/users/{user_id}/workouts
```

Поля `performed_at` и `amount` опциональны. Если `performed_at` не указан, используется текущее время.

### Статистика

```bash
curl http://localhost:8080/api/v1/users/{user_id}/statistics
```

Пример ответа:

```json
{
  "total": 42,
  "today": 3,
  "last_7_days": [
    { "date": "2025-06-04", "count": 5 },
    { "date": "2025-06-05", "count": 0 },
    { "date": "2026-06-06", "count": 2 }
  ]
}
```

## Структура проекта

```
cmd/server/           — точка входа
internal/
  config/             — конфигурация из env
  domain/             — сущности и доменные ошибки
  handler/            — HTTP-слой
  repository/         — доступ к PostgreSQL
  service/            — бизнес-логика
migrations/           — SQL-миграции
```

## Стек

- **Go 1.26 + chi/v5** — HTTP-роутер
- **pgx/v5** — драйвер PostgreSQL
- **PostgreSQL 17** — база данных
- **Docker Compose** — запуск (postgres + migrate + backend)
# gym
