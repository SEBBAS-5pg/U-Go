// backend/internal/api/response_utils.go
package api

import (
	"encoding/json"
	"net/http"
)

// RespondWithError maneja las respuestas de error estandarizadas.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON maneja las respuestas de éxito estandarizadas.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// 1. Establecer el tipo de contenido
	w.Header().Set("Content-Type", "application/json")

	// 2. Establecer el código de estado HTTP
	w.WriteHeader(code)

	// 3. Codificar el payload (datos) a JSON y escribir la respuesta
	if payload != nil {
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			// Si hay un error de codificación, simplemente lo logueamos o lo ignoramos.
			// No podemos cambiar el código de estado en este punto.
			// log.Printf("Error al codificar respuesta JSON: %v", err)
		}
	}
}

// ContextKey define un tipo para las claves de contexto, previniendo colisiones.
// Aunque el trip_handler usa models.ContextUserIDKey, esta función estaba en el código anterior
// y se define aquí por si otras partes del código la necesitan, aunque la eliminamos del CancelTrip corregido.
type ContextKey string
