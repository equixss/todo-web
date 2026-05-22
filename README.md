# Todo Web API

Go веб-приложение для управления задачами (Todo) с использованием PostgreSQL и JWT аутентификации.

## Архитектура проекта

Проект использует feature-based clean architecture:

```
internal/feature/
├── users/          # Модуль пользователей
│   ├── repository/     # Репозитории данных
│   ├── service/        # Бизнес-логика
│   └── transport/http/ # HTTP handlers и DTOs
├── tasks/          # Модуль задач
│   ├── repository/     # Репозитории данных
│   ├── service/        # Бизнес-логика
│   └── transport/http/ # HTTP handlers и DTOs
└── statistics/     # Модуль статистики
    ├── repository/     # Репозитории данных
    ├── service/        # Бизнес-логика
    └── transport/http/ # HTTP handlers и DTOs

internal/core/
├── auth/           # JWT аутентификация
├── config/         # Конфигурация приложения
├── domain/         # Domain модели (User, Task, Statistics)
├── errors/         # Обработка ошибок
├── logger/         # Логирование
├── repository/     # Базовые репозитории
└── transport/http/ # Общие HTTP компоненты (middleware, server)
```

## Установка и запуск

### Требования

- Go 1.26.3+
- PostgreSQL

### Настройка окружения

1. Скопируйте `.env.example` в `.env`:
   ```bash
   cp .env.example .env
   ```

2. Заполните переменные окружения:
   ```
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=todoapp
   JWT_SECRET=your_secure_random_secret_key_here
   JWT_EXPIRY=24h
   JWT_ISSUER=todoweb-api
   ```

3. Запустите PostgreSQL контейнер:
   ```bash
   make env_up
   ```

4. Примените миграции базы данных:
   ```bash
   make migrate-up
   ```

5. Запустите приложение:
   ```bash
   make todo-run
   ```

Приложение будет доступно по адресу `http://localhost:5050`

## API Документация

### Swagger UI

Откройте Swagger UI для интерактивной документации:
- URL: http://localhost:5050/swagger/index.html
- JSON спецификация: http://localhost:5050/swagger.json

### Endpoints

#### Пользователи (Users)

| Метод | Endpoint | Описание | Авторизация |
|-------|----------|----------|-------------|
| POST | `/users` | Создание нового пользователя | Нет |
| POST | `/users/login` | Вход в систему | Нет |
| POST | `/users/refresh` | Обновление токена доступа | Да |
| POST | `/users/logout` | Выход из системы | Да |
| GET | `/users/me` | Получение текущего пользователя | Да |
| GET | `/users` | Список всех пользователей | Да |
| GET | `/users/{id}` | Получить пользователя по ID | Да |
| PATCH | `/users/{id}` | Обновление пользователя | Да |
| DELETE | `/users/{id}` | Удаление пользователя | Да |

#### Задачи (Tasks)

| Метод | Endpoint | Описание | Авторизация |
|-------|----------|----------|-------------|
| POST | `/tasks` | Создание новой задачи | Да |
| GET | `/tasks` | Список задач пользователя | Да |
| GET | `/tasks/{id}` | Получить задачу по ID | Да |
| PATCH | `/tasks/{id}` | Обновление задачи | Да |
| DELETE | `/tasks/{id}` | Удаление задачи | Да |

#### Статистика (Statistics)

| Метод | Endpoint | Описание | Авторизация |
|-------|----------|----------|-------------|
| GET | `/statistics` | Получение статистики задач | Да |

## Команды Makefile

```bash
make env_up          # Запуск PostgreSQL контейнера
make env_down        # Остановка PostgreSQL контейнера
make migrate-up      # Применение миграций базы данных
make migrate-down    # Откат миграций
make migrate-create  # Создание новой миграции
make todo-run        # Локальный запуск приложения
make todo-deploy     # Сборка и запуск через docker-compose
make ps              # Отображение статуса контейнеров
make env_clean       # Очистка volumes (с подтверждением)
make env-port-forward  # Port forwarding PostgreSQL на localhost:5432
make swag            #генерация сваггера
```

---
