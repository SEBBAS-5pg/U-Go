import React, { useState, useEffect, useRef } from 'react';
import { 
  IonContent, IonHeader, IonPage, IonTitle, IonToolbar, 
  IonCard, IonCardHeader, IonCardTitle, IonCardContent,
  IonItem, IonLabel, IonNote, IonIcon, IonBadge, IonButton, IonInput, IonToast
} from '@ionic/react';
import { navigate, paperPlane, car } from 'ionicons/icons'; 
import { TripService, TripData } from '../services/TripService';
import { GeolocationService, Coordinates } from '../services/GeolocationService'; 
import { LocationTrackingService } from '../services/LocationTrackingService'; 
import './Tab1.css';

const UPDATE_INTERVAL_MS = 10000; 

const Tab1: React.FC = () => {
  const [activeTripId, setActiveTripId] = useState<string | null>(null);
  const [inputTripId, setInputTripId] = useState(""); 
  const [tripData, setTripData] = useState<TripData | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  const [currentLocation, setCurrentLocation] = useState<Coordinates | null>(null);
  const [isSending, setIsSending] = useState<boolean>(false);
  const [sendError, setSendError] = useState<string | null>(null);
  const latestCoords = useRef<Coordinates | null>(null);

  // 0. 🧠 PERSISTENCIA: Cargar ID al iniciar
  useEffect(() => {
    const savedTripId = localStorage.getItem('driver_active_trip_id');
    if (savedTripId) {
        console.log("Restaurando sesión de conductor:", savedTripId);
        setActiveTripId(savedTripId);
    }
  }, []);

  // 1. GPS Activo
  useEffect(() => {
    const watchId = GeolocationService.watchPosition(
        (coords) => { setCurrentLocation(coords); latestCoords.current = coords; },
        (err) => console.error(err)
    );
    return () => { if (watchId !== null) GeolocationService.clearWatch(watchId); };
  }, []);

  // 2. ACCIÓN: ACEPTAR VIAJE
  const handleAcceptTrip = async () => {
    const cleanId = inputTripId.replace(/[^a-f0-9-]/gi, '');

    if (cleanId.length < 10) {
        setToastMsg("El ID del viaje parece inválido (muy corto)");
        return;
    }

    setLoading(true);
    setError(null);
    try {
        const trip = await TripService.acceptTrip(cleanId, "V-CUALQUIERA");
        setTripData(trip);
        setActiveTripId(trip.id);
        
        // 💾 GUARDAR EN LOCALSTORAGE
        localStorage.setItem('driver_active_trip_id', trip.id);
        
        setToastMsg("¡Conectado! Enviando ubicación...");
    } catch (err: any) {
        const msg = err.message || "Error desconocido";
        setError(msg);
        setToastMsg("Error: " + msg);
    } finally {
        setLoading(false);
    }
  };

  // ACCIÓN: SALIR / DETENER
  const handleStopTrip = () => {
      setActiveTripId(null);
      setTripData(null);
      // 🗑️ BORRAR DE LOCALSTORAGE
      localStorage.removeItem('driver_active_trip_id');
  };

  // 3. MONITOREO
  useEffect(() => {
    if (!activeTripId) return;
    TripService.getTripData(activeTripId).then(setTripData).catch(console.error);
  }, [activeTripId]);

  // 4. TRACKING
  useEffect(() => {
    if (!activeTripId) return;
    const sendLocation = async () => {
        if (!latestCoords.current) return;
        setIsSending(true);
        try {
            await LocationTrackingService.sendLocationUpdate(activeTripId, latestCoords.current);
        } catch (e) { 
            if(e instanceof Error) setSendError(e.message); 
        } finally { setIsSending(false); }
    };
    const interval = setInterval(sendLocation, UPDATE_INTERVAL_MS);
    sendLocation();
    return () => clearInterval(interval);
  }, [activeTripId]);

  // --- RENDERIZADO ---
  return (
    <IonPage>
      <IonHeader><IonToolbar color="primary"><IonTitle>Conductor</IonTitle></IonToolbar></IonHeader>
      <IonContent fullscreen className="ion-padding">
        
        {/* MODO 1: FORMULARIO */}
        {!activeTripId && (
            <IonCard>
                <IonCardHeader><IonCardTitle>Aceptar Viaje</IonCardTitle></IonCardHeader>
                <IonCardContent>
                    <div className="ion-text-center"><IonIcon icon={car} size="large"/></div>
                    <IonItem lines="none">
                        <IonLabel position="stacked">Pegar ID del Viaje</IonLabel>
                        <IonInput 
                            value={inputTripId} 
                            placeholder="Ej: 550e84..."
                            onIonInput={e => setInputTripId(e.detail.value!)} 
                        />
                    </IonItem>
                    <IonButton expand="block" className="ion-margin-top" onClick={handleAcceptTrip} disabled={loading}>
                        {loading ? "Conectando..." : "Aceptar e Iniciar"}
                    </IonButton>
                    {error && <IonNote color="danger" className="ion-padding-top" style={{display:'block'}}>{error}</IonNote>}
                </IonCardContent>
            </IonCard>
        )}

        {/* MODO 2: MONITOREO */}
        {activeTripId && (
            <IonCard color="success">
                <IonCardContent className="ion-text-center">
                    <IonIcon icon={paperPlane} size="large" />
                    <h2>Viaje Activo</h2>
                    <p>Enviando tu ubicación al pasajero...</p>
                    {/* Botón modificado para usar handleStopTrip */}
                    <IonButton fill="outline" color="light" className="ion-margin-top" onClick={handleStopTrip}>Detener</IonButton>
                </IonCardContent>
            </IonCard>
        )}

        {/* GPS DEBUG */}
        <IonCard className="gps-card">
            <IonCardHeader>
                <IonCardTitle className="gps-title">
                    <IonIcon icon={navigate} /> Mi GPS
                </IonCardTitle>
            </IonCardHeader>
            <IonCardContent>
                {currentLocation ? (
                    <div className="coords-display">
                        <div className="coord-item">
                            <span className="label">Lat:</span>
                            <span className="value">{currentLocation.latitude.toFixed(5)}</span>
                        </div>
                        <div className="coord-item">
                            <span className="label">Lng:</span>
                            <span className="value">{currentLocation.longitude.toFixed(5)}</span>
                        </div>
                    </div>
                ) : (
                     <IonNote color="warning">Buscando señal...</IonNote>
                )}
            </IonCardContent>
        </IonCard>

        <IonToast isOpen={!!toastMsg} message={toastMsg || ""} duration={3000} onDidDismiss={() => setToastMsg(null)} />
      </IonContent>
    </IonPage>
  );
};

export default Tab1;