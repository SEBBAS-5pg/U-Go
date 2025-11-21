import React, { useState, useEffect } from 'react';
import { 
  IonContent, IonHeader, IonPage, IonTitle, IonToolbar, 
  IonCard, IonCardHeader, IonCardTitle, IonCardContent,
  IonItem, IonLabel, IonNote, IonBadge, IonIcon, IonButton, 
  IonInput, IonList, IonToast
} from '@ionic/react';
import { navigate, search, copy } from 'ionicons/icons'; 
import TripMap from '../components/MapComponent'; 
import { LocationReceiveService } from '../services/LocationReceiveService';
import { TripService, CreateTripPayload } from '../services/TripService';
import './Tab2.css';

const DEFAULT_CITY_LAT = 4.6097; 
const DEFAULT_CITY_LNG = -74.0817;

const Tab2: React.FC = () => {
    const [currentTripId, setCurrentTripId] = useState<string | null>(null);
    
    const [origin, setOrigin] = useState("Mi Casa");
    const [destination, setDestination] = useState("Oficina");
    const [isCreating, setIsCreating] = useState(false);
    const [toastMessage, setToastMessage] = useState<string | null>(null);

    const [driverLocation, setDriverLocation] = useState<[number, number] | null>(null);
    const [pollingError, setPollingError] = useState<string | null>(null);
    const [lastReceivedTime, setLastReceivedTime] = useState<Date | null>(null);

    // 0. 🧠 PERSISTENCIA: Cargar ID al iniciar
    useEffect(() => {
        const savedTripId = localStorage.getItem('passenger_active_trip_id');
        if (savedTripId) {
            console.log("Restaurando sesión de pasajero:", savedTripId);
            setCurrentTripId(savedTripId);
        }
    }, []);

    // CREAR VIAJE
    const handleRequestTrip = async () => {
        if (!origin || !destination) {
            setToastMessage("Por favor ingresa origen y destino");
            return;
        }
        setIsCreating(true);
        try {
            const payload: CreateTripPayload = {
                origin_name: origin,
                origin_lat: DEFAULT_CITY_LAT,
                origin_lng: DEFAULT_CITY_LNG,
                destination_name: destination,
                destination_lat: DEFAULT_CITY_LAT + 0.05, 
                destination_lng: DEFAULT_CITY_LNG + 0.02,
                conductor_id: "" 
            };

            const newTrip = await TripService.createTrip(payload);
            setCurrentTripId(newTrip.id);
            
            // 💾 GUARDAR EN LOCALSTORAGE
            localStorage.setItem('passenger_active_trip_id', newTrip.id);
            
            setToastMessage("¡Viaje solicitado! Copia el ID para el conductor.");

        } catch (error: any) {
            console.error(error);
            setToastMessage(error.message || "Error al crear viaje");
        } finally {
            setIsCreating(false);
        }
    };

    // ACCIÓN: NUEVO VIAJE (LIMPIAR)
    const handleNewTrip = () => {
        setCurrentTripId(null);
        setDriverLocation(null);
        setLastReceivedTime(null);
        // 🗑️ BORRAR DE LOCALSTORAGE
        localStorage.removeItem('passenger_active_trip_id');
    };

    // POLLING
    useEffect(() => {
        if (!currentTripId) return;
        const fetchLocation = async () => {
            try {
                const locationData = await LocationReceiveService.getDriverLocation(currentTripId);
                setDriverLocation([locationData.latitude, locationData.longitude]);
                setPollingError(null);
                setLastReceivedTime(new Date());
            } catch (err) { console.log("Esperando ubicación..."); }
        };
        fetchLocation();
        const intervalId = setInterval(fetchLocation, 3000);
        return () => clearInterval(intervalId);
    }, [currentTripId]);

    // VISTA FORMULARIO
    const renderRequestForm = () => (
        <div className="ion-padding">
            <IonCard>
                <IonCardHeader><IonCardTitle>Pedir Viaje</IonCardTitle></IonCardHeader>
                <IonCardContent>
                    <IonList>
                        <IonItem lines="none">
                            <IonIcon icon={navigate} slot="start" color="primary"/>
                            <IonInput label="Origen" labelPlacement="floating" value={origin} onIonInput={e => setOrigin(e.detail.value!)}/>
                        </IonItem>
                        <IonItem lines="none">
                            <IonIcon icon={search} slot="start" color="secondary"/>
                            <IonInput label="Destino" labelPlacement="floating" value={destination} onIonInput={e => setDestination(e.detail.value!)}/>
                        </IonItem>
                    </IonList>
                    <div className="ion-padding-top">
                        <IonButton expand="block" onClick={handleRequestTrip} disabled={isCreating}>
                            {isCreating ? "Solicitando..." : "Solicitar U-Go"}
                        </IonButton>
                    </div>
                </IonCardContent>
            </IonCard>
        </div>
    );

    // VISTA MAPA
    const renderTripMonitor = () => {
        const driverCoords = driverLocation ? [driverLocation[0], driverLocation[1]] as [number, number] : undefined;
        const myCoords = [DEFAULT_CITY_LAT, DEFAULT_CITY_LNG] as [number, number];

        return (
            <>
                <TripMap 
                    driverLocation={driverCoords}
                    passengerLocation={myCoords} 
                    center={myCoords} 
                    zoom={13}
                />
                
                <IonCard className="id-copy-box">
                    <IonCardContent>
                        <IonItem lines="none">
                            <IonLabel position="stacked" color="medium">ID para el conductor:</IonLabel>
                            <IonInput readonly value={currentTripId || ""} />
                            <IonButton slot="end" fill="clear" onClick={() => {
                                if(currentTripId) {
                                    navigator.clipboard.writeText(currentTripId);
                                    setToastMessage("¡ID copiado!");
                                }
                            }}>
                                <IonIcon icon={copy} />
                            </IonButton>
                        </IonItem>
                    </IonCardContent>
                </IonCard>

                <IonCard color="light" className="ion-margin-bottom">
                    <IonCardHeader>
                        <IonCardTitle className="ion-text-center">Estado del Viaje</IonCardTitle>
                    </IonCardHeader>
                    <IonCardContent className="ion-text-center">
                        <p className="status-text">Esperando conductor...</p>
                        {lastReceivedTime && (
                            <IonNote color="success" style={{fontWeight: 'bold'}}>
                                <br/>¡Conductor en movimiento!<br/>Actualizado: {lastReceivedTime.toLocaleTimeString()}
                            </IonNote>
                        )}
                    </IonCardContent>
                </IonCard>

                <IonButton expand="block" color="medium" fill="outline" onClick={handleNewTrip}>
                    Nuevo Viaje / Salir
                </IonButton>
            </>
        );
    };

    return (
        <IonPage>
            <IonHeader><IonToolbar color="secondary"><IonTitle>Pasajero</IonTitle></IonToolbar></IonHeader>
            <IonContent fullscreen>
                {currentTripId ? renderTripMonitor() : renderRequestForm()}
                <IonToast isOpen={!!toastMessage} message={toastMessage || ""} duration={2000} onDidDismiss={() => setToastMessage(null)} />
            </IonContent>
        </IonPage>
    );
};

export default Tab2;