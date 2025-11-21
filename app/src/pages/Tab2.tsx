import React, { useState, useEffect } from 'react';
import { 
  IonContent, 
  IonHeader, 
  IonPage, 
  IonTitle, 
  IonToolbar, 
  IonCard,
  IonCardHeader,
  IonCardTitle,
  IonCardContent,
  IonItem,
  IonLabel,
  IonNote,
  IonBadge,
  IonIcon
} from '@ionic/react';
import { warning } from 'ionicons/icons';
import TripMap from '../components/MapComponent'; 
import { LocationReceiveService, DriverLocation } from '../services/LocationReceiveService'; // Importamos el nuevo servicio
import './Tab2.css';

// ⚠️ SIMULACIÓN DE ID DE VIAJE (Pasajero) ⚠️
const TRIP_ID_MOCK = "85b57f0f-8979-4591-b0e2-d4b998f45a0b"; 

// Frecuencia de Polling (obtener ubicación del conductor)
const POLLING_INTERVAL_MS = 3000; // Cada 3 segundos

const mockTripData = {
    id: TRIP_ID_MOCK,
    status: 'en_curso', 
    originName: 'Plaza Cataluña',
    destinationName: 'Parque Güell',
    conductorName: 'Juan Pérez',
    vehicleModel: 'Seat Arona',
    vehiclePlate: 'ABC-123',
};

// Coordenadas iniciales para el mapa (simulación)
const PASSENGER_COORDS: [number, number] = [41.3851, 2.1734];


const Tab2: React.FC = () => {
    // Estado para la ubicación del conductor recibida del backend
    const [driverLocation, setDriverLocation] = useState<[number, number] | null>(null);
    const [pollingError, setPollingError] = useState<string | null>(null);
    const [lastReceivedTime, setLastReceivedTime] = useState<Date | null>(null);

    // 1. Efecto para el Polling (Consulta de Ubicación en Tiempo Real)
    useEffect(() => {
        const fetchLocation = async () => {
            try {
                // 1. Llamar al servicio para obtener la ubicación
                const locationData: DriverLocation = await LocationReceiveService.getDriverLocation(TRIP_ID_MOCK);
                
                // 2. Actualizar el estado con las coordenadas
                setDriverLocation([locationData.latitude, locationData.longitude]);
                setPollingError(null);
                setLastReceivedTime(new Date());

            } catch (err) {
                // Si hay error (ej: 404 o backend caído), lo almacenamos
                const errorMsg = err instanceof Error ? err.message : "Error desconocido al hacer polling.";
                setPollingError(errorMsg);
                console.error("Error de Polling:", errorMsg);
            }
        };

        // Iniciar el Polling: ejecutar inmediatamente y luego cada intervalo
        fetchLocation(); // Primera ejecución inmediata
        const intervalId = setInterval(fetchLocation, POLLING_INTERVAL_MS);

        // Limpieza: Detener el intervalo al desmontar el componente
        return () => {
            clearInterval(intervalId);
            console.log("Polling de ubicación del conductor detenido.");
        };
    }, []); // El array vacío asegura que solo se ejecute una vez al montar


    const getStatusColor = (status: string) => {
        switch (status) {
          case 'en_curso': return 'success';
          case 'solicitado': return 'warning';
          case 'finalizado': return 'tertiary'; 
          case 'cancelado': return 'danger'; 
          default: return 'medium';
        }
    };
    
    // Convertimos la ubicación recibida (si existe) al formato de MapComponent
    const driverCoords: [number, number] | undefined = driverLocation ? [driverLocation[0], driverLocation[1]] : undefined;

  return (
    <IonPage>
      <IonHeader>
        <IonToolbar color="secondary">
          <IonTitle>Ruta del Viaje</IonTitle>
        </IonToolbar>
      </IonHeader>
      
      <IonContent fullscreen className="ion-padding">

        {/* --- MAPA DE RASTREO --- */}
        {/* Le pasamos la ubicación del conductor (real o nula) y la del pasajero (mock) */}
        <TripMap 
            driverLocation={driverCoords}
            passengerLocation={PASSENGER_COORDS} // La ubicación del pasajero es mockeada aquí
            center={PASSENGER_COORDS} // Centramos el mapa en el pasajero
            zoom={14}
        />
        
        {/* --- ESTADO DEL POLLING --- */}
        <IonCard color={pollingError ? "danger" : "success"} className="ion-margin-bottom ion-text-center">
            <IonCardContent className="ion-no-padding">
                <IonItem color={pollingError ? "danger" : "success"} lines="none">
                    <IonIcon icon={warning} slot="start" />
                    <IonLabel className="ion-text-wrap">
                        {pollingError ? (
                            <IonNote color="light" style={{ fontWeight: 'bold' }}>
                                {pollingError}
                            </IonNote>
                        ) : (
                            <IonNote color="light">
                                {lastReceivedTime ? `Última ubicación recibida: ${lastReceivedTime.toLocaleTimeString()}` : "Conectando al servidor..."}
                            </IonNote>
                        )}
                    </IonLabel>
                </IonItem>
            </IonCardContent>
        </IonCard>

        {/* --- DETALLES DEL VIAJE --- */}
        <IonCard color="light" className="ion-margin-bottom">
            <IonCardHeader>
                <IonCardTitle className="ion-text-center">Viaje Activo</IonCardTitle>
                <div className="ion-text-center ion-padding-top">
                    <IonBadge color={getStatusColor(mockTripData.status)} style={{ fontSize: '1.2em', padding: '8px' }}>
                        {mockTripData.status.toUpperCase().replace('_', ' ')}
                    </IonBadge>
                </div>
            </IonCardHeader>
        </IonCard>

        <IonCard>
            <IonCardHeader>
                <IonCardTitle>Detalles del Conductor</IonCardTitle>
            </IonCardHeader>
            <IonCardContent>
                <IonItem lines="full">
                    <IonLabel>Conductor:</IonLabel>
                    <IonNote slot="end" color="dark">{mockTripData.conductorName}</IonNote>
                </IonItem>
                <IonItem lines="full">
                    <IonLabel>Vehículo:</IonLabel>
                    <IonNote slot="end" color="dark">{mockTripData.vehicleModel}</IonNote>
                </IonItem>
                <IonItem lines="none">
                    <IonLabel>Placa:</IonLabel>
                    <IonNote slot="end" color="dark">{mockTripData.vehiclePlate}</IonNote>
                </IonItem>
            </IonCardContent>
        </IonCard>

        <div className="ion-padding-top ion-text-center">
            <IonNote color="medium">Tu ubicación y la del conductor se actualizan en tiempo real.</IonNote>
        </div>

      </IonContent>
    </IonPage>
  );
};

export default Tab2;