package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	// --- Drivers de DB ---
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// --- Tu Stack ---
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Driver de Postgres

	// carpetas internas
	"github.com/SEBBAS-5pg/U-Go/backend/config"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/api"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

func main() {
	// 1. Cargar Configuración (desde config/config.go)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// 2. Conectar a la Base de Datos (Postgres)
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Ping para verificar la conexión
	if err := db.Ping(); err != nil {
		log.Fatal("Error (Ping) connecting to the database:", err)
	}
	log.Println("✅ Successful connection to Postgres (ugo_develop_db).")

	// --- Conectar a MongoDB ---
	ctxMongo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctxMongo, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping a MongoDB
	if err := mongoClient.Ping(ctxMongo, nil); err != nil {
		log.Fatal("Error (Ping) connecting to MongoDB:", err)
	}

	// Asignar la base de datos específica
	mongoDB := mongoClient.Database("ugo_mongo_db")
	log.Println("✅ Successful connection to MongoDB (ugo_mongo_db).")

	// -- Inyeccion de Dependencias (DI) ---

	// (REPOSITORIOS)
	// a. "musculo" (repository) - (necesita la BD)
	userRepo := repository.NewUserRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	locationRepo := repository.NewLocationRepository(mongoDB)
	tripRepo := repository.NewTripRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	//(SERVICIOS)
	// b. "cerebro" (service) - Necesia el repository
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	storageService := service.NewStorageService(mongoDB)
	userService := service.NewUserService(userRepo, locationRepo, storageService, tripRepo, ratingRepo)
	vehicleService := service.NewVehicleService(vehicleRepo, storageService)
	tripService := service.NewTripService(tripRepo, userService, locationRepo)
	ratingService := service.NewRatingService(ratingRepo, tripRepo, userRepo)

	//(HANDLER)
	// c. "mesero" (handler) - necesita el service
	authHandler := api.NewAuthHandler(authService)
	userHandler := api.NewUserHandler(userService)
	vehicleHandler := api.NewVehicleHandler(vehicleService)
	driverHandler := api.NewDriverHandler(userService)
	passengerHandler := api.NewPassengerHandler(userService)
	tripHandler := api.NewTripHandler(tripService)
	ratingHandler := api.NewRatingHandler(ratingService)
	locationHandler := api.NewLocationHandler(locationRepo)

	//(MIDDLEWARE)
	// d. "Guardian de seguridad"
	authMiddleware := api.NewAuthMiddleware(cfg.JWTSecret)

	// configurar el router y rutas (Gorilla Mux)
	router := mux.NewRouter()
	// Definir un sub-router para /api/v1
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Midlewares Directos al sub-router apiV1
	apiV1.Use(jsonContentTypeMiddleware)
	apiV1.Use(enableCORS)

	// -- Rutas Publicas (AUTENTIFICACION) ---
	// --- Endpoints de U-Go ---
	apiV1.HandleFunc("/health", healthHandler).Methods("GET")
	// --- Endpoint de Registro (HU-01) ---
	// Conecta la ruta POST /auth/register con la función authHandler.Register
	apiV1.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	// Ruta Login (T-08)
	apiV1.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	// (Aquí irán los otros endpoints: /vehicles, etc.)
	apiV1.HandleFunc("/trip/{tripId}/location", locationHandler.UpdateTripLocation).Methods("POST")
	apiV1.HandleFunc("/trip/{tripId}/location", locationHandler.GetTripLocation).Methods("GET")

	// -- Rutas Protegidas (REQUIERE TOKEN JWT) ---

	// se crea un sub-router separado para las rutas protegidas
	protectRoutes := apiV1.PathPrefix("").Subrouter()
	protectRoutes.Use(authMiddleware.Middleware) // se aplica el guardia

	// GET /api/v1/users/me
	protectRoutes.HandleFunc("/users/me", userHandler.GetMyProfile).Methods("GET")
	// PUT /users/me
	protectRoutes.HandleFunc("/users/me", userHandler.UpdateMyProfile).Methods("PUT")
	// POST Imagen
	protectRoutes.HandleFunc("/users/me/image", userHandler.UploadProfileImage).Methods("POST")
	// POST de Vehiculos
	protectRoutes.HandleFunc("/vehicles", vehicleHandler.RegisterVehicle).Methods("POST")
	protectRoutes.HandleFunc("/users/me/image", userHandler.GetProfileImage).Methods("GET")
	protectRoutes.HandleFunc("/users/me/image", userHandler.DeleteProfileImage).Methods("DELETE")
	protectRoutes.HandleFunc("/vehicles/{vehicleId}/image", vehicleHandler.GetVehicleImage).Methods("GET")
	protectRoutes.HandleFunc("/vehicles/{vehicleId}/image", vehicleHandler.DeleteVehicleImage).Methods("DELETE")
	// vehicles
	// GET /api/v1/vehicles
	protectRoutes.HandleFunc("/vehicles", vehicleHandler.GetVehicles).Methods("GET")
	protectRoutes.HandleFunc("/driver/location", driverHandler.UpdateLocation).Methods("POST")
	// GET /
	protectRoutes.HandleFunc("/passenger/drivers", passengerHandler.FindNearbyDrivers).Methods("GET").Queries("lat", "{lat}", "lng", "{lng}")
	// NUEVA RUTA PARA EL HISTORIAL DEL PASAJERO (T-14)
	protectRoutes.HandleFunc("/passenger/history", tripHandler.GetPassengerHistory).Methods("GET")

	// NUEVA RUTA PARA EL HISTORIAL (HU-09)
	protectRoutes.HandleFunc("/driver/history", userHandler.GetDriverHistory).Methods("GET")
	// POST /api/v1/vehicles/{vehicleId}/image
	protectRoutes.HandleFunc("/vehicles/{vehicleId}/image", vehicleHandler.UploadVehicleImage).Methods("POST")

	// --- Endpoints de Pasajero
	// POST /api/v1/trips (Solicitar un nuevo viaje)
	protectRoutes.HandleFunc("/trips", tripHandler.CreateTrip).Methods("POST")
	// GET /api/v1/trips/{tripId}
	protectRoutes.HandleFunc("/trips/{tripId}", tripHandler.GetTripByID).Methods("GET")
	// POST /api/v1/trips/{tripId}/accept
	protectRoutes.HandleFunc("/trips/{tripId}/accept", tripHandler.AcceptTrip).Methods("POST")
	// POST /api/v1/trips/{tripId}/cancel
	protectRoutes.HandleFunc("/trips/{tripId}/cancel", tripHandler.CancelTrip).Methods("POST")
	// POST /api/v1/driver/trips/{tripId}/finalize
	protectRoutes.HandleFunc("/driver/trips/{tripId}/finalize", tripHandler.FinalizeTrip).Methods("POST")

	// --- Endpoints de Calificaciones ---
	// POST /api/v1/ratings (Calificar un viaje finalizado)
	protectRoutes.HandleFunc("/ratings", ratingHandler.CreateRating).Methods("POST")
	// (iran las otras rutas protegidas)
	// ...

	// Middlewares (cors + json)
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// iniciar servidor
	log.Printf("✅ Servidor U-Go (Go) corriendo en http://localhost:%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, enhancedRouter))

}

// =============================
// ⚙️ Middlewares (Moveremos esto luego)
// =============================

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// ¡IMPORTANTE! Añadir "Authorization" para que acepte el JWT en el futuro
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// =============================
// 📂 Handlers (Moveremos esto luego)
// =============================
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "OK",
		"project": "U-Go Backend",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
