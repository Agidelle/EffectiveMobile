# EffectiveMobile

REST сервис для агрегации данных об онлайн-подписках пользователей.


## Структура проекта

- `cmd/` — точка входа, команды миграций и запуска сервера
- `internal/` — бизнес-логика, API, сервисы, хранилища
- `api/auth` — авторизация и аутентификация для примера, не интегрирована
- `docs/` — документация API (Swagger)
- `migrations/` — SQL-миграции для базы данных
- `entrypoint.sh` — скрипт запуска миграций и сервера
- `Dockerfile` — сборка контейнера приложения
- `docker-compose.yml` — запуск приложения и базы данных

## Быстрый старт

1. **Сборка и запуск сервиса:**
   ```sh
   docker compose up --build
   
2. Миграции применяются автоматически при запуске.

## Управление сервисом

Команды:\
Запуск сервиса: serve\
Запуск миграций: migration up\
Откат миграций: migration down

## Управлением миграциями

- Применить миграции вручную:
```sh
docker compose run --rm subs-app ./SUBS migration up
```
- Откатить миграции (очистить БД):
```sh
./SUBS migration down
```

## Переменные конфигурации
Используются переменные для подключения к БД (см. .env):

DB_HOST=postgres\
DB_PORT=5432\
DB_NAME=mydatabase\
DB_USER=user\
DB_PASSWORD=mysecretpassword\
APP_PORT=3000

## Документация
Swagger-описание API находится в docs/swagger.json/yaml.  <hr></hr> 
## Технологии
- Go 1.25+
- PostgreSQL
- Chi (HTTP роутер)
- Viper (конфигурации)
- Cobra (CLI)
- Testify (тесты)
- Slog (логирование)
- migrate (миграции БД)
- Docker, Docker Compose
- Swagger (документация API)