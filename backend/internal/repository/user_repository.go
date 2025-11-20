package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

	// 1. Construir la consulta de forma dinámica
	updates := []string{"updated_at = NOW()"} // Campo de actualización automática
	args := []interface{}{}
	paramCounter := 1 // Contador para $1, $2, etc.

	// Añadir FullName si se proporciona
	if req.FullName != "" {
		updates = append(updates, fmt.Sprintf("full_name = $%d", paramCounter))
		args = append(args, req.FullName)
		paramCounter++
	}

	// Añadir IsDriver si se proporciona (req.IsDriver no es nil)
	if req.IsDriver != nil {
		updates = append(updates, fmt.Sprintf("is_driver = $%d", paramCounter))
		args = append(args, *req.IsDriver) // Desreferenciar el puntero *bool
		paramCounter++
	}

	// Si no hay campos para actualizar además de updated_at
	if len(updates) <= 1 {
		// Podríamos devolver nil, nil para indicar que no hubo cambios o forzar la actualización de full_name
		return nil, errors.New("debe proporcionar al menos 'full_name' o 'is_driver' para actualizar")
	}

	setClause := strings.Join(updates, ", ")

	// El último parámetro siempre será el userID para la cláusula WHERE
	args = append(args, userID)

	// El query final usa el contador de parámetros actual para el WHERE id = $N
	query := fmt.Sprintf(`
        UPDATE users
        SET %s
        WHERE id = $%d
        RETURNING 
            id, email, full_name, password_hash, is_active, activation_token,
            is_driver, driver_status, profile_image_url, average_rating, created_at
    `, setClause, paramCounter)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &models.User{}

	// 2. Ejecutar la consulta con la lista dinámica de argumentos
	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(
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

// UpdateDriverStatus actualiza el campo driver_status de un usuario en PostgreSQL
func (r *UserRepository) UpdateDriverStatus(ctx context.Context, userID string, status string) error {
	query := `
        UPDATE users
        SET 
            driver_status = $1::driver_status, -- <<-- ¡CORRECCIÓN CLAVE! Hacemos CAST explícito al tipo ENUM
            is_active = ($1 = 'online') 
        WHERE id = $2
    `

	result, err := r.db.ExecContext(ctx, query, status, userID)
	if err != nil {
		log.Printf("Error al actualizar driver_status en DB: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Retornar error si el usuario no fue encontrado
	}

	return nil
}

// UpdateAverageRating recalcula y actualiza el promedio de calificación de un usuario (conductor)
func (r *UserRepository) UpdateAverageRating(ctx context.Context, userID string) error {
	// 1. Consulta SQL: Calcula el promedio de la columna 'rating' de la tabla 'ratings'
	// donde rated_id es el conductor, y actualiza la columna 'average_rating' en 'users'.
	query := `
		UPDATE users 
		SET average_rating = (
			SELECT COALESCE(AVG(rating), 0) 
			FROM ratings 
			WHERE rated_id = $1
		)
		WHERE id = $1
	`
	// NOTA: COALESCE(AVG(rating), 0) asegura que si no hay ratings, el promedio sea 0.

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error al actualizar el promedio de rating para el usuario %s: %v", userID, err)
		return errors.New("error al actualizar el promedio de calificación")
	}

	return nil
}
