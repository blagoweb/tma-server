# Telegram Mini App REST API

Простой REST-сервер на Go для Telegram Mini App с поддержкой PostgreSQL, JWT аутентификации и CORS.

## Возможности

- 🔐 JWT аутентификация для Telegram Mini App
- 🗄️ PostgreSQL база данных
- 🌐 CORS поддержка
- 📱 REST API для работы с данными
- 🐳 Docker поддержка
- ☁️ Готов к развертыванию на Railway

## Структура проекта

```
tma/
├── auth/           # JWT аутентификация
├── config/         # Конфигурация
├── database/       # Работа с БД
├── handlers/       # HTTP обработчики
├── middleware/     # Middleware
├── models/         # Модели данных
├── routes/         # Маршруты API
├── main.go         # Главный файл
├── go.mod          # Зависимости Go
├── Dockerfile      # Docker образ
└── README.md       # Документация
```

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/telegram` - Аутентификация через Telegram

### Пользователи (требует JWT)
- `GET /api/v1/user/profile` - Получить профиль пользователя

### Элементы (требует JWT)
- `GET /api/v1/items` - Получить все элементы пользователя
- `GET /api/v1/items/:id` - Получить элемент по ID
- `POST /api/v1/items` - Создать новый элемент
- `PUT /api/v1/items/:id` - Обновить элемент
- `DELETE /api/v1/items/:id` - Удалить элемент

### Система
- `GET /health` - Проверка состояния сервера

## Локальная разработка

### Требования
- Go 1.21+
- PostgreSQL

### Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd tma
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл `.env` на основе `env.example`:
```bash
cp env.example .env
```

4. Настройте переменные окружения в `.env`:
```env
DATABASE_URL=postgres://username:password@localhost:5432/tma_db?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
PORT=8080
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
ENV=development
```

5. Запустите PostgreSQL и создайте базу данных:
```sql
CREATE DATABASE tma_db;
```

6. Запустите сервер:
```bash
go run main.go
```

Сервер будет доступен по адресу `http://localhost:8080`

## Развертывание на Railway

### 1. Подготовка

1. Убедитесь, что у вас есть аккаунт на [Railway](https://railway.app/)
2. Установите Railway CLI:
```bash
npm install -g @railway/cli
```

### 2. Развертывание

1. Войдите в Railway:
```bash
railway login
```

2. Инициализируйте проект:
```bash
railway init
```

3. Добавьте PostgreSQL сервис:
```bash
railway add
# Выберите PostgreSQL
```

4. Настройте переменные окружения:
```bash
railway variables set JWT_SECRET=your-super-secret-jwt-key
railway variables set TELEGRAM_BOT_TOKEN=your-telegram-bot-token
railway variables set ENV=production
```

5. Разверните приложение:
```bash
railway up
```

### 3. Получение URL

После развертывания Railway предоставит URL вашего приложения. Используйте его в настройках Telegram Mini App.

## Использование API

### Аутентификация

1. Получите `initData` из Telegram Mini App
2. Отправьте POST запрос на `/api/v1/auth/telegram`:

```javascript
const response = await fetch('https://your-app.railway.app/api/v1/auth/telegram', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    init_data: initData
  })
});

const { token, user } = await response.json();
```

3. Используйте полученный токен в заголовке `Authorization`:

```javascript
const response = await fetch('https://your-app.railway.app/api/v1/items', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

### Примеры запросов

#### Создание элемента
```bash
curl -X POST https://your-app.railway.app/api/v1/items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Новая задача",
    "description": "Описание задачи"
  }'
```

#### Получение элементов
```bash
curl -X GET https://your-app.railway.app/api/v1/items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Безопасность

- Все JWT токены имеют срок действия 24 часа
- Валидация Telegram init data для предотвращения подделки
- CORS настроен для работы с Telegram Mini App
- Все запросы к защищенным эндпоинтам требуют валидный JWT токен

## Переменные окружения

| Переменная | Описание | Обязательная |
|------------|----------|--------------|
| `DATABASE_URL` | URL подключения к PostgreSQL | Да |
| `JWT_SECRET` | Секретный ключ для JWT | Да |
| `PORT` | Порт сервера | Нет (по умолчанию 8080) |
| `TELEGRAM_BOT_TOKEN` | Токен Telegram бота | Нет |
| `ENV` | Окружение (development/production) | Нет |

## Лицензия

MIT
