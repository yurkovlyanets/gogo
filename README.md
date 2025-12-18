# gogo

Pet‑project на Go с REST API и PostgreSQL, построенный с опорой на принципы **Clean Architecture**.

Проект создаётся как учебный, но с максимально «продовой» структурой: тонкий `main`, composition root, изоляция инфраструктуры и аккуратная работа с конфигурацией.

---

## 🚀 Цели проекта

* Написать REST API на Go
* Сохранять и читать данные пользователя из PostgreSQL
* Использовать **Viper** и `config.yaml` для конфигурации
* Отделить домен и юзкейсы от инфраструктуры
* Соблюдать принципы Clean Architecture
* Получить удобную базу для дальнейшего расширения (auth, фронт, тесты)

---

## 🧱 Архитектура

### Ключевые идеи

* `main.go` — только запуск приложения и graceful shutdown
* Вся сборка зависимостей находится в **composition root**
* Конфигурация читается один раз при старте
* Домен и usecase не знают ни про HTTP, ни про PostgreSQL

### Composition Root

Роль composition root выполняет пакет:

```
/server
```

Он отвечает за:

* подключение к PostgreSQL
* сборку HTTP‑слоя
* запуск и корректную остановку сервера

---

### Предварительные требования

* Go **1.22+**
* Docker + Docker Compose
* CLI для миграций **golang-migrate** (команда `migrate`)

    * macOS (Homebrew): `brew install golang-migrate`

---
## Запуск:
### 1) Клонировать репозиторий

```bash
  git clone https://github.com/yurkovlyanets/gogo.git
cd gogo
```

---

### 2) Поднять PostgreSQL

```bash
  make db-up
```

Проверка (опционально):

```bash
  docker ps | grep gogo-postgres
```

---

### 3) Применить миграции

> команда использует переменную окружения `DATABASE_URL`.

Самый простой способ:

```bash

export DATABASE_URL="postgres://gogo:gogo@localhost:5432/gogo?sslmode=disable"
make migrate-up
```

Если миграции применились успешно — ошибок не будет.

---

### 4) Запустить

```bash

go run ./cmd/api
```

В логах должно появиться что-то вроде:

```
listening on http://localhost:8080
```

---

## 🧪 Проверка API (curl)

### 1) Healthcheck

**Request**

```bash

curl -i http://localhost:8080/health
```

**Response (200)**

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"status":"ok"}
```

---

### 2) Создать пользователя

**Request**

```bash

curl -i -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Lera","email":"valeriya@example.com"}'
```

**Response (201)**

```http
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8

{
  "id": "<uuid>",
  "name": "Lera",
  "email": "lera@example.com",
  "createdAt": "2025-12-18T20:00:00+05:00"
}
```
---

### 3) Получить пользователя по ID

Подставь `id` из ответа выше.

**Request**

```bash

curl -i http://localhost:8080/api/users/<uuid>
```

**Response (200)**

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "id": "<uuid>",
  "name": "Lera",
  "email": "lera@example.com",
  "createdAt": "2025-12-18T20:00:00+05:00"
}
```

---

### 4) Ошибки и типовые ответы

#### Неверный JSON при создании

**Request**

```bash

curl -i -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name":}'
```

**Response (400)**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8

{"error":"invalid json"}
```

#### Неверный UUID в GET

**Request**

```bash
curl -i http://localhost:8080/api/users/not-a-uuid
```

**Response (400)**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8

{"error":"invalid id"}
```

#### Пользователь не найден

**Request**

```bash
  curl -i http://localhost:8080/api/users/00000000-0000-0000-0000-000000000000
```

**Response (404)**

```http
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8

{"error":"user not found"}
```

## 🧩 Что можно сделать дальше

* добавить интеграционные тесты (testcontainers)
* добавить авторизацию
* добавить логирование и метрики

---

## 📌 Статус

Проект находится в активной разработке и используется как учебный пример перехода с Java на Go с правильной архитектурой.

