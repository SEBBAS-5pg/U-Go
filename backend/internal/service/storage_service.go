// Ruta: /backend/internal/service/storage_service.go

package service

import (
	"context"
	"fmt"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StorageService maneja la subida de archivos a MongoDB GridFS
type StorageService struct {
	db *mongo.Database
}

// NewStorageService es la fábrica
func NewStorageService(db *mongo.Database) *StorageService {
	return &StorageService{db: db}
}

// UploadFile guarda un archivo en GridFS y devuelve su ID
// Usamos el 'filename' (ej. 'profile_user_UUID.jpg') para guardarlo
func (s *StorageService) UploadFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	// (Usamos el nombre de la DB por defecto para GridFS)
	bucket, err := gridfs.NewBucket(
		s.db,
		options.GridFSBucket().SetName("uploads"), // Nombre de la colección
	)
	if err != nil {
		log.Printf("Error al crear bucket GridFS: %v", err)
		return "", err
	}

	// Abrir el stream de subida
	uploadStream, err := bucket.OpenUploadStream(
		filename, // El nombre de archivo único que le damos
	)
	if err != nil {
		log.Printf("Error al abrir stream de subida: %v", err)
		return "", err
	}
	defer uploadStream.Close()

	// Copiar el contenido del 'file' (la petición HTTP) al 'stream' (Mongo)
	if _, err = io.Copy(uploadStream, file); err != nil {
		log.Printf("Error al copiar archivo a GridFS: %v", err)
		return "", err
	}

	// El 'FileID' es un 'ObjectID' de Mongo
	fileID := uploadStream.FileID.(primitive.ObjectID)
	// Devolvemos el ID como un string simple
	return fileID.Hex(), nil
}

// DeleteFile elimina un archivo de GridFS usando su fileID (que es un ObjectID en formato Hex)
func (s *StorageService) DeleteFile(ctx context.Context, fileIDHex string) error {
	bucket, err := gridfs.NewBucket(
		s.db,
		options.GridFSBucket().SetName("uploads"),
	)
	if err != nil {
		return fmt.Errorf("error al crear bucket GridFS para eliminar: %w", err)
	}

	// Convertir el string Hex (que guardamos como URL) de vuelta a primitive.ObjectID
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		return fmt.Errorf("ID de archivo inválido: %w", err)
	}

	// Ejecutar la eliminación
	if err := bucket.Delete(fileID); err != nil {
		if err == mongo.ErrNoDocuments {
			// Si el archivo no se encuentra en GridFS, no lo tratamos como error fatal
			log.Printf("Advertencia: El archivo GridFS %s no se encontró para eliminar.", fileIDHex)
			return nil
		}
		log.Printf("Error al eliminar archivo GridFS %s: %v", fileIDHex, err)
		return fmt.Errorf("error al eliminar archivo del storage: %w", err)
	}

	log.Printf("✅ Archivo GridFS eliminado con ID: %s", fileIDHex)
	return nil
}
