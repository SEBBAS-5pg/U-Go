package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
)

// UserService maneja la logica de negocio para los perfiles de usuario
type UserService struct {
	userRepo    *repository.UserRepository
	storageRepo *StorageService
}

// NewUserService es la "fabrica"
func NewUserService(userRepo *repository.UserRepository, storageRepo *StorageService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		storageRepo: storageRepo,
	}
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

// UpdateUserProfile valida y actualiza el perfil de un usario
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	// --- Validación de Lógica de Negocio ---
	if req.FullName == "" {
		return nil, errors.New("El nombre completo (full_name) no puede estar vacío")
	}
	if len(req.FullName) < 3 {
		return nil, errors.New("El nombre completo debe tener al menos 3 caracteres")
	}

	// Llamar al repositorio para hacer la actualización
	updatedUser, err := s.userRepo.UpdateUser(ctx, userID, req)
	if err != nil {
		log.Printf("Error al actualizar perfil (service): %v", err)
		return nil, errors.New("No se pudo actualizar el perfil del usuario")
	}

	return updatedUser, nil
}

// UpdateUserProfileImage maneja la lógica de subir una foto de perfil
func (s *UserService) UpdateUserProfileImage(ctx context.Context, userID string, file io.Reader, fileHeaderFilename string) (string, error) {

	// 1. Crear un nombre de archivo único
	// (En un sistema real, usaríamos un UUID + extensión, pero esto funciona)
	uniqueFilename := fmt.Sprintf("profile_%s_%s", userID, fileHeaderFilename)

	// 2. Subir el archivo a MongoDB (GridFS)
	fileID, err := s.storageRepo.UploadFile(ctx, file, uniqueFilename)
	if err != nil {
		log.Printf("Error al subir archivo a storage (service): %v", err)
		return "", errors.New("No se pudo guardar el archivo")
	}

	// Crear una URL pública. Por ahora, solo guardaremos el ID de Mongo.
	// En producción, esto sería: "https://cdn.u-go.com/images/" + fileID
	// Por ahora, el ID de Mongo es suficiente.
	imageURL := fileID

	// 4. Actualizar la URL en PostgreSQL
	err = s.userRepo.UpdateProfileImageURL(ctx, userID, imageURL)
	if err != nil {
		log.Printf("Error al actualizar URL en Postgres (service): %v", err)
		// (Aquí deberíamos borrar el archivo de Mongo...
		// pero lo dejaremos así por simplicidad del proyecto)
		return "", errors.New("No se pudo asociar la imagen al perfil")
	}

	// 5. Devolver la URL/ID
	return imageURL, nil
}
