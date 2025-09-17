package repository

import (
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/repository/interfaces"
	"github.com/google/uuid"
	"time"
)

type AuthRepository struct {
	db *database.DB
}

func NewAuthRepository(db *database.DB) interfaces.AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateSession(userID uuid.UUID, sessionToken string, expiresAt time.Time, userAgent, ipAddress string) error {
	_, err := r.db.Insert("sessions").
		Columns("user_id", "session_token", "expires_at", "user_agent", "ip_address").
		Values(userID, sessionToken, expiresAt, userAgent, ipAddress).
		Exec()
	return err
}

func (r *AuthRepository) DeleteSession(sessionToken string) error {
	_, err := r.db.Delete("sessions").
		Where("session_token = ?", sessionToken).
		Exec()
	return err
}

func (r *AuthRepository) GetSessionByToken(sessionToken string) (*models.Session, error) {
	var session models.Session
	err := r.db.Select("user_id", "session_token", "expires_at", "user_agent", "ip_address", "created_at").
		From("sessions").
		Where("session_token = ? AND expires_at > NOW()", sessionToken).
		QueryRow().
		Scan(&session.UserID, &session.SessionToken, &session.ExpiresAt,
			&session.UserAgent, &session.IPAddress, &session.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) ValidateSession(sessionToken string) (*models.User, error) {
	var user models.User
	err := r.db.Select("u.id", "u.email", "u.username", "u.first_name", "u.last_name", "u.password_hash", "u.role", "u.is_active", "u.created_at", "u.updated_at").
		From("users u").
		Join("sessions s ON u.id = s.user_id").
		Where("s.session_token = ? AND s.expires_at > NOW() AND u.is_active = true", sessionToken).
		QueryRow().
		Scan(&user.ID, &user.Email, &user.Username, &user.FirstName,
			&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
