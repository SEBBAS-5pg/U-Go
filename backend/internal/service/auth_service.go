package service

import (
	"context"
	"crypto/rand"  //genera tokens seguros
	"encoding/hex" //Convierte el token a string
	"errors"       //crea errores propios
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// AuthService es el "cerebro" que maneja la logica de negocio de autentificacion
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService es la. "fabrica" para nuestro servicio
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register es la funcion principal de logica de negocio para HU-01
func (s *AuthService) Register(ctx context.Context, email, fullName, password string) (*models.User, error) {
	// logica de validacion
	if len(password) < 8 {
		return nil, errors.New("The password must be at least 8 characters long")
	}
	// (aqui ira la validacion del email institucional con Regex)

	// logica de seguridad (bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing the password: %v", err)
		return nil, errors.New("Internal error processing password")
	}

	// Logica de Negocio (token de activacion)
	token, err := generateSecureToken(32) // 32 bytes = 64 caracteres
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, errors.New("Internal error generating token")
	}

	// Prepara el Modelo (datos por defecto)
	user := &models.User{
		Email:           email,
		FullName:        fullName,
		PasswordHash:    string(hashedPassword),
		ActivationToken: &token, // es un puntero
		IsActive:        false,
		IsDriver:        false,
		DriverStatus:    "offline",
	}
	// Lamma la capa de Repositori (Guarda en BD)
	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("Error saving user to database: %v", err)
		// codigo de error 23505 en Postgres
		return nil, errors.New("The email is already in use")
	}

	// Asignamos el ID que la BD nos devolvio
	user.ID = userID
	user.PasswordHash = "" // Se limpia la contraseña

	// Logica de Notificacion (Email)
	log.Printf("SIMULACIÓN DE EMAIL: Enviando token de activación a %s. Token: %s", user.Email, token)

	return user, nil
}

// --Funciones de Ayuda (Helpers)

// generateSecureToken crea un string aleatorio seguro para el token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
