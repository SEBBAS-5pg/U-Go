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

// CONSTANTES PARA EL ESTADO DEL CONDUCTOR (DriverStatus)
const (
	DriverStatusOnline  = "online"
	DriverStatusOffline = "offline"
	DriverStatusEnViaje = "en_viaje"
)

// Tipos ENUM para el estado del viaje (Coinciden con PostgreSQL)
const (
	TripStatusSolicitado = "solicitado"
	TripStatusAceptado   = "aceptado"
	TripStatusEnCurso    = "en_curso"
	TripStatusFinalizado = "finalizado"
	TripStatusCancelado  = "cancelado"
)

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
	ID              string     `json:"id"`
	PasajeroID      string     `json:"pasajeroId"`
	ConductorID     string     `json:"conductorId"` // Puede ser nulo hasta que se acepte
	VehicleID       string     `json:"vehicleId"`   // ID del vehículo usado
	Status          string     `json:"status"`      // Usará los ENUMs de arriba
	OriginLat       float64    `json:"originLat"`
	OriginLng       float64    `json:"originLng"`
	OriginName      string     `json:"originName,omitempty"`
	DestinationLat  float64    `json:"destinationLat"`
	DestinationLng  float64    `json:"destinationLng"`
	DestinationName string     `json:"destinationName,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	CompletedAt     *time.Time `json:"completedAt,omitempty"`
}

// CreateTripRequest (Petición) es la estructura que envía el pasajero
type CreateTripRequest struct {
	ConductorID     string  `json:"conductor_id"`
	OriginLat       float64 `json:"origin_lat"`
	OriginLng       float64 `json:"origin_lng"`
	OriginName      string  `json:"origin_name"`
	DestinationLat  float64 `json:"destination_lat"`
	DestinationLng  float64 `json:"destination_lng"`
	DestinationName string  `json:"destination_name"`
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
	IsDriver *bool  `json:"is_driver,omitempty"`
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
	TripID    string    `json:"trip_id,omitempty" bson:"trip_id,omitempty"`
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

// AcceptTripRequest es la estructura que usa el conductor para aceptar el viaje
type AcceptTripRequest struct {
	VehicleID string `json:"vehicle_id"`
}

// FinalizeTripRequest es la estructura que usa el conductor para finalizar el viaje
type FinalizeTripRequest struct {
	FinalLat float64 `json:"final_lat"` // Opcional: Para registrar dónde terminó
	FinalLng float64 `json:"final_lng"` // Opcional: Para registrar dónde terminó
}

// Añadir la estructura de solicitud de calificación
type CreateRatingRequest struct {
	TripID  string `json:"trip_id" validate:"required"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"` // Usamos Rating, consistente con el modelo y el JSON
	Comment string `json:"comment"`
}

// DriverHistoryResponse estructura el historial del conductor
type DriverHistoryResponse struct {
	UserID        string  `json:"user_id"`
	FullName      string  `json:"full_name"`
	AverageRating float64 `json:"average_rating"` // Del user

	// Lista de viajes
	Trips []Trip `json:"trips"` // Asumo que Trip ya está definido

	// Lista de calificaciones que ha recibido
	ReceivedRatings []Rating `json:"received_ratings"` // Asumo que Rating ya está definido
}
