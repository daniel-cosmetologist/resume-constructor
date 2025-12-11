
# Resume Builder

Resume Builder — это веб-приложение для конструирования резюме и генерации PDF на основе LaTeX-шаблона.  
Фронтенд написан на Vue + TypeScript, бэкенд — на Go, генерация PDF вынесена в отдельный LaTeX-сервис.  
Все компоненты запускаются через `docker-compose` и работают в виде отдельных контейнеров.

---

## Содержание

1. [Обзор проекта](#обзор-проекта)  
2. [Структура репозитория](#структура-репозитория)  
3. [Сервисы и архитектура](#сервисы-и-архитектура)  
4. [Запуск в режиме разработки](#запуск-в-режиме-разработки)  
5. [Деплой в продакшн-среду](#деплой-в-продакшн-среду)  
6. [Техническое задание (MVP)](#техническое-задание-mvp)  
   - [1. Общая информация](#11-общая-информация)  
   - [2. Архитектура и компоненты](#12-архитектура-и-компоненты)  
   - [3. Нефункциональные требования](#13-нефункциональные-требования)  
   - [4. Модель данных резюме](#14-модель-данных-резюме)  
   - [5. Внешний API Backend](#15-api-backend-внешний-api-для-фронтенда)  
   - [6. Внутренний API LaTeX-сервиса](#16-внутренний-api-latex-сервиса)  
   - [7. LaTeX-шаблон и плейсхолдеры](#17-latex-шаблон-и-плейсхолдеры)  
   - [8. Фронтенд (Vue + TypeScript)](#18-фронтенд-vue--typescript)  
   - [9. Логирование и мониторинг](#19-логирование-и-мониторинг)  
   - [10. Будущее расширение](#110-будущее-расширение)  

---

## Обзор проекта

**Цель проекта:** предоставить простое веб-приложение, в котором пользователь может:

- заполнить данные резюме на одной странице;
- добавить кастомные секции (например, описание homelab или личных проектов);
- опционально загрузить фотографию;
- увидеть предпросмотр PDF прямо на странице;
- скачать готовый PDF-файл резюме.

Вся логика генерации PDF реализована через LaTeX-шаблон, в который подставляются данные пользователя.

**Язык интерфейса и контента:** английский (UI и содержимое резюме).  
**Аутентификация:** отсутствует в MVP (приложение публичное, без регистрации).

---

## Структура репозитория

Репозиторий организован как монорепозиторий и содержит все сервисы и инфраструктуру.

```text
resume-builder/
├── docker-compose.yml          # единая точка запуска всех сервисов
├── README.md                   # текущее описание проекта и ТЗ
├── .gitignore
├── .github/
│   └── workflows/
│       └── ci.yml              # CI-конфигурация
├── docs/
│   ├── requirements.md         # подробное ТЗ (функциональные/нефункциональные требования)
│   ├── architecture.md         # описание архитектуры и взаимодействия сервисов
│   └── api.md                  # описание API
├── frontend/
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── Dockerfile
│   └── src/
│       ├── main.ts
│       ├── App.vue
│       ├── router/
│       │   └── index.ts
│       ├── views/
│       │   └── ResumeBuilderView.vue
│       ├── components/
│       │   ├── PersonalInfoForm.vue
│       │   ├── SummaryForm.vue
│       │   ├── ContactsForm.vue
│       │   ├── SkillsForm.vue
│       │   ├── ExperienceForm.vue
│       │   ├── EducationForm.vue
│       │   ├── CustomSectionsForm.vue
│       │   ├── PhotoUpload.vue
│       │   └── PdfPreview.vue
│       ├── api/
│       │   └── resumeApi.ts
│       └── types/
│           └── resume.ts
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   └── internal/
│       ├── config/
│       │   └── config.go
│       ├── http/
│       │   ├── router.go
│       │   ├── handlers.go
│       │   └── middleware.go
│       ├── resume/
│       │   ├── model.go
│       │   ├── validation.go
│       │   └── service.go
│       └── latexclient/
│           └── client.go
├── latex-service/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── cmd/
│   │   └── latexservice/
│   │       └── main.go
│   ├── templates/
│   │   └── resume_template.tex
│   └── internal/
│       ├── config/
│       │   └── config.go
│       ├── http/
│       │   ├── router.go
│       │   └── handlers.go
│       ├── latex/
│       │   ├── renderer.go
│       │   └── sanitizer.go
│       └── model/
│           └── resume.go
└── gateway-nginx/
    ├── Dockerfile
    └── nginx.conf
````

---

## Сервисы и архитектура

Проект состоит из четырёх основных сервисов:

1. **Frontend (Vue + TypeScript)**

   * SPA, собирается в статический бандл.
   * Работает за nginx (gateway).
   * Общается с backend API по HTTP (REST, JSON).
   * Реализует форму, предпросмотр PDF и скачивание.

2. **Backend API (Go)**

   * REST API.
   * Принимает JSON-данные резюме от фронтенда.
   * Валидирует входные данные.
   * Делегирует генерацию PDF LaTeX-сервису.
   * Возвращает PDF фронтенду и не сохраняет его на диск/в БД.

3. **LaTeX/PDF Service (Go + LaTeX)**

   * Отдельный контейнер с установленным LaTeX (`pdflatex`, `latexmk`).
   * Хранит LaTeX-шаблон с плейсхолдерами `{{...}}`.
   * Принимает JSON с данными резюме + фото (base64).
   * Обрабатывает фото (сжатие, кадрирование до 3×4, ограничение 2 МБ).
   * Подставляет данные в шаблон, компилирует PDF и возвращает его API.

4. **Gateway nginx**

   * Принимает внешний HTTP-трафик.
   * Проксирует статический фронтенд и запросы к backend API.
   * В будущем может быть заменён или обёрнут ingress-nginx / API Gateway при деплое в Kubernetes.

В текущем MVP:

* нет базы данных;
* PDF и временные файлы хранятся только в оперативной памяти и временных директориях на время обработки запроса;
* в перспективе возможна интеграция с MinIO и Postgres.

---

## Запуск в режиме разработки

1. Установить Docker и docker-compose.

2. Клонировать репозиторий:

   ```bash
   git clone https://github.com/<your-username>/resume-builder.git
   cd resume-builder
   ```

3. Запустить все сервисы:

   ```bash
   docker-compose up --build
   ```

4. Открыть браузер и перейти по адресу:

   ```text
   http://localhost:8080
   ```

   (порт 8080 проброшен из контейнера `gateway` в `docker-compose.yml`).

5. Разработка:

   * изменять код фронтенда/бэкенда/latex-сервиса;
   * по необходимости перезапускать отдельные сервисы через `docker-compose restart <service>`.

---

## Деплой в продакшн-среду

Базовый сценарий деплоя на выделенную виртуальную машину:

1. Поднять виртуальную машину (Linux).

2. Установить Docker и docker-compose.

3. Склонировать репозиторий на сервер.

4. При необходимости настроить переменные окружения (`HTTP_ADDR`, `LATEX_SERVICE_URL`, `TEMPLATE_PATH` и т.п.).

5. Запустить:

   ```bash
   docker-compose up -d --build
   ```

6. Настроить внешний ingress / балансировщик, если требуется:

   * HTTPS-терминация;
   * маршрутизация по доменам;
   * интеграция с Kubernetes (на следующем этапе развития).

---

## Техническое задание (MVP)

Ниже приведён консолидированный вариант ТЗ, использованный при проектировании.

### 1.1. Общая информация

**Название проекта (рабочее):** `resume-builder`

**Назначение:**
Веб-приложение для интерактивной генерации резюме в формате PDF на основе LaTeX-шаблона:
пользователь заполняет форму, добавляет кастомные блоки, опционально загружает фото и получает PDF.

**Основной пользовательский сценарий:**

1. Пользователь открывает веб-приложение.
2. Заполняет форму (на английском языке).
3. При необходимости добавляет кастомные секции (например, описание homelab или проектов).
4. Опционально загружает фото.
5. Видит предпросмотр PDF-резюме на странице.
6. Нажимает кнопку для скачивания и получает PDF.

**Язык интерфейса и контента:** английский.
**Регистрация, аутентификация, авторизация:** отсутствуют в MVP.

---

### 1.2. Архитектура и компоненты

Состав системы:

* **Frontend (Vue + TS)** — UI, форма, предпросмотр PDF.
* **Backend API (Go)** — REST API, валидация, вызов LaTeX-сервиса.
* **LaTeX Service (Go + LaTeX)** — работа с шаблоном, генерация PDF.
* **Gateway nginx** — reverse proxy для фронта и API.
* **Инфраструктура** — `docker-compose` + монорепозиторий.

Хранение данных в MVP:

* PDF и временные файлы создаются только в процессе обработки запроса и удаляются после ответа;
* отсутствие постоянного хранения (БД, объектное хранилище);
* на будущее — MinIO (PDF/картинки) и Postgres.

---

### 1.3. Нефункциональные требования

* Время генерации PDF: ориентир до 5 секунд при типовом объёме данных.
* Предпросмотр должен использовать дебаунс (не на каждое нажатие).
* Ограничение размера JSON (ориентировочно до 256 КБ).
* Фото:

  * если размер > 2 МБ — автоматическое сжатие до ≤ 2 МБ;
  * кадрирование до соотношения сторон 3×4 (портрет).
* Ограничения по длинам полей:

  * `fullName` — до 100 символов;
  * `position` — до 100;
  * `summary` — до 1500;
  * `skills` — до 50 элементов, каждый до 50 символов;
  * bullet — до 300 символов;
  * `experience` / `education` / `customSections` — до 10 записей.
* Валидация email и URL.
* Экранирование LaTeX-символов перед генерацией `.tex`.

---

### 1.4. Модель данных резюме

JSON-модель, используемая во внешнем API:

```jsonc
{
  "fullName": "John Doe",
  "position": "Senior Software Engineer",
  "summary": "Experienced engineer...",
  "contacts": {
    "email": "john.doe@example.com",
    "phone": "+1 555 123 4567",
    "location": "San Francisco, CA, USA",
    "links": [
      { "label": "GitHub", "url": "https://github.com/johndoe" },
      { "label": "LinkedIn", "url": "https://www.linkedin.com/in/johndoe" }
    ]
  },
  "skills": ["Go", "Kubernetes"],
  "experience": [
    {
      "company": "Example Corp",
      "position": "Senior Backend Engineer",
      "location": "Remote",
      "startDate": "2020-05",
      "endDate": "",
      "description": "Working on high-load APIs...",
      "bullets": [
        "Designed and implemented a scalable microservice architecture."
      ]
    }
  ],
  "education": [
    {
      "institution": "Tech University",
      "degree": "BSc in Computer Science",
      "location": "City, Country",
      "startDate": "2012-09",
      "endDate": "2016-06",
      "details": "Graduated with honors."
    }
  ],
  "customSections": [
    {
      "title": "Homelab",
      "bulletSymbol": "•",
      "items": [
        "Built a Kubernetes cluster at home.",
        "Self-hosted monitoring stack."
      ]
    }
  ],
  "photo": {
    "mimeType": "image/jpeg",
    "data": "BASE64_ENCODED_JPEG_DATA"
  }
}
```

---

### 1.5. API Backend (внешний API для фронтенда)

**Endpoint:** `POST /api/v1/resume/pdf`

* Вход: JSON по модели выше.

* Выход (успех): `200 OK`, `Content-Type: application/pdf`, бинарный PDF.

* Ошибка валидации: `400 Bad Request`, JSON вида:

  ```json
  {
    "error": "validation_error",
    "message": "Invalid resume data",
    "details": {
      "fullName": "Full name is required"
    }
  }
  ```

* Ошибка генерации PDF: `500 Internal Server Error`, JSON:

  ```json
  {
    "error": "pdf_generation_failed",
    "message": "Failed to generate PDF"
    }
  ```

Тот же endpoint используется и для предпросмотра: фронтенд получает PDF как `blob` и встраивает его в `<object>`.

---

### 1.6. Внутренний API LaTeX-сервиса

**Endpoint:** `POST /internal/v1/render`

* Доступен только во внутренней сети Docker.
* Принимает JSON с данными резюме (структура аналогична внешнему API).
* Возвращает `application/pdf` при успехе.
* При ошибках возвращает JSON с кодом/сообщением.

---

### 1.7. LaTeX-шаблон и плейсхолдеры

В `latex-service/templates/resume_template.tex` используется шаблон с плейсхолдерами:

* `{{FullName}}`
* `{{Position}}`
* `{{Summary}}`
* `{{Contacts}}`
* `{{Skills}}`
* `{{Experience}}`
* `{{Education}}`
* `{{CustomSections}}`
* `{{Photo}}`

LaTeX-сервис подставляет в них уже готовые LaTeX-фрагменты, формируемые по структурам данных, с предварительным экранированием спецсимволов.

---

### 1.8. Фронтенд (Vue + TypeScript)

Основная страница `ResumeBuilderView` содержит:

* блоки форм: Personal Info, Summary, Contacts, Skills, Experience, Education;
* Custom Sections с настраиваемым символом bullet и списком элементов;
* Photo Upload c конвертацией файла в base64;
* PdfPreview с live-предпросмотром;
* кнопку Download PDF.

Фронтенд:

* хранит состояние резюме в виде объекта `ResumeRequest`;
* наблюдает за изменениями и с дебаунсом вызывает backend для обновления PDF;
* отображает ошибки API пользователю.

---

### 1.9. Логирование и мониторинг

* **Backend**: логирует входящие запросы, ошибки валидации, ошибки вызовов LaTeX-сервиса.
* **LaTeX-service**: логирует запросы на рендер и ошибки компиляции LaTeX.
* **nginx-gateway**: ведёт access/error логи.

---

### 1.10. Будущее расширение

* Подключение MinIO для хранения PDF и фотографий.
* Подключение Postgres для пользователей и резюме.
* Несколько шаблонов резюме и выбор шаблона на фронтенде.
* Мультиязычность (интерфейс и содержимое резюме).
* Перевод деплоя в Kubernetes:

  * отдельные `Deployment`/`Service` на каждый компонент;
  * `ingress-nginx` / API Gateway для HTTPS и маршрутизации.

---

Этот `README.md` является основным источником правды по структуре проекта и ТЗ для MVP. Подробности по архитектуре и API дополнительно раскрыты в `docs/architecture.md` и `docs/api.md`.
