
docs/api.md
```markdown
# API

Этот документ описывает основные HTTP-эндпоинты проекта **Resume Builder**.

---

## 1. Внешний API (backend)

### 1.1. `POST /api/v1/resume/pdf`

**Назначение:**  
Сгенерировать PDF-резюме по данным, полученным от фронтенда.

**Запрос:**

- Метод: `POST`
- Путь: `/api/v1/resume/pdf`
- Заголовки:
  - `Content-Type: application/json`
- Тело: JSON с моделью `ResumeRequest`:

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
      { "label": "GitHub", "url": "https://github.com/johndoe" }
    ]
  },
  "skills": ["Go", "Kubernetes"],
  "experience": [],
  "education": [],
  "customSections": [],
  "photo": {
    "mimeType": "image/jpeg",
    "data": "BASE64_ENCODED_JPEG_DATA"
  }
}
