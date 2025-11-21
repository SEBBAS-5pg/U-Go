// src/services/TripService.ts
import axios, { isAxiosError } from "axios";
import { AuthService } from "./AuthService"; 

const BASE_URL = import.meta.env.VITE_API_BASE_URL;

// --- Interfaz de la respuesta de GET /trips/{tripId} ---
// Basada estrictamente en el modelo models.Trip de Go
export interface TripData {
    id: string;
    pasajeroId: string;
    conductorId: string; 
    vehicleId: string;   
    status: 'solicitado' | 'aceptado' | 'en_curso' | 'finalizado' | 'cancelado';
    
    originLat: number;
    originLng: number;
    originName?: string; 
    destinationLat: number;
    destinationLng: number;
    destinationName?: string;
    
    createdAt: string; 
    startedAt?: string;
    completedAt?: string; 
}

// Estructura para solicitar un viaje
export interface CreateTripPayload {
    origin_name: string;
    origin_lat: number;
    origin_lng: number;
    destination_name: string;
    destination_lat: number;
    destination_lng: number;
    conductor_id?: string; // Opcional por ahora
}

// Llama al endpoint GET /trips/{tripId}
async function getTripData(tripId: string): Promise<TripData> {
    
    const token = AuthService.getToken();

    if (!token) {
        throw new Error("Usuario no autenticado.");
    }

    try {
        // Usamos TripData como el tipo de respuesta esperado
        const response = await axios.get<TripData>(
            `${BASE_URL}/trips/${tripId}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        
        return response.data;
        
    } catch (error) {
        if (isAxiosError(error) && error.response) {
            throw new Error(error.response.data.error || `Error al obtener datos del viaje ${tripId}`);
        }
        throw new Error("Error de red o desconocido.");
    }
}

async function createTrip(payload: CreateTripPayload): Promise<TripData> {
    const token = AuthService.getToken();
    if (!token) throw new Error("No hay sesión activa");

    try {
        const response = await axios.post<TripData>(
            `${BASE_URL}/trips`,
            payload,
            {
                headers: { Authorization: `Bearer ${token}` }
            }
        );
        return response.data;
    } catch (error) {
        if (isAxiosError(error) && error.response) {
            throw new Error(error.response.data.error || "Error al solicitar viaje");
        }
        throw new Error("Error de red al crear viaje");
    }
}
// Función para que un conductor acepte un viaje
async function acceptTrip(tripId: string, vehicleId: string): Promise<TripData> {
    const token = AuthService.getToken();
    if (!token) throw new Error("No hay sesión activa");

    try {
        // El backend espera { vehicle_id: "..." }
        const payload = { vehicle_id: vehicleId };
        
        const response = await axios.post<TripData>(
            `${BASE_URL}/trips/${tripId}/accept`,
            payload,
            {
                headers: { Authorization: `Bearer ${token}` }
            }
        );
        return response.data;
    } catch (error) {
        if (isAxiosError(error) && error.response) {
            throw new Error(error.response.data.error || "Error al aceptar el viaje");
        }
        throw new Error("Error de red al aceptar viaje");
    }
}

export const TripService = {
    getTripData,
    createTrip,
    acceptTrip, // <--- ¡No olvides exportarla!
};