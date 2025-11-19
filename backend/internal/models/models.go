package models

import (
	"time"
)

// User (usuarios) coincide con la tabla 'users' y tu MER
type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	FullName        string    `json:"full_name"`
	PasswordHash    string    `json:"-"` // Ocultar en JSON
	IsActive        bool      `json:"is_active"`
	ActivationToken *string   `json:"-"` // Ocultar en JSON
	IsDriver        bool      `json:"is_driver"`
	DriverStatus    string    `json:"driver_status"`     // 'offline', 'online', 'en_viaje'
	ProfileImageURL *string   `json:"profile_image_url"` // Puntero para NULOS
	AverageRating   float64   `json:"average_rating"`
	CreatedAt       time.Time `json:"created_at"`
}

// Vehicle (Vehículo) coincide con la tabla 'vehicles'
type Vehicle struct {
	ID              string    `json:"id"`
	ConductorID     string    `json:"conductorId"`
	Plate           string    `json:"plate"`
	Model           string    `json:"model"`
	Color           string    `json:"color"`
	Status          string    `json:"status"` // 'pendiente', 'aprobado', 'rechazado'
	VehicleImageURL *string   `json:"vehicleImageUrl"`
	CreatedAt       time.Time `json:"createdAt"`
}

// Trip (Viaje) coincide con la tabla 'trips'
type Trip struct {
	ID             string     `json:"id"`
	PasajeroID     string     `json:"pasajeroId"`
	ConductorID    *string    `json:"conductorId"` // Puntero para NULOS
	VehicleID      *string    `json:"vehicleId"`   // Puntero para NULOS
	Status         string     `json:"status"`
	OriginLat      float64    `json:"originLat"`
	OriginLng      float64    `json:"originLng"`
	DestinationLat float64    `json:"destinationLat"`
	DestinationLng float64    `json:"destinationLng"`
	CreatedAt      time.Time  `json:"createdAt"`
	StartedAt      *time.Time `json:"startedAt"`
	CompletedAt    *time.Time `json:"completedAt"`
}

// Rating (Calificacion) coincide con la tabla 'ratings'
type Rating struct {
	ID        string    `json:"id"`
	TripID    string    `json:"tripId"`
	RaterID   string    `json:"raterId"`
	RatedID   string    `json:"ratedId"`
	Rating    int       `json:"rating"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
}

type ContextKey string

const ContextUserIDKey ContextKey = "userID"

// --- Modelos Geoespaciales y de Geolocalización (HU-08) ---

// GeoJson representa el formato GeoJSON Point [longitud, latitud]
type GeoJson struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

// DriverLocation representa el documento que se guardará en MongoDB
type DriverLocation struct {
	UserID    string    `json:"user_id" bson:"userid"`
	Status    string    `json:"status" bson:"status"`
	Location  GeoJson   `json:"location" bson:"location"`
	UpdatedAt time.Time `json:"updated_at" bson:"updatedat"`
}

// UpdateLocationRequest define la estructura esperada del JSON de entrada para la ubicación
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"` // ej: "online", "offline"
}
