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
	query := `
        INSERT INTO sessions (user_id, session_token, expires_at, user_agent, ip_address)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.DB.Exec(query, userID, sessionToken, expiresAt, userAgent, ipAddress)
	return err
}

func (r *AuthRepository) DeleteSession(sessionToken string) error {
	query := "DELETE FROM sessions WHERE session_token = $1"
	_, err := r.db.DB.Exec(query, sessionToken)
	return err
}

func (r *AuthRepository) GetSessionByToken(sessionToken string) (*models.Session, error) {
	var session models.Session
	query := `
        SELECT user_id, session_token, expires_at, user_agent, ip_address, created_at
        FROM sessions WHERE session_token = $1 AND expires_at > NOW()`

	err := r.db.DB.QueryRow(query, sessionToken).Scan(
		&session.UserID, &session.SessionToken, &session.ExpiresAt,
		&session.UserAgent, &session.IPAddress, &session.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) ValidateSession(sessionToken string) (*models.User, error) {
	var user models.User
	query := `
        SELECT u.id, u.email, u.username, u.first_name, u.last_name, u.password_hash, u.role, u.is_active, u.created_at, u.updated_at
        FROM users u
        JOIN sessions s ON u.id = s.user_id
        WHERE s.session_token = $1 AND s.expires_at > NOW() AND u.is_active = true`

	err := r.db.DB.QueryRow(query, sessionToken).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.PasswordHash, &user.Role, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}