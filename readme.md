# Avito Merch

**Avito Merch** — сервис для управления покупкой мерча, тестовое задание для стажёра Backend-направления (зимняя волна 2025)

## Реализованные возможности

- Авторизация
- Покупка мерча
- Получение информации о балансе, купленном мерче и истории транзакций
- Отправка монет другому пользователю

## Стек технологий

- **Go** — язык разработки
- **PostgreSQL** — база данных
- **pgx** — драйвер для PostgreSQL
- **chi** — роутер
- **Docker Compose** — контейнеризация и оркестрация
- **Testconteiners** - для проведения тестирования

## Запуск проекта с помощью Docker Compose

### Предварительные требования

- Docker и Docker Compose установлены в системе

### Шаги для запуска

1. **Клонируйте репозиторий**:

   ```bash
   git clone https://github.com/plasmatrip/avito_merch.git
   cd avito_merch
   ```

2. **Запустите проект** с помощью Docker Compose:

   ```bash
   docker-compose up --build
   ```

3. **API будет доступно по адресу**: `http://localhost:8080`

## Структура проекта

- `build/Dockerfile` - Docker-файл для сборки приложения
- `docs/diagram/er_diagram.png` - ER-диаграмма базы данных
- `docs/swagger/openapi: 3.yml` - Swagger документ с описание REST API сервиса
- `cmd` — точка входа
- `internal` — основная логика и хендлеры
- `internal/storage/db/init` — SQL-скрипт создания базы данных
- `internal/storage/db/init_test` — bash-скрипт создания тестовой базы данных с помощью библтотеки Testcontainers
- `internal/storage/db/migrations` — миграции базы данных
- `docker-compose.yml` — конфигурация Docker Compose

## API Эндпоинты

- `POST /buy/{item}` — покупка мерча
- `POST /send-coin` — перевод монет между пользователями
- `GET /info` — получение информации о пользователе

## Тестирование

```bash
go test -v ./... -coverprofile=cover.out -covermode=atomic
go tool cover -html cover.out -o cover.html
```

## Контакты

- **GitHub**: [plasmatrip](https://github.com/plasmatrip)
