package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	// Importa tus modelos
	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
)

// UserRepository maneja la comunicación con la tabla 'users'
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository ("es la fabrica")
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
	// se da un tiempo a la consulta para que se cuelgue
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

// Funcion para T-08 (login)

// GetByEmail busca un usuario por email
// Devuelve sql.ErrNoRows si no se encuentra

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT
	 id, email, full_name, password_hash, is_active, activation_token,
	 is_driver, driver_status, profile_image_url, average_rating, created_at
    	FROM users
		WHERE email = $1
	`

	// Damos un timeout al contexto
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Prepara el struct para escanear
	user := &models.User{}

	// QueryRowCOntext ejecuta la consulta
	err := r.db.QueryRowContext(ctxTimeout, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.IsActive,
		&user.ActivationToken,
		&user.IsDriver,
		&user.DriverStatus,
		&user.ProfileImageURL,
		&user.AverageRating,
		&user.CreatedAt,
	)

	if err != nil {
		// sql.ErrNoRows es el error esperado si el email no existe
		if err == sql.ErrNoRows {
			return nil, err
		}
		// otro error mas
		log.Printf("Error al escanear usuario por email: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
	SELECT
	id, email, full_name, password_hash, is_active, activation_token,
			is_driver, driver_status, profile_image_url, average_rating, created_at
		FROM users
		WHERE id = $1
	`
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &models.User{}

	err := r.db.QueryRowContext(ctxTimeout, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.IsActive,
		&user.ActivationToken,
		&user.IsDriver,
		&user.DriverStatus,
		&user.ProfileImageURL,
		&user.AverageRating,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error scanning user by ID: %v", err)
		return nil, err
	}
	// se limpian los datos sostenibles antes de devolverlos
	user.PasswordHash = ""
	user.ActivationToken = nil

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	query := `
		UPDATE users
		SET 
			full_name = $1,
			updated_at = now()
		WHERE 
			id = $2
		RETURNING 
			id, email, full_name, password_hash, is_active, activation_token,
			is_driver, driver_status, profile_image_url, average_rating, created_at
	`
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &models.User{}

	err := r.db.QueryRowContext(ctxTimeout, query, req.FullName, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
		&user.IsActive,
		&user.ActivationToken,
		&user.IsDriver,
		&user.DriverStatus,
		&user.ProfileImageURL,
		&user.AverageRating,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found to update")
		}
		log.Printf("Error scanning updated user (repo): %v", err)
		return nil, err
	}

	// Limpia datos sensibles antes de devolver
	user.PasswordHash = ""
	user.ActivationToken = nil

	return user, nil
}

// UpdateProfileImageURL actualiza solo la columna profile_image_url
func (r *UserRepository) UpdateProfileImageURL(ctx context.Context, userID string, imageURL string) error {
	query := `
		UPDATE users
		SET 
			profile_image_url = $1
		WHERE 
			id = $2
	`
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctxTimeout, query, imageURL, userID)
	if err != nil {
		log.Printf("Error al actualizar profile_image_url (repo): %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // El usuario no fue encontrado
	}

	return nil
}
