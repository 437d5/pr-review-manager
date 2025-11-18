# PR Review Manager

Система для управления процессом ревью pull request'ов в командах разработки.

## 📦 Установка и запуск

1. **Склонировать репозиторий**
```bash
git clone https://github.com/437d5/pr-review-manager.git
```

2. **Запуск приложения**
```
make up
```

Можно запустить через `docker compose up`, но необходимо будет создать `.env` файл
1. **Создать `.env` файл**
```bash
cat > .env << EOF
MODE=dev
# MODE can be dev | prod

REVIEWER_ADDRESS=:8080
REVIEWER_READ_TIMEOUT=15
REVIEWER_WRITE_TIMEOUT=15
REVIEWER_IDLE_TIMEOUT=60

DB_NAME=pr_reviewer
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASS=password
EOF
```
2. **Запуск приложения**
```bash
docker compose up
```

## Envs
В директории env находится файл пример `.env.example`.

В docker-compose.yaml и Makefile прописан env/.env.example файл, 
который уже есть в репозитории, поэтому его можно использовать готовый.
С этим файлом сервер запустится на `http://127.0.0.1:8080`.

```bash
MODE=dev # Настраивает уровень логгирования может быть
# или dev или prod, если dev, включаются отладочные логи.

REVIEWER_ADDRESS=:8080 # Порт для запуска сервера
# Таймауты
REVIEWER_READ_TIMEOUT=15 
REVIEWER_WRITE_TIMEOUT=15
REVIEWER_IDLE_TIMEOUT=60

# Переменные для подключения к бд
DB_NAME=pr_reviewer
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASS=password
```

### Юнит-тесты

```bash
go test -v ./...
```
