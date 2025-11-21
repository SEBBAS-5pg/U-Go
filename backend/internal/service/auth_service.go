package service

import (
	"context"
	"crypto/rand"  //genera tokens seguros
	"database/sql" //Importa para sql.ErrNoRows
	"encoding/hex" //Convierte el token a string
	"errors"       //crea errores propios
	"log"
	"time" // Importa time para el JWT

	// Importa el JWT
	"github.com/golang-jwt/jwt/v5"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// Define un tipo personalizado para nuestra llave de contexto
type contextKey string

// Esta es la llave que usa para guardar/leer el userID del contexto
const ContextKeyUserID contextKey = "userID"

// AuthService es el "cerebro" que maneja la logica de negocio de autentificacion
type AuthService struct {
	userRepo *repository.UserRepository

	// secret para firmar los tokens
	jwtSecret string
}

// NewAuthService es la. "fabrica" para nuestro servicio
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		// JWT
		jwtSecret: jwtSecret,
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
	// Lama la capa de Repositori (Guarda en BD)
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

// Funcion de Login (T-08)

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Busca el usuario por email (usando la funcion del repositorio)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Invalid credentials")
		}
		log.Printf("Error searching for user by email: %v", err)
		return "", errors.New("Internal Server Error")
	}

	// Comparar la contrasela (la de la BD con la del envio del usuario)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// si err != nnil, las contraseñas no coinciden
		return "", errors.New("Invalid credentials")
	}
	/*
		// verifica si la cuenta esta activa
		if !user.IsActive {
			return "", errors.New("Please activate your account (check your email)")
		}
	*/

	// contraseña incorrecta Generar el JWT
	tokenString, err := s.generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return "", errors.New("Internal error generating token")
	}
	return tokenString, nil
}

// funcion helper para JWT

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	// crea los "Claims" (informacion dentro del token)
	claims := jwt.MapClaims{
		"sub": user.ID, // "subject" el ID del usuario
		"nam": user.FullName,
		"rol": user.IsDriver,                             // Rol (true/false)
		"iat": time.Now().Unix(),                         // "Issud At" Cuando se creó
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Expira en 7 dias
	}

	// Crea un nuevo token con los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con nuestro JWT_SECRET del .env
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
