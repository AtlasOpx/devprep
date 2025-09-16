Структура как на проде

```go
├── cmd/
│   └── devprep/
│       └── main.go              ← ТОЧКА ВХОДА (запуск приложения)
├── internal/
│   ├── config/
│   │   └── config.go            ← НАСТРОЙКИ (порты, БД, секреты)
│   ├── database/
│   │   ├── connection.go        ← ПОДКЛЮЧЕНИЕ К БД
│   │   └── migrations.go        ← ЗАПУСК МИГРАЦИЙ
│   ├── dto/                     ← DTO = Data Transfer Objects
│   │   ├── user_dto.go          ← ЧТО ПРИХОДИТ/УХОДИТ В JSON
│   │   └── response_dto.go      ← СТАНДАРТНЫЕ ОТВЕТЫ API
│   ├── models/                  ← МОДЕЛИ = структуры для БД
│   │   └── user.go              ← ТАБЛИЦА users в виде Go структуры
│   ├── repository/              ← РЕПОЗИТОРИЙ = работа с БД
│   │   ├── interfaces/          ← ИНТЕРФЕЙСЫ для тестов
│   │   │   └── user_repo.go     ← ЧТО УМЕЕТ ДЕЛАТЬ репозиторий
│   │   └── user_repository.go   ← КАК работать с БД (SQL запросы)
│   ├── service/                 ← СЕРВИС = бизнес-логика
│   │   ├── interfaces/
│   │   │   └── user_service.go  ← ЧТО УМЕЕТ ДЕЛАТЬ сервис
│   │   └── user_service.go      ← БИЗНЕС-ЛОГИКА (валидация, хеширование и т.д.)
│   ├── handlers/                ← ХЕНДЛЕРЫ = обработка HTTP запросов
│   │   └── user_handler.go      ← ПОЛУЧАЕТ запрос → ОТДАЕТ ответ
│   ├── middleware/              ← МИДДЛВАРЫ = проверки перед хендлером
│   │   ├── auth.go              ← ПРОВЕРКА токенов
│   │   └── validation.go        ← ВАЛИДАЦИЯ входных данных
│   └── routes/                  ← РОУТЫ = какой URL к какому хендлеру
│       └── routes.go            ← НАСТРОЙКА всех путей
├── migrations/                  ← SQL ФАЙЛЫ для создания таблиц
│   ├── 000001_create_users.up.sql
│   └── 000001_create_users.down.sql
├── pkg/                         ← ПЕРЕИСПОЛЬЗУЕМЫЙ КОД
│   ├── validator/
│   │   └── validator.go         ← КАСТОМНАЯ ВАЛИДАЦИЯ
│   └── errors/
│       └── errors.go            ← КАСТОМНЫЕ ОШИБКИ
├── .env                         ← СЕКРЕТЫ И НАСТРОЙКИ
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```



### DTO (Data Transfer Objects) - что приходит/уходит

**internal/dto/user\_dto.go**

```go
package dto

// CreateUserRequest - что ПРИХОДИТ когда создают пользователя
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// UserResponse - что ОТДАЕМ обратно (без пароля!)
type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// LoginRequest - что нужно для входа
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```


**internal/dto/response\_dto.go**

```go
package dto

// APIResponse - стандартный ответ API
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// PaginatedResponse - для списков с пагинацией
type PaginatedResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Total   int         `json:"total"`
    Page    int         `json:"page"`
    Limit   int         `json:"limit"`
}
```



### Models - структуры для БД

**internal/models/user.go**

```go
package models

import "time"

// User - точная копия таблицы users в БД
type User struct {
    ID        int       `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Email     string    `db:"email" json:"email"`
    Password  string    `db:"password" json:"-"` // json:"-" = не показывать в API
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
```



### Repository Interface - что УМЕЕТ делать репозиторий

**internal/repository/interfaces/user\_repo.go**

```go
package interfaces

import (
    "context"
    "devprep/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id int) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, offset, limit int) ([]*models.User, error)
}
```



### Repository - КАК работать с БД

**internal/repository/user\_repository.go**

```go
package repository

import (
    "context"
    "database/sql"
    "time"
  
    "devprep/internal/models"
    "devprep/internal/repository/interfaces"
)

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
    query := `
        INSERT INTO users (name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at`
  
    now := time.Now()
    err := r.db.QueryRowContext(ctx, query, 
        user.Name, user.Email, user.Password, now, now,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
  
    if err != nil {
        return nil, err
    }
  
    return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, name, email, password, created_at, updated_at 
              FROM users WHERE email = $1`
  
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID, &user.Name, &user.Email, &user.Password, 
        &user.CreatedAt, &user.UpdatedAt,
    )
  
    if err != nil {
        return nil, err
    }
  
    return user, nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
    query := `SELECT id, name, email, created_at, updated_at 
              FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
  
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
  
    var users []*models.User
    for rows.Next() {
        user := &models.User{}
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
  
    return users, nil
}

```



### Service Interface - что УМЕЕТ делать сервис

**internal/service/interfaces/user\_service.go**

```go
package interfaces

import (
    "context"
    "devprep/internal/dto"
    "devprep/internal/models"
)

type UserService interface {
    CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
    LoginUser(ctx context.Context, req *dto.LoginRequest) (string, error) // возвращает JWT токен
    GetUserByID(ctx context.Context, id int) (*dto.UserResponse, error)
    ListUsers(ctx context.Context, page, limit int) ([]*dto.UserResponse, error)
}
```



### Service - БИЗНЕС-ЛОГИКА

**internal/service/user\_service.go**

```go
package service

import (
    "context"
    "errors"
    "time"
  
    "devprep/internal/dto"
    "devprep/internal/models"
    "devprep/internal/repository/interfaces"
    serviceInterfaces "devprep/internal/service/interfaces"
  
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

type userService struct {
    userRepo interfaces.UserRepository
    jwtSecret string
}

func NewUserService(userRepo interfaces.UserRepository, jwtSecret string) serviceInterfaces.UserService {
    return &userService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. Проверяем, что пользователь не существует
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
  
    // 2. Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
  
    // 3. Создаем модель для БД
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
    }
  
    // 4. Сохраняем в БД
    createdUser, err := s.userRepo.Create(ctx, user)
    if err != nil {
        return nil, err
    }
  
    // 5. Возвращаем DTO (без пароля)
    return &dto.UserResponse{
        ID:    createdUser.ID,
        Name:  createdUser.Name,
        Email: createdUser.Email,
    }, nil
}

func (s *userService) LoginUser(ctx context.Context, req *dto.LoginRequest) (string, error) {
    // 1. Находим пользователя
    user, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
  
    // 2. Проверяем пароль
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return "", errors.New("invalid credentials")
    }
  
    // 3. Создаем JWT токен
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
  
    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }
  
    return tokenString, nil
}
```



### Handler - обработка HTTP

**internal/handlers/user\_handler.go**

```go
package handlers

import (
    "strconv"
  
    "devprep/internal/dto"
    "devprep/internal/service/interfaces"
  
    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
)

type UserHandler struct {
    userService interfaces.UserService
    validator   *validator.Validate
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
        validator:   validator.New(),
    }
}

// POST /api/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req dto.CreateUserRequest
  
    // Парсим JSON
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(dto.APIResponse{
            Success: false,
            Error:   "Invalid JSON format",
        })
    }
  
    // Валидируем
    if err := h.validator.Struct(&req); err != nil {
        return c.Status(400).JSON(dto.APIResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
  
    // Вызываем сервис
    user, err := h.userService.CreateUser(c.Context(), &req)
    if err != nil {
        return c.Status(400).JSON(dto.APIResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
  
    return c.Status(201).JSON(dto.APIResponse{
        Success: true,
        Message: "User created successfully",
        Data:    user,
    })
}

// POST /api/auth/login
func (h *UserHandler) Login(c *fiber.Ctx) error {
    var req dto.LoginRequest
  
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(dto.APIResponse{
            Success: false,
            Error:   "Invalid JSON format",
        })
    }
  
    if err := h.validator.Struct(&req); err != nil {
        return c.Status(400).JSON(dto.APIResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
  
    token, err := h.userService.LoginUser(c.Context(), &req)
    if err != nil {
        return c.Status(401).JSON(dto.APIResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
  
    return c.JSON(dto.APIResponse{
        Success: true,
        Message: "Login successful",
        Data:    fiber.Map{"token": token},
    })
}

// GET /api/users?page=1&limit=10
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
  
    users, err := h.userService.ListUsers(c.Context(), page, limit)
    if err != nil {
        return c.Status(500).JSON(dto.APIResponse{
            Success: false,
            Error:   err.Error(),
        })
    }
  
    return c.JSON(dto.APIResponse{
        Success: true,
        Data:    users,
    })
}
```



### Routes - настройка путей

**internal/routes/routes.go**

```go
package routes

import (
    "devprep/internal/handlers"
    "devprep/internal/middleware"
  
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, jwtSecret string) {
    // Глобальные middleware
    app.Use(logger.New())
    app.Use(cors.New())
  
    // API группа
    api := app.Group("/api")
  
    // Публичные роуты (без токена)
    auth := api.Group("/auth")
    auth.Post("/login", userHandler.Login)
  
    users := api.Group("/users")
    users.Post("/", userHandler.CreateUser) // регистрация
  
    // Защищенные роуты (нужен токен)
    protected := api.Group("/", middleware.AuthMiddleware(jwtSecret))
    protected.Get("/users", userHandler.ListUsers)
  
    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })
}
```



### Middleware - проверки

**internal/middleware/auth.go**

```go
package middleware

import (
    "strings"
  
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Получаем токен из заголовка
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Authorization header required",
            })
        }
  
        // Проверяем формат "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid authorization header format",
            })
        }
  
        tokenString := parts[1]
  
        // Парсим токен
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
  
        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
  
        // Извлекаем данные пользователя
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Locals("user_id", int(claims["user_id"].(float64)))
            c.Locals("email", claims["email"].(string))
        }
  
        return c.Next()
    }
}

```



## Зачем такая структура?

1. **DTO** - четко определяет что приходит и уходит из API
2. **Models** - точное отражение БД таблиц
3. **Repository** - вся работа с БД в одном месте
4. **Service** - вся бизнес-логика в одном месте
5. **Handler** - только обработка HTTP, никакой логики
6. **Middleware** - переиспользуемые проверки
7. **Interfaces** - для тестов и замены реализации

**Поток данных:**

```
HTTP запрос → Handler → Service → Repository → БД
БД → Repository → Service → Handler → HTTP ответ
```



**Плюсы:**

* Легко тестировать каждый слой отдельно
* Можно заменить БД, не трогая бизнес-логику
* Код четко разделен по ответственности
* Легко добавлять новые фичи
