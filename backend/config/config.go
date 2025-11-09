package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Esta carpeta config almacena toda la configuración de la aplicación

type Config struct {
	DatabaseURL string
	Port        string
	MongoURI    string
}

// LoadConfig carga la configuracion desde un archivo .env
func LoadConfig() (*Config, error) {
	// carga el archivo.env desde la raiz dek proyecto
	// (Asume que .env esta en el repo principal, un nivel arribba de backend)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: The .env file could not be loaded. Using system environment variables.")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		MongoURI:    os.Getenv("MONGO_URI"),
	}
	// Valores por defecto
	if cfg.Port == "" {
		cfg.Port = "8080" // Puerto por defecto
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return cfg, nil
}
