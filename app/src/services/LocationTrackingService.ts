import { Coordinates } from './GeolocationService';

// Define la URL base de nuestro backend
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Interfaz para el payload que se enviará al backend
 */
interface LocationPayload {
    latitude: number;
    longitude: number;
    timestamp: number;
    accuracy: number;
}

/**
 * Servicio para enviar la ubicación del dispositivo al backend
 */
export const LocationTrackingService = {
    
    /**
     * Envia la ubicación del conductor al servidor para un viaje específico.
     * * @param tripId El ID del viaje que se está monitoreando.
     * @param coords Las coordenadas de la ubicación actual.
     * @returns Promesa que resuelve en void si es exitoso.
     */
    async sendLocationUpdate(tripId: string, coords: Coordinates): Promise<void> {
        const url = `${API_BASE_URL}/trip/${tripId}/location`;
        
        // Creamos el payload con la estructura requerida por el backend
        const payload: LocationPayload = {
            latitude: coords.latitude,
            longitude: coords.longitude,
            timestamp: coords.timestamp,
            accuracy: coords.accuracy
        };

        try {
            // Utilizamos el método POST o PUT (dependiendo de la implementación del backend)
            // Usaremos POST por simplicidad, asumiendo que es un nuevo 'evento' de ubicación
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    // Aquí iría la autenticación (Bearer Token, si estuviera implementada)
                },
                body: JSON.stringify(payload),
            });

            // Si la respuesta no es 2xx, lanzamos un error
            if (!response.ok) {
                const errorText = await response.text();
                // Si el backend envía un mensaje de error claro, lo usamos
                let errorMessage = `Error al actualizar la ubicación para el viaje ${tripId}. Estado: ${response.status}`;
                if (errorText) {
                    // Intentamos parsear el JSON si existe
                    try {
                        const errorJson = JSON.parse(errorText);
                        if (errorJson.message) {
                            errorMessage = errorJson.message;
                        } else {
                            errorMessage = errorText;
                        }
                    } catch {
                        // Si no es JSON, usamos el texto plano
                        errorMessage = errorText;
                    }
                }
                throw new Error(errorMessage);
            }
            
            // Si la respuesta es 200/204, se considera exitoso y no devolvemos nada
            //console.log(`Ubicación de viaje ${tripId} actualizada exitosamente.`);

        } catch (error) {
            // Capturamos errores de red o errores lanzados desde el bloque try
            console.error("Fallo en LocationTrackingService:", error);
            if (error instanceof Error) {
                 throw new Error(`Fallo en la comunicación: ${error.message}`);
            }
             throw new Error("Error desconocido al enviar ubicación.");
        }
    },
};