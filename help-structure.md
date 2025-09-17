# Go Production Best Practices Guide

## Когда НЕ использовать интерфейсы

В Go есть правило: **"Accept interfaces, return structs"**. Но в продакшене часто интерфейсы избыточны:

### ❌ Плохо: Преждевременная абстракция
```go
// Ненужный интерфейс для простого CRUD
type UserRepository interface {
    GetByID(id uuid.UUID) (*models.User, error)
    Create(user *models.User) error
}

type UserService interface {
    GetProfile(userID uuid.UUID) (*models.User, error)
}
```

### ✅ Хорошо: Прямые зависимости
```go
type UserService struct {
    userRepo *repository.UserRepository  // Конкретный тип
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}
```

## Структура продакшен проекта

```
internal/
├── app/             # Композиция зависимостей
├── config/          # Конфигурация приложения
├── database/        # Подключение к БД
├── dto/            # Data Transfer Objects
├── handlers/       # HTTP обработчики
├── middleware/     # Middleware для HTTP
├── models/         # Бизнес модели
├── repository/     # Слой данных (БД)
├── routes/         # Только маршрутизация
├── service/        # Бизнес логика
└── utils/          # Утилиты

cmd/
└── app/
    └── main.go     # Точка входа + инициализация
```

## Принципы инъекции зависимостей

### Правильная композиция зависимостей

```go
// internal/app/app.go - Слой композиции
type Dependencies struct {
    AuthHandler    *handlers.AuthHandler
    UserHandler    *handlers.UserHandler
    AuthMiddleware *middleware.AuthMiddleware
}

func NewDependencies(db *database.DB, cfg *config.Config) *Dependencies {
    // Репозитории
    userRepo := repository.NewUserRepository(db)
    authRepo := repository.NewAuthRepository(db)

    // Сервисы
    authService := service.NewAuthService(userRepo, authRepo)
    userService := service.NewUserService(userRepo)

    // Handlers и middleware
    authHandler := handlers.NewAuthHandler(authService, authRepo, cfg)
    userHandler := handlers.NewUserHandler(userService)
    authMiddleware := middleware.NewAuthMiddleware(authRepo)

    return &Dependencies{
        AuthHandler:    authHandler,
        UserHandler:    userHandler,
        AuthMiddleware: authMiddleware,
    }
}
```

```go
// cmd/devprep/main.go - Инициализация
func main() {
    // Настройка конфигурации и БД
    cfg := config.Load()
    db := database.Connect(cfg)
    fiberApp := fiber.New()

    // Создание всех зависимостей
    deps := app.NewDependencies(db, cfg)

    // Настройка маршрутов
    routes.SetupRoutes(fiberApp, deps)

    fiberApp.Listen(":3000")
}
```

```go
// internal/routes/routes.go - Только маршрутизация
func SetupRoutes(fiberApp *fiber.App, deps *app.Dependencies) {
    api := fiberApp.Group("/api/v1")

    SetupAuthRoutes(api, deps.AuthHandler, deps.AuthMiddleware)
    SetupUserRoutes(api, deps.UserHandler, deps.AuthMiddleware)
}
```

## Когда использовать интерфейсы

Интерфейсы нужны только при:

1. **Множественных реализациях**
```go
// Разные провайдеры уведомлений
type Notifier interface {
    Send(message string) error
}

type EmailNotifier struct{}
type SMSNotifier struct{}
type SlackNotifier struct{}
```

2. **Тестировании сложной логики**
```go
// Только если мокинг критичен
type PaymentGateway interface {
    ProcessPayment(amount int) error
}
```

3. **Стандартные интерфейсы Go**
```go
func HandleFile(w io.Writer, r io.Reader) error {
    // Используем стандартные интерфейсы
}
```

## Тестирование без интерфейсов

### Интеграционные тесты с реальной БД
```go
func TestUserService(t *testing.T) {
    db := setupTestDB(t) // SQLite в памяти
    defer db.Close()

    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)

    // Тестируем с реальной БД
    user, err := userService.GetProfile(userID)
    assert.NoError(t, err)
}
```

### Unit тесты - оставляем старые моки
```go
// Если у вас уже есть unit тесты с моками - оставьте их
// Проект может содержать legacy тесты с интерфейсами
// Главное - не создавать НОВЫЕ интерфейсы для тестирования

// Пример: тесты уже использующие моки можно не переписывать
type MockUserRepository struct {
    // старые моки
}
```

### Использование testcontainers для полноценных тестов
```go
func TestWithPostgres(t *testing.T) {
    ctx := context.Background()
    container, err := postgres.RunContainer(ctx)
    require.NoError(t, err)
    defer container.Terminate(ctx)

    // Подключаемся к реальной БД в контейнере
    // Более точное тестирование реального поведения
}
```

## Обработка ошибок в продакшене

### Стандартные ошибки вместо кастомных интерфейсов
```go
import (
    "database/sql"
    "errors"
)

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(...)

    if err == sql.ErrNoRows {
        return nil, sql.ErrNoRows  // Стандартная ошибка
    }

    return &user, err
}
```

## Конфигурация и запуск

### Graceful shutdown без лишних абстракций
```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    // Простая инициализация
    cfg := config.Load()
    db := database.Connect(cfg)
    app := setupApp(db, cfg)

    // Graceful shutdown
    go func() {
        app.Listen(":8080")
    }()

    <-ctx.Done()
    app.Shutdown()
}
```

## Мониторинг и логирование

### Структурированные логи
```go
import "log/slog"

func (s *AuthService) Login(req *LoginRequest) error {
    slog.Info("login attempt",
        "email", req.Email,
        "ip", req.IP,
    )

    // Логика авторизации

    slog.Info("login successful", "user_id", user.ID)
    return nil
}
```

### Health checks
```go
func (app *App) setupHealthChecks() {
    app.router.GET("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "db": app.db.Ping() == nil,
        })
    })
}
```

## Основные принципы

1. **KISS** - Keep It Simple, Stupid
2. **YAGNI** - You Aren't Gonna Need It
3. **Конкретные типы** вместо интерфейсов по умолчанию
4. **Composition over inheritance**
5. **Явные зависимости** в конструкторах
6. **Ошибки как значения**, а не исключения

Помните: в Go интерфейсы должны быть **маленькими** и **специфичными**. Большие интерфейсы - признак плохого дизайна.
