import React, { useState, useEffect, useRef } from 'react';
import { 
  IonContent, 
  IonHeader, 
  IonPage, 
  IonTitle, 
  IonToolbar, 
  IonLoading, // Lo usamos aquí, fuera de renderMainContent
  IonCard, 
  IonCardHeader, 
  IonCardTitle, 
  IonCardContent,
  IonItem,
  IonLabel,
  IonNote,
  IonIcon,
  IonBadge
} from '@ionic/react';
import { location, navigate, paperPlane, warning } from 'ionicons/icons'; 
import { TripService, TripData } from '../services/TripService';
import { GeolocationService, Coordinates } from '../services/GeolocationService'; 
import { LocationTrackingService } from '../services/LocationTrackingService'; 
import './Tab1.css';

// ⚠️ SIMULACIÓN DE ID DE VIAJE ⚠️
const TRIP_ID_MOCK = "85b57f0f-8979-4591-b0e2-d4b998f45a0b"; 

// Frecuencia de envío de ubicación al servidor (en milisegundos)
const UPDATE_INTERVAL_MS = 10000; // 10 segundos

const Tab1: React.FC = () => {
  // Estados del Viaje (Backend)
  const [tripData, setTripData] = useState<TripData | null>(null);
  const [loading, setLoading] = useState<boolean>(true); // Se inicializa en true para la carga inicial
  const [error, setError] = useState<string | null>(null);
  
  // Estados del GPS (Frontend/Nativo)
  const [currentLocation, setCurrentLocation] = useState<Coordinates | null>(null);
  const [gpsError, setGpsError] = useState<string | null>(null);
  
  // Estado de envío de ubicación al servidor
  const [lastSentTime, setLastSentTime] = useState<Date | null>(null);
  const [isSending, setIsSending] = useState<boolean>(false);
  const [sendError, setSendError] = useState<string | null>(null);

  // Referencia para guardar la última ubicación recibida del GPS
  const latestCoords = useRef<Coordinates | null>(null);

  // 1. Efecto para cargar datos del viaje (Backend)
  useEffect(() => {
    const loadTripData = async () => {
      setLoading(true); // 👈 Iniciamos la carga
      setError(null);
      try {
        const data = await TripService.getTripData(TRIP_ID_MOCK);
        setTripData(data);
        // NO PONEMOS setLoading(false) aquí, lo hacemos en el finally
      } catch (err) {
        console.error("Error al cargar el viaje:", err);
        if (err instanceof Error) {
            setError(err.message);
        } else {
            setError("Error desconocido.");
        }
      } finally {
        // 🚨 CRUCIAL: Detenemos la carga SÍ o SÍ, ya sea que haya éxito o error.
        setLoading(false); 
      }
    };
    loadTripData();
  }, []);

  // 2. Efecto para el GPS (Geolocalización en tiempo real)
  useEffect(() => {
    const watchId = GeolocationService.watchPosition(
        (coords) => {
            setCurrentLocation(coords);
            latestCoords.current = coords; 
            setGpsError(null); 
        },
        (err) => {
            setGpsError(err.message);
        }
    );

    return () => {
        if (watchId !== null) {
            GeolocationService.clearWatch(watchId);
        }
    };
  }, []);

  // 3. Efecto para el Envío periódico de ubicación (Backend Update)
  useEffect(() => {
    const sendLocation = async () => {
        if (!latestCoords.current) {
            return;
        }

        setIsSending(true);
        setSendError(null);
        try {
            await LocationTrackingService.sendLocationUpdate(TRIP_ID_MOCK, latestCoords.current);
            setLastSentTime(new Date());
            setSendError(null); 
        } catch (err) {
            console.error("Fallo al enviar ubicación al servidor:", err);
            if (err instanceof Error) {
                setSendError(err.message);
            } else {
                setSendError("Error desconocido en el envío.");
            }
        } finally {
            setIsSending(false);
        }
    };

    const intervalId = setInterval(sendLocation, UPDATE_INTERVAL_MS);

    return () => {
        clearInterval(intervalId);
    };
  }, []);


  // --- Renderizado Condicional del Contenido Principal ---
  const renderMainContent = () => {
    // 🚨 Ya NO comprobamos 'loading' aquí, esa lógica se movió al final del componente.

    // Si hay error en el viaje, mostramos mensaje amigable
    if (error || !tripData) {
        return (
            <div className="error-container ion-padding ion-text-center">
                <IonIcon icon={location} color="medium" style={{ fontSize: '64px' }} />
                <h3>No se pudo cargar el viaje</h3>
                <p className="ion-text-wrap text-medium">
                    {error || "Verifica que el ID del viaje sea correcto."}
                </p>
                <IonNote color="danger" className="id-note">{TRIP_ID_MOCK}</IonNote>
            </div>
        );
    }

    // Si hay datos, mostramos la info del viaje
    const getStatusColor = (status: string) => {
        switch (status) {
          case 'aceptado': return 'success';
          case 'en_curso': return 'success';
          case 'solicitado': return 'warning';
          case 'finalizado': return 'tertiary'; 
          case 'cancelado': return 'danger'; 
          default: return 'medium';
        }
    };

    return (
        <>
            <IonCard color="light" className="ion-margin-bottom">
                <IonCardHeader>
                    <IonCardTitle className="ion-text-center">Estado Actual</IonCardTitle>
                    <div className="ion-text-center ion-padding-top">
                        <IonBadge color={getStatusColor(tripData.status)} style={{ fontSize: '1.2em', padding: '8px' }}>
                            {tripData.status.toUpperCase().replace('_', ' ')}
                        </IonBadge>
                    </div>
                </IonCardHeader>
            </IonCard>

            <IonCard>
                <IonCardContent>
                    <IonItem lines="full">
                        <IonLabel>Origen</IonLabel>
                        <IonNote slot="end">{tripData.originName || "Coordenadas GPS"}</IonNote>
                    </IonItem>
                    <IonItem lines="none">
                        <IonLabel>Destino</IonLabel>
                        <IonNote slot="end">{tripData.destinationName || "Coordenadas GPS"}</IonNote>
                    </IonItem>
                    
                    <div className="ion-padding-top">
                        <IonNote style={{ fontSize: '0.8em', display: 'block' }}>Conductor: {tripData.conductorId ? 'Asignado' : 'Pendiente'}</IonNote>
                        <IonNote style={{ fontSize: '0.8em', display: 'block' }}>Vehículo: {tripData.vehicleId ? 'Asignado' : 'Pendiente'}</IonNote>
                    </div>
                </IonCardContent>
            </IonCard>
        </>
    );
  };

  // --- Renderizado de la Tarjeta de Envío de Datos ---
  const renderTrackingCard = () => {
    // ... (El resto del código de la tarjeta de tracking es el mismo)
    let color: string;
    let icon: string;
    let message: string;
    
    if (sendError) {
        color = 'danger';
        icon = warning;
        message = sendError.includes("404") ? 
            "⚠️ Error 404: Ruta de envío no encontrada." : 
            `Error de red: ${sendError}`;
    } else if (isSending) {
        color = 'warning';
        icon = paperPlane;
        message = 'Enviando ubicación...';
    } else if (lastSentTime) {
        color = 'success';
        icon = paperPlane;
        message = `Último envío: ${lastSentTime.toLocaleTimeString()}`;
    } else {
        color = 'medium';
        icon = paperPlane;
        message = 'Esperando el primer envío...';
    }

    return (
        <IonCard className={`tracking-card ion-text-center card-${color}`}>
            <IonCardContent>
                <IonIcon icon={icon} color={color} style={{ fontSize: '24px' }} />
                <p style={{ marginTop: '5px' }}>
                    <IonNote color={color}>{message}</IonNote>
                </p>
            </IonCardContent>
        </IonCard>
    );
  }


  return (
    <IonPage>
      <IonHeader>
        <IonToolbar color="primary">
          <IonTitle>Monitoreo Activo</IonTitle>
        </IonToolbar>
      </IonHeader>
      
      <IonContent fullscreen className="ion-padding">
        
        {/* 🚨 IonLoading se renderiza AQUÍ y SÓLO si loading es true */}
        <IonLoading isOpen={loading} message={'Cargando datos del viaje...'} />
        
        {/* 🚨 Si NO estamos cargando, mostramos el contenido */}
        {!loading && (
          <>
            {/* --- TARJETA DE ENVÍO AL BACKEND --- */}
            {renderTrackingCard()}
            
            {/* --- TARJETA DE GPS (Datos Nativos) --- */}
            <IonCard className="gps-card">
                <IonCardHeader>
                    <IonCardTitle className="gps-title">
                        <IonIcon icon={navigate} /> Mi Ubicación GPS
                    </IonCardTitle>
                </IonCardHeader>
                <IonCardContent>
                    {currentLocation ? (
                        <div className="coords-display">
                            <div className="coord-item">
                                <span className="label">Latitud:</span>
                                <span className="value">{currentLocation.latitude.toFixed(6)}</span>
                            </div>
                            <div className="coord-item">
                                <span className="label">Longitud:</span>
                                <span className="value">{currentLocation.longitude.toFixed(6)}</span>
                            </div>
                            <div className="coord-item">
                                <span className="label">Precisión:</span>
                                <span className="value">+/- {currentLocation.accuracy.toFixed(1)} m</span>
                            </div>
                            <div className="last-update">
                                Actualizado: {new Date(currentLocation.timestamp).toLocaleTimeString()}
                            </div>
                        </div>
                    ) : (
                        <div className="ion-text-center ion-padding">
                            {gpsError ? (
                                <IonNote color="danger">{gpsError}</IonNote>
                            ) : (
                                <div className="scanning-text">
                                    <IonLoading isOpen={true} duration={0} message={undefined} cssClass="spinner-inline"/>
                                    Buscando satélites...
                                </div>
                            )}
                        </div>
                    )}
                </IonCardContent>
            </IonCard>

            {/* --- DATOS DEL BACKEND (Viaje) --- */}
            {renderMainContent()}
          </>
        )}
      </IonContent>
    </IonPage>
  );
};

export default Tab1;