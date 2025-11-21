// Define la estructura de las coordenadas que devolverá el servicio
export interface Coordinates {
    latitude: number;
    longitude: number;
    accuracy: number;
    timestamp: number;
}

/**
 * Servicio para obtener la geolocalización del dispositivo.
 * Por ahora, usa la API nativa del navegador.
 */
export const GeolocationService = {

    /**
     * Obtiene la posición actual del dispositivo (una sola vez).
     * @returns Promise<Coordinates>
     */
    getCurrentPosition(): Promise<Coordinates> {
        return new Promise((resolve, reject) => {
            if (!navigator.geolocation) {
                reject(new Error("La geolocalización no es compatible con este navegador o dispositivo."));
                return;
            }

            navigator.geolocation.getCurrentPosition(
                (position) => {
                    const coords: Coordinates = {
                        latitude: position.coords.latitude,
                        longitude: position.coords.longitude,
                        accuracy: position.coords.accuracy,
                        timestamp: position.timestamp,
                    };
                    resolve(coords);
                },
                (error) => {
                    let errorMessage: string;
                    switch (error.code) {
                        case error.PERMISSION_DENIED:
                            errorMessage = "Permiso de ubicación denegado por el usuario.";
                            break;
                        case error.POSITION_UNAVAILABLE:
                            errorMessage = "Información de ubicación no disponible.";
                            break;
                        case error.TIMEOUT:
                            errorMessage = "La solicitud para obtener la ubicación ha caducado.";
                            break;
                        default:
                            errorMessage = `Error de geolocalización desconocido (Código: ${error.code}).`;
                            break;
                    }
                    reject(new Error(errorMessage));
                },
                // Opciones de configuración
                {
                    enableHighAccuracy: true,
                    timeout: 5000,           
                    maximumAge: 0,           
                }
            );
        });
    },

    /**
     * Inicia el monitoreo continuo de la posición.
     * @param callback Función que se llama con las nuevas coordenadas.
     * @param errorCallback Función que se llama en caso de error.
     * @returns ID del watcher para detener el monitoreo.
     */
    watchPosition(
        callback: (coords: Coordinates) => void,
        errorCallback: (error: Error) => void
    ): number | null {
        if (!navigator.geolocation) {
            errorCallback(new Error("La geolocalización no es compatible."));
            return null;
        }

        const watchId = navigator.geolocation.watchPosition(
            (position) => {
                const coords: Coordinates = {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude,
                    accuracy: position.coords.accuracy,
                    timestamp: position.timestamp,
                };
                callback(coords);
            },
            (error) => {
                let errorMessage: string;
                switch (error.code) {
                    case error.PERMISSION_DENIED:
                        errorMessage = "Permiso de ubicación denegado por el usuario.";
                        break;
                    case error.POSITION_UNAVAILABLE:
                        errorMessage = "Información de ubicación no disponible.";
                        break;
                    case error.TIMEOUT:
                        errorMessage = "La solicitud para obtener la ubicación ha caducado.";
                        break;
                    default:
                        errorMessage = `Error de geolocalización desconocido (Código: ${error.code}).`;
                        break;
                }
                errorCallback(new Error(errorMessage));
            },
            {
                enableHighAccuracy: true,
                timeout: 10000, 
                maximumAge: 0, 
            }
        );

        return watchId;
    },

    /**
     * Detiene el monitoreo de posición.
     * @param watchId ID devuelto por watchPosition.
     */
    clearWatch(watchId: number | null): void {
        if (watchId !== null && navigator.geolocation) {
            navigator.geolocation.clearWatch(watchId);
        }
    }
};