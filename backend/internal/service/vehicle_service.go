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

// VehicleService maneja la lógica de negocio para los vehículos
type VehicleService struct {
	vehicleRepo *repository.VehicleRepository
	storageRepo *StorageService // Necesario para la HU-07 completa (subida de fotos)
}

// NewVehicleService es la "fábrica"
func NewVehicleService(vehicleRepo *repository.VehicleRepository, storageRepo *StorageService) *VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepo,
		storageRepo: storageRepo,
	}
}

// RegisterVehicle registra un nuevo vehículo en el sistema
// La lógica de la subida de imagen se manejará en el Handler,
// pero el servicio es el que coordina ambos repositorios.
func (s *VehicleService) RegisterVehicle(ctx context.Context, conductorID string, plate string, model string, color string) (*models.Vehicle, error) {

	// --- Validación de Lógica de Negocio ---
	if plate == "" || model == "" || color == "" {
		return nil, errors.New("los campos Plate, Model y Color son obligatorios")
	}

	// 1. Crear el objeto Vehicle a registrar
	newVehicle := &models.Vehicle{
		ConductorID: conductorID, // Este ID viene del JWT (AuthMiddleware)
		Plate:       plate,
		Model:       model,
		Color:       color,
		// Status y CreatedAt son gestionados por el repositorio
	}

	// 2. Llamar al repositorio para guardar en PostgreSQL
	registeredVehicle, err := s.vehicleRepo.CreateVehicle(ctx, newVehicle)
	if err != nil {
		// Asumimos que si hay un error de DB (ej. placa duplicada)
		// el error se logueó en el repositorio, aquí devolvemos un error genérico
		log.Printf("Error al registrar vehículo en DB (service): %v", err)
		return nil, errors.New("no se pudo registrar el vehículo")
	}

	return registeredVehicle, nil
}

// GetVehiclesByConductor obtiene la lista de vehículos registrados por un conductor
func (s *VehicleService) GetVehiclesByConductor(ctx context.Context, conductorID string) ([]models.Vehicle, error) {
	if conductorID == "" {
		return nil, errors.New("ID de conductor no proporcionado")
	}

	vehicles, err := s.vehicleRepo.GetVehiclesByConductorID(ctx, conductorID)
	if err != nil {
		log.Printf("Error al obtener vehículos (service): %v", err)
		return nil, errors.New("no se pudo obtener la lista de vehículos")
	}

	// Lógica de negocio adicional aquí, si fuera necesaria (ej. filtrar por estado)

	return vehicles, nil
}

// UpdateVehicleImage maneja la lógica de subir una foto de vehículo
func (s *VehicleService) UpdateVehicleImage(ctx context.Context, vehicleID string, file io.Reader, fileHeaderFilename string) (string, error) {

	// 1. Crear un nombre de archivo único
	// Usamos el ID del vehículo + el nombre original
	uniqueFilename := fmt.Sprintf("vehicle_%s_%s", vehicleID, fileHeaderFilename)

	// 2. Subir el archivo a MongoDB (GridFS)
	fileID, err := s.storageRepo.UploadFile(ctx, file, uniqueFilename)
	if err != nil {
		log.Printf("Error al subir archivo de vehículo a storage (service): %v", err)
		return "", errors.New("no se pudo guardar el archivo de la imagen del vehículo")
	}

	// 3. Crear la URL pública (en este caso, es el ID de Mongo)
	imageURL := fileID

	// 4. Actualizar la URL en PostgreSQL
	err = s.vehicleRepo.UpdateVehicleImageURL(ctx, vehicleID, imageURL)
	if err != nil {
		log.Printf("Error al actualizar URL en Postgres (service): %v", err)
		// En un sistema real, aquí se borraría el archivo de Mongo si falla Postgres.
		return "", errors.New("no se pudo asociar la imagen al vehículo")
	}

	// 5. Devolver la URL/ID
	return imageURL, nil
}

// GetVehicleImageURL obtiene la URL (o ID de Mongo) de la imagen de un vehículo.
func (s *VehicleService) GetVehicleImageURL(ctx context.Context, vehicleID string) (string, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return "", errors.New("vehículo no encontrado")
	}
	if *vehicle.VehicleImageURL == "" {
		return "", errors.New("imagen del vehículo no configurada")
	}
	return *vehicle.VehicleImageURL, nil
}

// DeleteVehicleImage elimina la imagen del vehículo del storage y resetea la URL en PostgreSQL.
func (s *VehicleService) DeleteVehicleImage(ctx context.Context, vehicleID string) error {
	// 1. Obtener la URL/ID actual
	vehicle, err := s.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return errors.New("vehículo no encontrado")
	}

	if *vehicle.VehicleImageURL != "" {
		// 2. Eliminar el archivo del Storage (MongoDB GridFS)
		err = s.storageRepo.DeleteFile(ctx, *vehicle.VehicleImageURL)
		if err != nil {
			log.Printf("Advertencia: Falló la eliminación del archivo %s de MongoDB: %v", *vehicle.VehicleImageURL, err)
		}
	}

	// 3. Actualizar la URL a "" en PostgreSQL
	err = s.vehicleRepo.UpdateVehicleImageURL(ctx, vehicleID, "")
	if err != nil {
		log.Printf("Error al actualizar URL en Postgres (service): %v", err)
		return errors.New("no se pudo eliminar la asociación de la imagen del vehículo")
	}

	return nil
}
