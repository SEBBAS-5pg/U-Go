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


export const TripService = {
    getTripData,
};