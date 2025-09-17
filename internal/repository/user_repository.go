package repository

import (
	"fmt"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/repository/interfaces"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
        INSERT INTO users (id, email, username, first_name, last_name, password_hash, role, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.DB.Exec(query, user.ID, user.Email, user.Username, user.FirstName,
		user.LastName, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, email, username, first_name, last_name, password_hash, role, is_active, created_at, updated_at
        FROM users WHERE id = $1`

	err := r.db.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, email, username, first_name, last_name, password_hash, role, is_active, created_at, updated_at
        FROM users WHERE email = $1`

	err := r.db.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, email, username, first_name, last_name, password_hash, role, is_active, created_at, updated_at
        FROM users WHERE username = $1`

	err := r.db.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(id uuid.UUID, req *models.UpdateProfileRequest) error {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.FirstName != "" {
		setParts = append(setParts, "first_name = $"+strconv.Itoa(argIndex))
		args = append(args, req.FirstName)
		argIndex++
	}

	if req.LastName != "" {
		setParts = append(setParts, "last_name = $"+strconv.Itoa(argIndex))
		args = append(args, req.LastName)
		argIndex++
	}

	if req.Username != "" {
		setParts = append(setParts, "username = $"+strconv.Itoa(argIndex))
		args = append(args, req.Username)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.DB.Exec(query, args...)
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.DB.Exec(query, id)
	return err
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
        SELECT id, email, username, first_name, last_name, password_hash, role, is_active, created_at, updated_at
        FROM users ORDER BY created_at DESC`

	rows, err := r.db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.FirstName,
			&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
