# efmob — сервис подписок

REST API подписок (Chi, PostgreSQL, goose). Точка входа: `cmd/subscriptionservice`.

## Запуск

### Переменные окружения

| Переменная     | Описание |
|----------------|----------|
| `DATABASE_URI` | DSN PostgreSQL для pgx, например `postgres://user:pass@localhost:5432/efmob?sslmode=disable` |
| `PORT`         | Порт HTTP-сервера |

### Docker

Для поднятия сервиса можно воспользоваться
```
make up
```

### Полезные команды

```bash
make test          # go test ./...
make mock          # go generate ./... (моки репозитория)
make swagger       # перегенерация Swagger в docs/
```

Swagger UI после запуска: `http://localhost:<PORT>/swagger/index.html`. Базовый путь API в спецификации: `/api/v1`.

---

## Схема PostgreSQL

Миграции лежат в `migrations/`. Итоговая модель после всех апгрейдов:

### `services`

| Колонка      | Тип              | Описание |
|-------------|------------------|----------|
| `service_id`| `SERIAL` PK      | Идентификатор сервиса |
| `name`      | `VARCHAR(255)` NOT NULL, **UNIQUE** | Название сервиса |

### `subscriptions`

| Колонка       | Тип           | Описание |
|--------------|---------------|----------|
| `user_id`    | `UUID`        | Часть составного первичного ключа |
| `service_id` | `INT`         | FK → `services(service_id)`, часть PK |
| `start_date` | `DATE` NOT NULL | Дата начала |
| `end_date`   | `DATE`        | Дата окончания (может быть NULL) |
| `price`      | `INT` NOT NULL | Цена; ограничение `CHECK (price > 0)` |

**Первичный ключ:** `(user_id, service_id)`.

Миграции подтягиваются в коде, при запуске сервиса. они расположены в директории ```migrations```