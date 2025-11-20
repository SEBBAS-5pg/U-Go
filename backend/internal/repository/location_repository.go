package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
)

// LocationRepository maneja la interacción con MongoDB para la geolocalización
type LocationRepository struct {
	Collection *mongo.Collection
}

// NewLocationRepository es la fábrica
func NewLocationRepository(db *mongo.Database) *LocationRepository {
	// Definimos la colección que usaremos
	collection := db.Collection("driver_locations")

	// ⚠️ IMPORTANTE: Crear índice Geo-Especial (2dsphere)
	// Esto es crucial para hacer consultas eficientes por cercanía (ej. buscar taxis cercanos)
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"location", "2dsphere"}},
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Error al crear índice 2dsphere en MongoDB: %v", err)
	}

	return &LocationRepository{Collection: collection}
}

// UpsertDriverLocation inserta o actualiza la posición del conductor
func (r *LocationRepository) UpsertDriverLocation(ctx context.Context, userID string, lat, lng float64, status string) error {

	// 1. Definir la estructura GeoJSON Point y el documento a guardar
	locationDoc := models.DriverLocation{
		UserID: userID,
		Status: status,
		Location: models.GeoJson{
			Type:        "Point",
			Coordinates: []float64{lng, lat}, // GeoJSON usa [longitud, latitud]
		},
		UpdatedAt: time.Now(),
	}

	// 2. Definir el filtro (busca por userID) y la operación de actualización
	filter := bson.M{"userid": userID}

	// Usamos $set para actualizar todos los campos si el documento existe
	update := bson.M{"$set": locationDoc}

	// 3. Opciones: Upsert = true (inserta si no existe, actualiza si existe)
	opts := options.Update().SetUpsert(true)

	// 4. Ejecutar la operación
	_, err := r.Collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error al hacer Upsert de la ubicación del conductor %s: %v", userID, err)
		return err
	}

	return nil
}

// FindNearbyDrivers busca conductores online dentro de un radio (en metros) de un punto
func (r *LocationRepository) FindNearbyDrivers(ctx context.Context, latitude float64, longitude float64, maxDistanceMeters int) ([]models.DriverLocation, error) {

	// 1. Definir el punto de búsqueda (Pasajero)
	centerPoint := models.GeoJson{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude}, // MongoDB usa [longitud, latitud]
	}

	// 2. Definir los criterios de búsqueda ($geoNear)
	// El $geoNear requiere que el índice '2dsphere' exista, lo cual ya configuramos en UpsertDriverLocation.
	pipeline := []bson.M{
		{
			"$geoNear": bson.M{
				"near":          centerPoint,                // Punto central de la búsqueda (Pasajero)
				"distanceField": "distance",                 // Nombre del nuevo campo que contendrá la distancia al conductor
				"maxDistance":   maxDistanceMeters,          // Radio máximo en metros
				"spherical":     true,                       // Indica que la distancia es calculada en una esfera (Tierra)
				"query":         bson.M{"status": "online"}, // ¡Solo buscar conductores con status: "online"!
			},
		},
	}

	// 3. Ejecutar la agregación
	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error al ejecutar $geoNear en MongoDB: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// 4. Decodificar los resultados
	var drivers []models.DriverLocation
	if err = cursor.All(ctx, &drivers); err != nil {
		log.Printf("Error al decodificar resultados de $geoNear: %v", err)
		return nil, err
	}

	return drivers, nil
}
