# Инструкции по развертыванию на Railway

## Быстрое развертывание

### 1. Подготовка репозитория

1. Создайте новый репозиторий на GitHub
2. Загрузите код в репозиторий:
```bash
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/yourusername/tma.git
git push -u origin main
```

### 2. Развертывание на Railway

1. Перейдите на [Railway](https://railway.app/)
2. Нажмите "New Project"
3. Выберите "Deploy from GitHub repo"
4. Выберите ваш репозиторий
5. Railway автоматически определит, что это Go проект и начнет сборку

### 3. Настройка базы данных

1. В проекте Railway нажмите "New"
2. Выберите "Database" → "PostgreSQL"
3. Railway автоматически создаст базу данных и установит переменную `DATABASE_URL`

### 4. Настройка переменных окружения

В настройках проекта Railway добавьте следующие переменные:

| Переменная | Значение | Описание |
|------------|----------|----------|
| `JWT_SECRET` | `bla` | Секретный ключ для JWT (сгенерируйте случайную строку) |
| `TELEGRAM_BOT_TOKEN` | `your-bot-token` | Токен вашего Telegram бота (опционально) |
| `ENV` | `production` | Окружение |
| `PORT` | `8080` | Порт (Railway автоматически переопределит) |

### 5. Получение URL

После развертывания Railway предоставит URL вида:
```
https://your-app-name.railway.app
```

## Настройка Telegram Mini App

### 1. Создание бота

1. Найдите @BotFather в Telegram
2. Отправьте `/newbot`
3. Следуйте инструкциям для создания бота
4. Сохраните токен бота

### 2. Создание Mini App

1. Отправьте @BotFather `/newapp`
2. Выберите вашего бота
3. Укажите название и описание приложения
4. Укажите URL вашего приложения на Railway
5. Получите токен Mini App

### 3. Настройка веб-приложения

Создайте веб-приложение и разместите его на хостинге (например, GitHub Pages, Vercel, Netlify).

Пример простого веб-приложения:

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Telegram Mini App</title>
    <script src="https://telegram.org/js/telegram-web-app.js"></script>
</head>
<body>
    <div id="app">
        <h1>Welcome to My App!</h1>
        <button onclick="initApp()">Initialize App</button>
    </div>

    <script>
        let tg = window.Telegram.WebApp;
        tg.ready();

        async function initApp() {
            try {
                const response = await fetch('https://your-app.railway.app/api/v1/auth/telegram', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        init_data: tg.initData
                    })
                });

                const data = await response.json();
                console.log('Authenticated:', data);
                
                // Теперь можно использовать API с токеном
                localStorage.setItem('token', data.token);
                
            } catch (error) {
                console.error('Auth error:', error);
            }
        }
    </script>
</body>
</html>
```

## Проверка развертывания

### 1. Health Check

Откройте в браузере:
```
https://your-app.railway.app/health
```

Должен вернуться ответ:
```json
{
  "status": "ok",
  "message": "Telegram Mini App API is running"
}
```

### 2. Тестирование API

Используйте файл `test_api.http` для тестирования API через REST Client в VS Code или Postman.

### 3. Мониторинг

В Railway Dashboard вы можете:
- Просматривать логи приложения
- Мониторить использование ресурсов
- Настраивать автоматическое масштабирование

## Обновление приложения

Для обновления приложения:

1. Внесите изменения в код
2. Загрузите изменения в GitHub:
```bash
git add .
git commit -m "Update description"
git push
```

3. Railway автоматически пересоберет и развернет приложение

## Troubleshooting

### Проблемы с подключением к БД

1. Проверьте переменную `DATABASE_URL` в Railway
2. Убедитесь, что PostgreSQL сервис запущен
3. Проверьте логи приложения в Railway Dashboard

### Проблемы с CORS

1. Убедитесь, что CORS middleware подключен
2. Проверьте, что домен вашего веб-приложения разрешен

### Проблемы с JWT

1. Проверьте переменную `JWT_SECRET`
2. Убедитесь, что токен передается в заголовке `Authorization: Bearer <token>`

### Проблемы с Telegram Auth

1. Проверьте переменную `TELEGRAM_BOT_TOKEN`
2. Убедитесь, что `init_data` корректно передается
3. Проверьте логи для деталей ошибок

## Безопасность

### Рекомендации для продакшена

1. **JWT_SECRET**: Используйте длинную случайную строку (минимум 32 символа)
2. **HTTPS**: Railway автоматически предоставляет SSL сертификаты
3. **Rate Limiting**: Рассмотрите добавление rate limiting middleware
4. **Logging**: Настройте структурированное логирование
5. **Monitoring**: Настройте алерты в Railway

### Генерация JWT_SECRET

```bash
# Генерация случайного секрета
openssl rand -base64 32
```

## Стоимость

Railway предоставляет:
- **Free Tier**: $5 кредитов в месяц
- **Pro Plan**: $20/месяц за дополнительные ресурсы

Для небольшого Mini App обычно достаточно Free Tier.

## Поддержка

- [Railway Documentation](https://docs.railway.app/)
- [Telegram Mini Apps Documentation](https://core.telegram.org/bots/webapps)
- [Go Documentation](https://golang.org/doc/)
