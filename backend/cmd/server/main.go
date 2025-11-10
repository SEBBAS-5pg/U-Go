package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

	// -- Inyeccion de Dependencias (DI)
	// Conectar en el orden correcto:

	// a. "musculo" (repository) - (necesita la BD)
	userRepo := repository.NewUserRepository(db)
	// b. "cerebro" (service) - Necesia el repository
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	// c. "mesero" (handler) - necesita el service
	authHandler := *api.NewAuthHandler(authService)

	// configurar el router y rutas (Gorilla Mux)
	router := mux.NewRouter()
	// Definir un sub-router para /api/v1
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Midlewares Directos al sub-router apiV1
	apiV1.Use(jsonContentTypeMiddleware)
	apiV1.Use(enableCORS)

	// --- Endpoints de U-Go ---
	apiV1.HandleFunc("/health", healthHandler).Methods("GET")
	// --- Endpoint de Registro (HU-01) ---
	// Conecta la ruta POST /auth/register con la función authHandler.Register
	apiV1.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	// Ruta Login (T-08)
	apiV1.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	// (Aquí irán los otros endpoints: /vehicles, etc.)
	// ...

	// Middlewares (cors + json)
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// iniciar servidor
	log.Printf("✅ Servidor U-Go (Go) corriendo en http://localhost:%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, enhancedRouter))

	// (Aquí conectaríamos a Mongo también)
	// ...
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
