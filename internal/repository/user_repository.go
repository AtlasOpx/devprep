package repository

import (
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/repository/interfaces"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.Insert("users").
		Columns("id", "email", "username", "first_name", "last_name", "password_hash", "role", "is_active", "created_at", "updated_at").
		Values(user.ID, user.Email, user.Username, user.FirstName, user.LastName, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt).
		Exec()
	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Select("id", "email", "username", "first_name", "last_name", "password_hash", "role", "is_active", "created_at", "updated_at").
		From("users").
		Where("id = ?", id).
		QueryRow().
		Scan(&user.ID, &user.Email, &user.Username, &user.FirstName,
			&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Select("id", "email", "username", "first_name", "last_name", "password_hash", "role", "is_active", "created_at", "updated_at").
		From("users").
		Where("email = ?", email).
		QueryRow().
		Scan(&user.ID, &user.Email, &user.Username, &user.FirstName,
			&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Select("id", "email", "username", "first_name", "last_name", "password_hash", "role", "is_active", "created_at", "updated_at").
		From("users").
		Where("username = ?", username).
		QueryRow().
		Scan(&user.ID, &user.Email, &user.Username, &user.FirstName,
			&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(id uuid.UUID, req *models.UpdateProfileRequest) error {
	update := r.db.Update("users").Set("updated_at", "NOW()").Where("id = ?", id)

	if req.FirstName != "" {
		update = update.Set("first_name", req.FirstName)
	}

	if req.LastName != "" {
		update = update.Set("last_name", req.LastName)
	}

	if req.Username != "" {
		update = update.Set("username", req.Username)
	}

	_, err := update.Exec()
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Delete("users").
		Where("id = ?", id).
		Exec()
	return err
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Select("id", "email", "username", "first_name", "last_name", "password_hash", "role", "is_active", "created_at", "updated_at").
		From("users").
		OrderBy("created_at DESC").
		Query()
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
