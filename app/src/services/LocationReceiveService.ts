// Define la URL base de nuestro backend
const API_BASE_URL = 'http://localhost:8080/api/v1';

/**
 * Interfaz para la ubicación que recibiremos del backend.
 */
export interface DriverLocation {
    latitude: number;
    longitude: number;
    timestamp: number;
    tripId: string;
    // Podríamos añadir más datos como la velocidad, etc.
}

/**
 * Servicio para obtener la ubicación del conductor del backend.
 */
export const LocationReceiveService = {
    
    /**
     * Consulta la última ubicación conocida del conductor para un viaje específico.
     * @param tripId El ID del viaje a monitorear.
     * @returns Promesa que resuelve con la ubicación del conductor.
     */
    async getDriverLocation(tripId: string): Promise<DriverLocation> {
        const url = `${API_BASE_URL}/trip/${tripId}/location`;
        
        try {
            // Utilizamos el método GET para obtener la última ubicación
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) {
                // Si el backend aún no tiene la ubicación, puede devolver 404/400.
                // Lo manejamos lanzando un error.
                let errorMessage = `Error al obtener ubicación para el viaje ${tripId}. Estado: ${response.status}`;
                if (response.status === 404) {
                    errorMessage = `Viaje no encontrado o conductor sin enviar ubicación (404).`;
                }
                throw new Error(errorMessage);
            }
            
            // Suponemos que la respuesta es un JSON con la estructura DriverLocation
            const data: DriverLocation = await response.json();
            return data;

        } catch (error) {
            console.error("Fallo en LocationReceiveService:", error);
            if (error instanceof Error) {
                 throw new Error(`Fallo en la comunicación: ${error.message}`);
            }
             throw new Error("Error desconocido al obtener ubicación.");
        }
    },
};