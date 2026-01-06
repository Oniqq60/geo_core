# Ядро системы геооповещений

Backend-сервис на Go для системы геооповещений. Сервис интегрируется с новостным порталом (Django) через вебхуки.

## Архитектура

Проект использует Clean Architecture:
- **Handler** - HTTP handlers (Gin)
- **Service** - бизнес-логика
- **Repository** - работа с БД (PostgreSQL) и кэшем (Redis)

## Требования

- Go 1.23+
- PostgreSQL 15+
- Redis 7+
- Docker и Docker Compose (для запуска через Docker)

## Быстрый старт

### 1. Запуск через Docker Compose

```bash
# Клонируем репозиторий
git clone <repository-url>
cd Geo_system_core

# Запускаем все сервисы
docker-compose up -d

# Проверяем статус
docker-compose ps
```

Сервис будет доступен на `http://localhost:8080`

### 2. Локальный запуск

#### Установка зависимостей

```bash
go mod download
```

#### Настройка базы данных

```bash
# Создайте базу данных
createdb geo_system

# Примените миграции
psql -U postgres -d geo_system -f migrations/001_init.sql
```

#### Настройка переменных окружения

Скопируйте `.env.example` в `.env` и настройте параметры:

```bash
cp .env.example .env
```

#### Запуск сервиса

```bash
go run cmd/server/main.go
```

## Настройка ngrok для тестирования вебхуков

1. Установите ngrok: https://ngrok.com/download

2. Запустите Go-заглушку вебхука на порту 9090:
```bash
# Вариант 1: через docker-compose (сервис webhook-stub)
docker-compose up webhook-stub

# Вариант 2: локально
go run webhook-stub/main.go
```

3. Запустите ngrok:
```bash
ngrok http 9090
```

4. Скопируйте полученный URL (например, `https://abc123.ngrok.io`) и установите в `.env`:
```env
WEBHOOK_URL=https://abc123.ngrok.io/webhook
```

## API Документация

### Health Check

```bash
GET /api/v1/system/health
```

**Ответ:**
```json
{
  "status": "ok",
  "service": "geo_system_core"
}
```

### Управление инцидентами (требует API-key)

Все эндпоинты требуют заголовок `X-API-Key` с валидным API ключом.

#### Создание инцидента

```bash
POST /api/v1/incidents
Content-Type: application/json
X-API-Key: your-api-key

{
  "title": "Пожар в лесу",
  "description": "Обнаружен пожар в лесном массиве",
  "latitude": 55.7558,
  "longitude": 37.6173,
  "radius": 500,
  "severity": "high",
  "status": "active"
}
```

**Ответ:**
```json
{
  "id": "uuid",
  "title": "Пожар в лесу",
  "description": "Обнаружен пожар в лесном массиве",
  "latitude": 55.7558,
  "longitude": 37.6173,
  "radius": 500,
  "severity": "high",
  "status": "active",
  "is_active": true,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### Получение списка инцидентов (с пагинацией)

```bash
GET /api/v1/incidents?page=1&limit=10
X-API-Key: your-api-key
```

**Ответ:**
```json
{
  "data": [...],
  "page": 1,
  "limit": 10,
  "total": 100,
  "total_pages": 10
}
```

#### Получение инцидента по ID

```bash
GET /api/v1/incidents/{id}
X-API-Key: your-api-key
```

#### Обновление инцидента

```bash
PUT /api/v1/incidents/{id}
Content-Type: application/json
X-API-Key: your-api-key

{
  "title": "Обновленное название",
  "severity": "critical"
}
```

#### Удаление (деактивация) инцидента

```bash
DELETE /api/v1/incidents/{id}
X-API-Key: your-api-key
```

### Проверка координат (публичный)

```bash
POST /api/v1/location/check
Content-Type: application/json

{
  "latitude": 55.7558,
  "longitude": 37.6173,
  "user_id": "user123"
}
```

**Ответ:**
```json
{
  "nearby_incidents": [
    {
      "id": "uuid",
      "title": "Пожар в лесу",
      "description": "...",
      "latitude": 55.7558,
      "longitude": 37.6173,
      "radius": 500,
      "severity": "high",
      "distance": 250.5
    }
  ],
  "has_danger": true
}
```

После ответа сервис:
1. Сохраняет факт проверки в БД
2. Если обнаружены опасности, ставит задачу на отправку вебхука

### Статистика по зонам

```bash
GET /api/v1/incidents/stats
```

**Ответ:**
```json
{
  "zones": [
    {
      "incident_id": "uuid",
      "title": "Пожар в лесу",
      "user_count": 15
    }
  ],
  "total": 15
}
```

Возвращает количество уникальных пользователей (`user_count`) для каждой зоны за последние N минут (настраивается через `STATS_TIME_WINDOW_MINUTES`).

## Примеры запросов (curl)

### Создание инцидента

```bash
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -H "X-API-Key: default-api-key-change-in-production" \
  -d '{
    "title": "Пожар",
    "description": "Пожар в лесу",
    "latitude": 55.7558,
    "longitude": 37.6173,
    "radius": 500,
    "severity": "high",
    "status": "active"
  }'
```

### Проверка координат

```bash
curl -X POST http://localhost:8080/api/v1/location/check \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 55.7558,
    "longitude": 37.6173,
    "user_id": "user123"
  }'
```

### Получение статистики

```bash
curl http://localhost:8080/api/v1/incidents/stats
```

### Получение списка инцидентов

```bash
curl -X GET "http://localhost:8080/api/v1/incidents?page=1&limit=10" \
  -H "X-API-Key: default-api-key-change-in-production"
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_HOST` | Хост сервера | `0.0.0.0` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres` |
| `DB_NAME` | Имя БД | `geo_system` |
| `DB_SSLMODE` | SSL режим | `disable` |
| `REDIS_HOST` | Хост Redis | `localhost` |
| `REDIS_PORT` | Порт Redis | `6379` |
| `REDIS_PASSWORD` | Пароль Redis | (пусто) |
| `REDIS_DB` | Номер БД Redis | `0` |
| `WEBHOOK_URL` | URL для отправки вебхуков | `http://localhost:9090/webhook` |
| `WEBHOOK_RETRY_ATTEMPTS` | Количество попыток retry | `3` |
| `WEBHOOK_RETRY_DELAY` | Задержка между попытками | `5s` |
| `WEBHOOK_TIMEOUT` | Таймаут запроса | `10s` |
| `STATS_TIME_WINDOW_MINUTES` | Окно времени для статистики | `60` |
| `API_KEY` | API ключ для операторов | `default-api-key-change-in-production` |

## Особенности реализации

### Асинхронная отправка вебхуков

- Вебхуки отправляются асинхронно через Redis очередь
- Worker обрабатывает очередь в фоновом режиме
- При ошибках доставки выполняется retry с экспоненциальной задержкой

### Кэширование

- Активные инциденты кэшируются в Redis для быстрого доступа
- TTL кэша настраивается

### Валидация

- Валидация координат (широта: -90 до 90, долгота: -180 до 180)
- Защита от SQL-инъекций через параметризованные запросы
- Валидация входных данных через Gin binding

### Безопасность

- API-key аутентификация для операторов
- Параметризованные SQL-запросы
- Валидация всех входных данных

## Структура проекта

```
Geo_system_core/
├── cmd/
│   └── server/
│       └── main.go          # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # Middleware (auth)
│   ├── models/              # Модели данных
│   ├── repository/          # Репозитории
│   │   ├── postgres/        # PostgreSQL репозитории
│   │   └── redis/           # Redis репозитории
│   ├── router/              # Настройка роутера
│   └── service/             # Бизнес-логика
├── migrations/              # SQL миграции
│   └── 001_init.sql
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Тестирование

### Запуск тестов

```bash
go test ./...
```

### Тестирование API

Используйте готовые скрипты для тестирования API:

**Linux/Mac:**
```bash
chmod +x scripts/test_api.sh
./scripts/test_api.sh
```

**Windows (PowerShell):**
```powershell
.\scripts\test_api.ps1
```

### Тестирование вебхуков

1. Запустите HTTP-сервер-заглушку на порту 9090
2. Настройте ngrok (см. раздел выше)
3. Выполните проверку координат в опасной зоне
4. Проверьте логи вебхука

## Разработка

### Добавление новых миграций

Создайте новый файл в `migrations/` с номером версии:
```
migrations/002_add_new_field.sql
```

### Запуск в режиме разработки

```bash
# Установите air для hot reload
go install github.com/cosmtrek/air@latest

# Запустите с hot reload
air
```

## Лицензия

MIT

