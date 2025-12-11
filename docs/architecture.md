# Architecture

Этот документ описывает архитектуру проекта **Resume Builder** на уровне сервисов и их взаимодействия.

## 1. Общий обзор

Система состоит из четырёх контейнеров:

1. **frontend** — SPA на Vue + TypeScript.
2. **backend** — HTTP API на Go.
3. **latex-service** — сервис генерации PDF через LaTeX.
4. **gateway** — nginx, выступающий reverse proxy.

Связи:

- браузер ↔ `gateway` (HTTP, порт 8080 на хосте);
- `gateway` ↔ `frontend` (статические файлы);
- `gateway` ↔ `backend` (REST API `/api/...`);
- `backend` ↔ `latex-service` (внутренний HTTP `/internal/v1/render`).

Все контейнеры находятся в одной сети docker-compose.

## 2. Поток запроса

1. Пользователь открывает `http://localhost:8080`.
2. `gateway` отдаёт статический фронтенд из контейнера `frontend`.
3. Пользователь заполняет форму резюме, фронтенд собирает данные в объект `ResumeRequest`.
4. Для предпросмотра и скачивания фронтенд вызывает:

   ```http
   POST /api/v1/resume/pdf
   Host: <gateway>
   Content-Type: application/json
