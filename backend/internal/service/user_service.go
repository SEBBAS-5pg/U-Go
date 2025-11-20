package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
)

// UserService maneja la logica de negocio para los perfiles de usuario
// UserService maneja la lógica de negocio para los usuarios
type UserService struct {
	userRepo     *repository.UserRepository
	locationRepo *repository.LocationRepository
	storageRepo  *StorageService
}

// NewUserService es la "fábrica"
func NewUserService(
	userRepo *repository.UserRepository,
	locationRepo *repository.LocationRepository,
	storageRepo *StorageService, // ¡Ahora recibe 3 argumentos!
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		locationRepo: locationRepo,
		storageRepo:  storageRepo, // Asignación correcta
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
	if req.FullName != "" {
		if len(req.FullName) < 3 {
			return nil, errors.New("El nombre completo debe tener al menos 3 caracteres")
		}
	}

	// Si solo se proporciona IsDriver y es nil, esto pasaría.
	if req.FullName == "" && req.IsDriver == nil {
		return nil, errors.New("Debe proporcionar al menos 'full_name' o 'is_driver' para actualizar")
	}

	// Llamar al repositorio para hacer la actualización (ahora maneja full_name e is_driver)
	updatedUser, err := s.userRepo.UpdateUser(ctx, userID, req)
	if err != nil {
		log.Printf("Error al actualizar perfil (service): %v", err)
		return nil, errors.New("No se pudo actualizar el perfil del usuario")
	}
	if req.IsDriver != nil && *req.IsDriver {

		// 1. Inicializar ubicación en MongoDB
		// Asignamos una ubicación inicial (0, 0) y el estado 'offline' por defecto
		locationErr := s.locationRepo.UpsertDriverLocation(ctx, userID,
			0.0,
			0.0,
			models.DriverStatusOffline) // Estado inicial por defecto en Mongo

		if locationErr != nil {
			// Esto no debería ser un error fatal, pero lo registramos.
			// El usuario ya está actualizado en Postgres, solo falló Mongo.
			log.Printf("ADVERTENCIA: Falló la inicialización de MongoDB para el conductor %s: %v", userID, locationErr)
			// Se puede optar por devolver el error fatal si la inicialización de Mongo es crítica.
			// Por simplicidad, aquí permitimos que continúe.
		}
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

// UpdateDriverStatusAndLocation maneja la actualización del estado (Postgres) y la geolocalización (Mongo)
func (s *UserService) UpdateDriverStatusAndLocation(ctx context.Context, userID string, req models.UpdateLocationRequest) error {

	// 1. Validar la solicitud
	if userID == "" {
		return errors.New("ID de conductor no proporcionado")
	}
	if req.Status == "" || (req.Status != "online" && req.Status != "offline") {
		return errors.New("Estado (status) inválido. Debe ser 'online' o 'offline'")
	}

	// 2. Actualizar el estado del conductor en PostgreSQL
	err := s.userRepo.UpdateDriverStatus(ctx, userID, req.Status)
	if err != nil {
		log.Printf("Error al actualizar el estado del conductor en DB (Postgres): %v", err)
		if err == sql.ErrNoRows {
			return errors.New("Conductor no encontrado")
		}
		return errors.New("No se pudo actualizar el estado del conductor")
	}

	// 3. Actualizar la geolocalización en MongoDB
	err = s.locationRepo.UpsertDriverLocation(ctx, userID, req.Latitude, req.Longitude, req.Status)
	if err != nil {
		log.Printf("Error al actualizar la ubicación en MongoDB: %v", err)
		return errors.New("No se pudo registrar la ubicación del conductor")
	}

	return nil
}

// FindNearbyDrivers coordina la busqueda de conductores online en MongoDB
func (s *UserService) FindNearbyDrivers(ctx context.Context, latitude float64, longitude float64) ([]models.DriverLocation, error) {

	// 1. Validación básica de coordenadas (opcional pero buena práctica)
	if math.Abs(latitude) > 90 || math.Abs(longitude) > 180 {
		return nil, errors.New("Coordenadas geográficas inválidas")
	}

	// Definimos la distancia máxima para la búsqueda (ej: 10 km = 10000 metros)
	const maxDistanceMeters = 10000

	// 2. Llamar al repositorio de MongoDB para obtener la lista
	drivers, err := s.locationRepo.FindNearbyDrivers(ctx, latitude, longitude, maxDistanceMeters)
	if err != nil {
		log.Printf("Error al buscar conductores cercanos (service): %v", err)
		return nil, errors.New("Error interno al buscar conductores")
	}

	// 3. Devolver la lista
	return drivers, nil
}

// UpdateDriverStatus solo actualiza el estado del conductor en PostgreSQL (sin tocar Mongo)
func (s *UserService) UpdateDriverStatus(ctx context.Context, userID string, status string) error {
	if userID == "" {
		return errors.New("ID de conductor no proporcionado")
	}

	// 1. Actualizar el estado del conductor en PostgreSQL
	err := s.userRepo.UpdateDriverStatus(ctx, userID, status)
	if err != nil {
		log.Printf("Error al actualizar el estado del conductor en DB (Postgres): %v", err)
		if err == sql.ErrNoRows {
			return errors.New("Conductor no encontrado")
		}
		return errors.New("No se pudo actualizar el estado del conductor")
	}

	return nil
}
