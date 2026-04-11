# Wishlist API

Сервис для управления вишлистами. Пользователь может зарегистрироваться, создать вишлист к празднику, наполнить его подарками и открыть доступ по уникальной ссылке.


## Запуск

```bash
docker-compose up --build
```

Сервис запустится на `http://localhost:8080`.

## Конфигурация

Все параметры передаются через переменные окружения. Пример в файле `.env.example`

## API

### Авторизация

**Регистрация**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@mail.ru", "password": "secret123"}'
```
```json
{"message": "user registered successfully"}
```

**Вход**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@mail.ru", "password": "password123"}'
```
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400
}
```

Полученный `access_token` передавай в заголовке `Authorization: Bearer <token>` для всех закрытых эндпоинтов. Действие `access_token` 24 часа, потом он утрачивает свою силу.
Правильнее было бы сделать его действующим не более часа, и восстанавливать по `refresh_token`, но в рамках этого задания мне показалось это излишним.

---

### Вишлисты

**Создать вишлист**
```bash
curl -X POST http://localhost:8080/wishlist \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "День рождения",
    "description": "Мой вишлист",
    "event_date": "2025-06-01T00:00:00Z"
  }'
```
```json
{
  "id": "uuid",
  "title": "День рождения",
  "description": "Мой вишлист",
  "event_date": "2025-06-01T00:00:00Z",
  "share_token": "abc123xyz",
  "created_at": "2025-04-01T12:00:00Z",
  "updated_at": "2025-04-01T12:00:00Z"
}
```

**Список своих вишлистов**
```bash
curl http://localhost:8080/wishlist \
  -H "Authorization: Bearer <token>"
```

**Получить вишлист по ID**
```bash
curl http://localhost:8080/wishlist/{id} \
  -H "Authorization: Bearer <token>"
```

**Обновить вишлист** (все поля опциональны)
```bash
curl -X PUT http://localhost:8080/wishlist/{id} \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Новое название"}'
```

**Удалить вишлист**
```bash
curl -X DELETE http://localhost:8080/wishlist/{id} \
  -H "Authorization: Bearer <token>"
```

---

### Подарки

**Добавить подарок**
```bash
curl -X POST http://localhost:8080/wishlist/{id}/items \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "PS5",
    "description": "Белая",
    "product_url": "https://shop.com/ps5",
    "priority": 10
  }'
```
```json
{
  "id": "uuid",
  "title": "PS5",
  "description": "Белая",
  "product_url": "https://shop.com/ps5",
  "priority": 10,
  "is_reserved": false,
  "created_at": "2025-04-01T12:00:00Z",
  "updated_at": "2025-04-01T12:00:00Z"
}
```

**Список подарков**
```bash
curl http://localhost:8080/wishlist/{id}/items \
  -H "Authorization: Bearer <token>"
```

**Обновить подарок** (все поля опциональны)
```bash
curl -X PUT http://localhost:8080/wishlist/{id}/items/{itemId} \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "PS5 Pro", "priority": 9}'
```

**Удалить подарок**
```bash
curl -X DELETE http://localhost:8080/wishlist/{id}/items/{itemId} \
  -H "Authorization: Bearer <token>"
```

---

### Публичный доступ

Эти эндпоинты доступны без авторизации — по уникальному токену.

**Просмотр вишлиста по share_token**
```bash
curl http://localhost:8080/public/{share_token}
```
```json
{
  "id": "uuid",
  "title": "День рождения",
  "description": "Мой вишлист",
  "event_date": "2025-06-01T00:00:00Z",
  "items": [
    {
      "id": "uuid",
      "title": "PS5",
      "priority": 10,
      "is_reserved": false
    }
  ]
}
```

**Забронировать подарок**
```bash
curl -X POST http://localhost:8080/public/{share_token}/items/{itemId}/reserve
```
```json
{"message": "item reserved successfully"}
```

Если подарок уже забронирован — вернёт `409 Conflict`.
