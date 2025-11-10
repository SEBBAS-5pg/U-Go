package service

import (
	"context"
	"errors"
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
)

// UserService maneja la logica de negocio para los perfiles de usuario
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService es la "fabrica"
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserProfile obtiene el perfil basado en el ID del token
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("User ID not provided")
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("Error obtaining profile (service): %v", err)
		return nil, errors.New("Usuario no encontrado")
	}

	return user, nil
}
