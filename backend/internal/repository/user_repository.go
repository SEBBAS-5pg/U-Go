package repository

import (
	"context"
	"database/sql"
	"time"

	// Importa tus modelos
	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
)

// UserRepository maneja la comunicación con la tabla 'users'
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserta un nuevo usuario en la base de datos
func (r *UserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	var userID string

	//Este es el sql que usa driver(lob/pq)
	query := `
		INSERT INTO users 
			(email, full_name, password_hash, activation_token, is_active, is_driver, driver_status)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Ejecuta la consulta(query)
	err := r.db.QueryRowContext(ctxTimeout, query,
		user.Email,
		user.FullName,
		user.PasswordHash,
		user.ActivationToken,
		user.IsActive,
		user.IsDriver,
		user.DriverStatus,
	).Scan(&userID) //.scan() lee el returnun id

	if err != nil {
		return "", err
	}

	return userID, nil
}
