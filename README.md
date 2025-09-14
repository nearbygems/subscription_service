# Subscription Service

Сервис для хранения и подсчёта пользовательских онлайн-подписок.  
Реализован на **Go + PostgreSQL**, запускается через **Docker Compose**.

---

## 🚀 Запуск

1. Запустить сервис:

```bash
docker compose up --build
```
2. Приложение будет доступно по адресу:

```arduino
http://localhost:8080
```

## 📦 API

### Создать подписку
```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```
### Получить список подписок
```bash
curl http://localhost:8080/subscriptions
```

### Получить подписку по ID
```bash
curl http://localhost:8080/subscriptions/{id}
```

### Обновить подписку
```bash
curl -X PUT http://localhost:8080/subscriptions/{id} \
  -H "Content-Type: application/json" \
  -d '{"service_name": "Spotify", "price": 600}'
```

### Удалить подписку
```bash
curl -X DELETE http://localhost:8080/subscriptions/{id}
```

### Суммарная стоимость за период
```bash
curl "http://localhost:8080/subscriptions/summary?period_from=07-2025&period_to=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus"
```

## 📑 Swagger

Документация API доступна по адресу:
```bash
http://localhost:8080/swagger
```

## 🛠 Технологии
- Go 1.21 
- PostgreSQL 15 
- Docker + Docker Compose 
- Chi (router), SQLX, Logrus 
- Swagger (OpenAPI)